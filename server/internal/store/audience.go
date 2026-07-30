package store

import (
	"database/sql"
	"errors"
)

var ErrNoOneWaiting = errors.New("no one is waiting in the queue")

const audienceCols = "id, script_id, number, name, device_token, status, called_at, called_section_id, created_at"

type audienceRow struct {
	AudienceMember
	CalledSectionID *int64 `json:"called_section_id"`
}

func scanAudience(row interface{ Scan(...any) error }) (audienceRow, error) {
	var m audienceRow
	err := row.Scan(&m.ID, &m.ScriptID, &m.Number, &m.Name, &m.DeviceToken, &m.Status, &m.CalledAt, &m.CalledSectionID, &m.CreatedAt)
	return m, err
}

// AudienceView is an audience member plus queue-derived fields.
type AudienceView struct {
	AudienceMember
	CalledSectionID *int64 `json:"called_section_id"`
	Position        int    `json:"position"` // people waiting ahead; 0 when not waiting
}

// JoinAudience finds or creates the audience member for a device. Rejoining
// (page reload, Wi-Fi blip) returns the same number.
func (s *Store) JoinAudience(scriptID int64, deviceToken, name string) (AudienceView, error) {
	row := s.DB.QueryRow("SELECT "+audienceCols+" FROM audience_members WHERE device_token = ? AND script_id = ?", deviceToken, scriptID)
	m, err := scanAudience(row)
	if errors.Is(err, sql.ErrNoRows) {
		var next int
		if err := s.DB.QueryRow(
			"SELECT COALESCE(MAX(number), 0) + 1 FROM audience_members WHERE script_id = ?", scriptID).Scan(&next); err != nil {
			return AudienceView{}, err
		}
		_, err := s.DB.Exec(
			"INSERT INTO audience_members (script_id, number, name, device_token, created_at) VALUES (?, ?, ?, ?, ?)",
			scriptID, next, name, deviceToken, Now())
		if err != nil {
			return AudienceView{}, err
		}
		row = s.DB.QueryRow("SELECT "+audienceCols+" FROM audience_members WHERE device_token = ? AND script_id = ?", deviceToken, scriptID)
		m, err = scanAudience(row)
		if err != nil {
			return AudienceView{}, err
		}
	} else if err != nil {
		return AudienceView{}, err
	}
	return s.withPosition(m)
}

// AudienceByDevice returns the member for a device in the active script.
func (s *Store) AudienceByDevice(scriptID int64, deviceToken string) (*AudienceView, error) {
	row := s.DB.QueryRow("SELECT "+audienceCols+" FROM audience_members WHERE device_token = ? AND script_id = ?", deviceToken, scriptID)
	m, err := scanAudience(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	v, err := s.withPosition(m)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Store) withPosition(m audienceRow) (AudienceView, error) {
	v := AudienceView{AudienceMember: m.AudienceMember, CalledSectionID: m.CalledSectionID}
	if m.Status == "waiting" {
		if err := s.DB.QueryRow(
			"SELECT COUNT(*) FROM audience_members WHERE script_id = ? AND status = 'waiting' AND number < ?",
			m.ScriptID, m.Number).Scan(&v.Position); err != nil {
			return v, err
		}
	}
	return v, nil
}

// ListAudience returns all members for a script in queue order.
func (s *Store) ListAudience(scriptID int64) ([]AudienceView, error) {
	rows, err := s.DB.Query("SELECT "+audienceCols+" FROM audience_members WHERE script_id = ? ORDER BY number", scriptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AudienceView{}
	for rows.Next() {
		m, err := scanAudience(rows)
		if err != nil {
			return nil, err
		}
		v, err := s.withPosition(m)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *Store) GetAudienceMember(id int64) (AudienceView, error) {
	m, err := scanAudience(s.DB.QueryRow("SELECT "+audienceCols+" FROM audience_members WHERE id = ?", id))
	if err != nil {
		return AudienceView{}, err
	}
	return s.withPosition(m)
}

// CallNext calls the lowest-numbered waiting member to a section.
func (s *Store) CallNext(scriptID, sectionID int64) (AudienceView, error) {
	row := s.DB.QueryRow("SELECT "+audienceCols+" FROM audience_members WHERE script_id = ? AND status = 'waiting' ORDER BY number LIMIT 1", scriptID)
	m, err := scanAudience(row)
	if errors.Is(err, sql.ErrNoRows) {
		return AudienceView{}, ErrNoOneWaiting
	}
	if err != nil {
		return AudienceView{}, err
	}
	_, err = s.DB.Exec(
		"UPDATE audience_members SET status = 'called', called_at = ?, called_section_id = ? WHERE id = ?",
		Now(), sectionID, m.ID)
	if err != nil {
		return AudienceView{}, err
	}
	return s.GetAudienceMember(m.ID)
}

// SetAudienceStatus moves a member through the queue lifecycle. Requeueing
// clears the call.
func (s *Store) SetAudienceStatus(id int64, status string) error {
	if status == "waiting" {
		_, err := s.DB.Exec(
			"UPDATE audience_members SET status = 'waiting', called_at = NULL, called_section_id = NULL WHERE id = ?", id)
		return err
	}
	_, err := s.DB.Exec("UPDATE audience_members SET status = ? WHERE id = ?", status, id)
	return err
}

// TurnAudienceMemberID returns the audience member writing a turn, if the
// turn's writer came from the audience queue.
func (s *Store) TurnAudienceMemberID(turnID int64) (*int64, error) {
	var id *int64
	err := s.DB.QueryRow(`
		SELECT w.audience_member_id FROM sub_sections ss
		JOIN writers w ON w.id = ss.writer_id
		WHERE ss.id = ?`, turnID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return id, err
}
