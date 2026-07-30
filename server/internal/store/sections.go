package store

import (
	"database/sql"
	"errors"
)

const sectionCols = "id, script_id, name, description, ordering, status"

func scanSection(row interface{ Scan(...any) error }) (Section, error) {
	var sec Section
	err := row.Scan(&sec.ID, &sec.ScriptID, &sec.Name, &sec.Description, &sec.Ordering, &sec.Status)
	return sec, err
}

func (s *Store) ListSections(scriptID int64) ([]Section, error) {
	rows, err := s.DB.Query(
		"SELECT "+sectionCols+" FROM sections WHERE script_id = ? ORDER BY ordering, id", scriptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Section{}
	for rows.Next() {
		sec, err := scanSection(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sec)
	}
	return out, rows.Err()
}

func (s *Store) GetSection(id int64) (Section, error) {
	return scanSection(s.DB.QueryRow("SELECT "+sectionCols+" FROM sections WHERE id = ?", id))
}

// CreateSection inserts a section, its first sub_section (which will hold the
// priming line), and attaches the script's VOSD. Pass ordering nil to append.
func (s *Store) CreateSection(scriptID int64, name, description string, ordering *int) (Section, error) {
	ord := 0
	if ordering != nil {
		ord = *ordering
	} else {
		var max sql.NullInt64
		if err := s.DB.QueryRow("SELECT MAX(ordering) FROM sections WHERE script_id = ?", scriptID).Scan(&max); err != nil {
			return Section{}, err
		}
		if max.Valid {
			ord = int(max.Int64) + 1
		}
	}
	now := Now()
	res, err := s.DB.Exec(
		"INSERT INTO sections (script_id, name, description, ordering, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		scriptID, name, description, ord, now, now)
	if err != nil {
		return Section{}, err
	}
	id, _ := res.LastInsertId()
	if _, err := s.DB.Exec("INSERT INTO sub_sections (section_id, ordering) VALUES (?, 0)", id); err != nil {
		return Section{}, err
	}
	vosd, err := s.ensureVOSD(scriptID)
	if err != nil {
		return Section{}, err
	}
	if err := s.SetCharacterSection(vosd, id, true, false); err != nil {
		return Section{}, err
	}
	return s.GetSection(id)
}

func (s *Store) UpdateSection(id int64, name, description string, ordering int, status string) (Section, error) {
	_, err := s.DB.Exec(
		"UPDATE sections SET name = ?, description = ?, ordering = ?, status = ?, updated_at = ? WHERE id = ?",
		name, description, ordering, status, Now(), id)
	if err != nil {
		return Section{}, err
	}
	return s.GetSection(id)
}

func (s *Store) DeleteSection(id int64) error {
	_, err := s.DB.Exec("DELETE FROM sections WHERE id = ?", id)
	return err
}

// --- character_sections ---

func (s *Store) ListCharacterSections(sectionID int64) ([]CharacterSection, error) {
	rows, err := s.DB.Query(
		"SELECT id, character_id, section_id, on_stage FROM character_sections WHERE section_id = ?", sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []CharacterSection{}
	for rows.Next() {
		var cs CharacterSection
		if err := rows.Scan(&cs.ID, &cs.CharacterID, &cs.SectionID, &cs.OnStage); err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, rows.Err()
}

// SetCharacterSection attaches/detaches a character to a section. With
// attached=true it upserts the on_stage flag; false removes the link (VOSD
// stays attached).
func (s *Store) SetCharacterSection(characterID, sectionID int64, attached, onStage bool) error {
	if !attached {
		c, err := s.GetCharacter(characterID)
		if err != nil {
			return err
		}
		if c.Role == "vosd" {
			return ErrVOSDProtected
		}
		_, err = s.DB.Exec(
			"DELETE FROM character_sections WHERE character_id = ? AND section_id = ?", characterID, sectionID)
		return err
	}
	_, err := s.DB.Exec(`
		INSERT INTO character_sections (character_id, section_id, on_stage) VALUES (?, ?, ?)
		ON CONFLICT(character_id, section_id) DO UPDATE SET on_stage = excluded.on_stage`,
		characterID, sectionID, boolToInt(onStage))
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// --- props ---

func (s *Store) ListProps(sectionID int64) ([]Prop, error) {
	rows, err := s.DB.Query(
		"SELECT id, section_id, name, description FROM props WHERE section_id = ? ORDER BY id", sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Prop{}
	for rows.Next() {
		var p Prop
		if err := rows.Scan(&p.ID, &p.SectionID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) CreateProp(sectionID int64, name, description string) (Prop, error) {
	res, err := s.DB.Exec(
		"INSERT INTO props (section_id, name, description) VALUES (?, ?, ?)", sectionID, name, description)
	if err != nil {
		return Prop{}, err
	}
	id, _ := res.LastInsertId()
	var p Prop
	err = s.DB.QueryRow("SELECT id, section_id, name, description FROM props WHERE id = ?", id).
		Scan(&p.ID, &p.SectionID, &p.Name, &p.Description)
	return p, err
}

func (s *Store) DeleteProp(id int64) error {
	_, err := s.DB.Exec("DELETE FROM props WHERE id = ?", id)
	return err
}

// --- actors ---

func (s *Store) ListActors() ([]Actor, error) {
	rows, err := s.DB.Query("SELECT id, name, bio FROM actors ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Actor{}
	for rows.Next() {
		var a Actor
		if err := rows.Scan(&a.ID, &a.Name, &a.Bio); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) CreateActor(name, bio string) (Actor, error) {
	now := Now()
	res, err := s.DB.Exec(
		"INSERT INTO actors (name, bio, created_at, updated_at) VALUES (?, ?, ?, ?)", name, bio, now, now)
	if err != nil {
		return Actor{}, err
	}
	id, _ := res.LastInsertId()
	return Actor{ID: id, Name: name, Bio: bio}, nil
}

func (s *Store) UpdateActor(id int64, name, bio string) (Actor, error) {
	_, err := s.DB.Exec("UPDATE actors SET name = ?, bio = ?, updated_at = ? WHERE id = ?", name, bio, Now(), id)
	return Actor{ID: id, Name: name, Bio: bio}, err
}

func (s *Store) DeleteActor(id int64) error {
	_, err := s.DB.Exec("DELETE FROM actors WHERE id = ?", id)
	return err
}

// ErrNotFound reports whether err is sql.ErrNoRows.
func IsNotFound(err error) bool { return errors.Is(err, sql.ErrNoRows) }
