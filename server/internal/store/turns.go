package store

import (
	"database/sql"
	"errors"
	"time"
)

// GraceSeconds is how long after ends_at the server still accepts the line a
// writer was mid-sentence on when the timer hit zero.
const GraceSeconds = 15

var (
	ErrTurnActive = errors.New("a writer is already going in this section")
	ErrTurnOver   = errors.New("time is up for this turn")
)

const subSectionCols = "id, section_id, ordering, writer_id, started_at, ends_at, completed_at"

func scanSubSection(row interface{ Scan(...any) error }) (SubSection, error) {
	var ss SubSection
	err := row.Scan(&ss.ID, &ss.SectionID, &ss.Ordering, &ss.WriterID, &ss.StartedAt, &ss.EndsAt, &ss.CompletedAt)
	return ss, err
}

// ActiveTurn returns the section's in-flight sub_section (started, not
// completed), or nil.
func (s *Store) ActiveTurn(sectionID int64) (*SubSection, error) {
	row := s.DB.QueryRow(
		"SELECT "+subSectionCols+" FROM sub_sections WHERE section_id = ? AND started_at IS NOT NULL AND completed_at IS NULL",
		sectionID)
	ss, err := scanSubSection(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ss, nil
}

// ActiveTurns returns every in-flight sub_section across all sections (used
// for crash recovery and the admin dashboard).
func (s *Store) ActiveTurns() ([]SubSection, error) {
	rows, err := s.DB.Query(
		"SELECT " + subSectionCols + " FROM sub_sections WHERE started_at IS NOT NULL AND completed_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SubSection{}
	for rows.Next() {
		ss, err := scanSubSection(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ss)
	}
	return out, rows.Err()
}

func (s *Store) GetTurn(id int64) (SubSection, error) {
	return scanSubSection(s.DB.QueryRow("SELECT "+subSectionCols+" FROM sub_sections WHERE id = ?", id))
}

// StartTurn creates and starts the section's next sub_section for a new
// writer. Fails if the section already has an active turn.
func (s *Store) StartTurn(sectionID int64, writerName string, audienceMemberID *int64, durationSeconds int) (SubSection, error) {
	active, err := s.ActiveTurn(sectionID)
	if err != nil {
		return SubSection{}, err
	}
	if active != nil {
		return SubSection{}, ErrTurnActive
	}
	res, err := s.DB.Exec(
		"INSERT INTO writers (name, audience_member_id, created_at) VALUES (?, ?, ?)",
		writerName, audienceMemberID, Now())
	if err != nil {
		return SubSection{}, err
	}
	writerID, _ := res.LastInsertId()

	var maxOrd int
	if err := s.DB.QueryRow("SELECT COALESCE(MAX(ordering), -1) FROM sub_sections WHERE section_id = ?", sectionID).Scan(&maxOrd); err != nil {
		return SubSection{}, err
	}
	now := time.Now().UTC()
	ends := now.Add(time.Duration(durationSeconds) * time.Second)
	res, err = s.DB.Exec(
		"INSERT INTO sub_sections (section_id, ordering, writer_id, started_at, ends_at) VALUES (?, ?, ?, ?, ?)",
		sectionID, maxOrd+1, writerID, now.Format(time.RFC3339), ends.Format(time.RFC3339))
	if err != nil {
		return SubSection{}, err
	}
	id, _ := res.LastInsertId()
	if _, err := s.DB.Exec("UPDATE sections SET status = 'writing', updated_at = ? WHERE id = ?", Now(), sectionID); err != nil {
		return SubSection{}, err
	}
	return s.GetTurn(id)
}

// CompleteTurn marks a turn done (writer finished, admin ended it, or the
// timer expired). Idempotent.
func (s *Store) CompleteTurn(id int64) error {
	_, err := s.DB.Exec(
		"UPDATE sub_sections SET completed_at = ? WHERE id = ? AND completed_at IS NULL", Now(), id)
	return err
}

// AddLine appends a line to an active turn, enforcing the deadline plus
// grace period server-side.
func (s *Store) AddLine(turnID, characterID int64, text string) (Line, error) {
	turn, err := s.GetTurn(turnID)
	if err != nil {
		return Line{}, err
	}
	if turn.CompletedAt != nil || turn.StartedAt == nil {
		return Line{}, ErrTurnOver
	}
	if turn.EndsAt != nil {
		ends, err := time.Parse(time.RFC3339, *turn.EndsAt)
		if err == nil && time.Now().UTC().After(ends.Add(GraceSeconds*time.Second)) {
			return Line{}, ErrTurnOver
		}
	}
	var maxOrd int
	if err := s.DB.QueryRow("SELECT COALESCE(MAX(ordering), 0) FROM lines WHERE sub_section_id = ?", turnID).Scan(&maxOrd); err != nil {
		return Line{}, err
	}
	res, err := s.DB.Exec(
		"INSERT INTO lines (sub_section_id, character_id, text, ordering, created_at) VALUES (?, ?, ?, ?, ?)",
		turnID, characterID, text, maxOrd+1, Now())
	if err != nil {
		return Line{}, err
	}
	id, _ := res.LastInsertId()
	var l Line
	err = s.DB.QueryRow("SELECT id, sub_section_id, character_id, text, ordering, created_at FROM lines WHERE id = ?", id).
		Scan(&l.ID, &l.SubSectionID, &l.CharacterID, &l.Text, &l.Ordering, &l.CreatedAt)
	return l, err
}

// LineView is a line joined with its speaker, ready for display.
type LineView struct {
	ID            int64  `json:"id"`
	Text          string `json:"text"`
	Ordering      int    `json:"ordering"`
	CharacterID   int64  `json:"character_id"`
	CharacterName string `json:"character_name"`
	CharacterRole string `json:"character_role"`
}

func (s *Store) scanLineViews(query string, args ...any) ([]LineView, error) {
	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []LineView{}
	for rows.Next() {
		var lv LineView
		if err := rows.Scan(&lv.ID, &lv.Text, &lv.Ordering, &lv.CharacterID, &lv.CharacterName, &lv.CharacterRole); err != nil {
			return nil, err
		}
		out = append(out, lv)
	}
	return out, rows.Err()
}

// TurnLines returns the lines written in one turn, in order.
func (s *Store) TurnLines(turnID int64) ([]LineView, error) {
	return s.scanLineViews(`
		SELECT l.id, l.text, l.ordering, c.id, c.name, c.role
		FROM lines l JOIN characters c ON c.id = l.character_id
		WHERE l.sub_section_id = ? ORDER BY l.ordering`, turnID)
}

// LastVisibleLine returns the exquisite-corpse hand-off: the single most
// recent line written in the section before the given turn (the previous
// writer's last line, or the priming line).
func (s *Store) LastVisibleLine(sectionID int64, beforeOrdering int) (*LineView, error) {
	lines, err := s.scanLineViews(`
		SELECT l.id, l.text, l.ordering, c.id, c.name, c.role
		FROM lines l
		JOIN characters c ON c.id = l.character_id
		JOIN sub_sections ss ON ss.id = l.sub_section_id
		WHERE ss.section_id = ? AND ss.ordering < ?
		ORDER BY ss.ordering DESC, l.ordering DESC LIMIT 1`, sectionID, beforeOrdering)
	if err != nil || len(lines) == 0 {
		return nil, err
	}
	return &lines[0], nil
}

// OnStageCharacters returns the section's pickable speakers: VOSD first,
// then characters flagged on_stage.
func (s *Store) OnStageCharacters(sectionID int64) ([]Character, error) {
	rows, err := s.DB.Query(`
		SELECT `+characterCols+` FROM characters
		WHERE id IN (
			SELECT character_id FROM character_sections
			WHERE section_id = ? AND (on_stage = 1 OR character_id IN
				(SELECT id FROM characters WHERE role = 'vosd'))
		)
		ORDER BY role = 'vosd' DESC, name`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Character{}
	for rows.Next() {
		c, err := scanCharacter(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// TurnCount reports how many writer turns a section has had (excluding the
// priming sub_section 0).
func (s *Store) TurnCount(sectionID int64) (int, error) {
	var n int
	err := s.DB.QueryRow(
		"SELECT COUNT(*) FROM sub_sections WHERE section_id = ? AND writer_id IS NOT NULL", sectionID).Scan(&n)
	return n, err
}

// SectionScript returns the script a section belongs to.
func (s *Store) SectionScript(sectionID int64) (Script, error) {
	sec, err := s.GetSection(sectionID)
	if err != nil {
		return Script{}, err
	}
	return s.GetScript(sec.ScriptID)
}
