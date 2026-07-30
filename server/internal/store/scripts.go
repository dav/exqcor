package store

import (
	"database/sql"
	"errors"
	"fmt"
)

const vosdName = "VOSD"
const vosdDescription = "The scenic narration: all stage directions are spoken by Voice of Stage Directions."

var ErrVOSDProtected = errors.New("the VOSD character cannot be deleted or renamed")

func scanScript(row interface{ Scan(...any) error }) (Script, error) {
	var s Script
	err := row.Scan(&s.ID, &s.Title, &s.Description, &s.Theme, &s.WritingSeconds, &s.StationMode, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

const scriptCols = "id, title, description, theme, writing_seconds, station_mode, created_at, updated_at"

func (s *Store) ListScripts() ([]Script, error) {
	rows, err := s.DB.Query("SELECT " + scriptCols + " FROM scripts ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Script{}
	for rows.Next() {
		sc, err := scanScript(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

func (s *Store) GetScript(id int64) (Script, error) {
	return scanScript(s.DB.QueryRow("SELECT "+scriptCols+" FROM scripts WHERE id = ?", id))
}

// CreateScript inserts a script and its VOSD character.
func (s *Store) CreateScript(title, description, theme string, writingSeconds int, stationMode string) (Script, error) {
	if writingSeconds <= 0 {
		writingSeconds = 300
	}
	if stationMode == "" {
		stationMode = "station"
	}
	now := Now()
	res, err := s.DB.Exec(
		"INSERT INTO scripts (title, description, theme, writing_seconds, station_mode, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		title, description, theme, writingSeconds, stationMode, now, now)
	if err != nil {
		return Script{}, err
	}
	id, _ := res.LastInsertId()
	if _, err := s.ensureVOSD(id); err != nil {
		return Script{}, err
	}
	return s.GetScript(id)
}

func (s *Store) UpdateScript(id int64, title, description, theme string, writingSeconds int, stationMode string) (Script, error) {
	_, err := s.DB.Exec(
		"UPDATE scripts SET title = ?, description = ?, theme = ?, writing_seconds = ?, station_mode = ?, updated_at = ? WHERE id = ?",
		title, description, theme, writingSeconds, stationMode, Now(), id)
	if err != nil {
		return Script{}, err
	}
	return s.GetScript(id)
}

func (s *Store) DeleteScript(id int64) error {
	_, err := s.DB.Exec("DELETE FROM scripts WHERE id = ?", id)
	return err
}

func (s *Store) ensureVOSD(scriptID int64) (int64, error) {
	var id int64
	err := s.DB.QueryRow("SELECT id FROM characters WHERE script_id = ? AND role = 'vosd'", scriptID).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	now := Now()
	res, err := s.DB.Exec(
		"INSERT INTO characters (script_id, name, description, role, created_at, updated_at) VALUES (?, ?, ?, 'vosd', ?, ?)",
		scriptID, vosdName, vosdDescription, now, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// DuplicateScript copies a script as a reusable template: characters (with
// actor links), sections with on_stage flags, and each section's priming line
// (first line of the first sub_section, re-attributed to the new VOSD).
func (s *Store) DuplicateScript(id int64, title string) (Script, error) {
	src, err := s.GetScript(id)
	if err != nil {
		return Script{}, err
	}
	if title == "" {
		title = src.Title + " (copy)"
	}
	dst, err := s.CreateScript(title, src.Description, src.Theme, src.WritingSeconds, src.StationMode)
	if err != nil {
		return Script{}, err
	}

	chars, err := s.ListCharacters(id)
	if err != nil {
		return Script{}, err
	}
	charMap := map[int64]int64{} // old character id -> new character id
	newVOSD, err := s.ensureVOSD(dst.ID)
	if err != nil {
		return Script{}, err
	}
	for _, c := range chars {
		if c.Role == "vosd" {
			charMap[c.ID] = newVOSD
			continue
		}
		nc, err := s.CreateCharacter(dst.ID, c.ActorID, c.Name, c.Description)
		if err != nil {
			return Script{}, err
		}
		charMap[c.ID] = nc.ID
	}

	sections, err := s.ListSections(id)
	if err != nil {
		return Script{}, err
	}
	for _, sec := range sections {
		nsec, err := s.CreateSection(dst.ID, sec.Name, sec.Description, &sec.Ordering)
		if err != nil {
			return Script{}, err
		}
		css, err := s.ListCharacterSections(sec.ID)
		if err != nil {
			return Script{}, err
		}
		for _, cs := range css {
			newChar, ok := charMap[cs.CharacterID]
			if !ok {
				continue
			}
			if err := s.SetCharacterSection(newChar, nsec.ID, true, cs.OnStage); err != nil {
				return Script{}, err
			}
		}
		if line, err := s.PrimingLine(sec.ID); err != nil {
			return Script{}, err
		} else if line != nil {
			if err := s.SetPrimingLine(nsec.ID, newVOSD, line.Text); err != nil {
				return Script{}, err
			}
		}
	}
	return dst, nil
}

// PrimingLine returns the first line of a section's first sub_section, if any.
func (s *Store) PrimingLine(sectionID int64) (*Line, error) {
	row := s.DB.QueryRow(`
		SELECT l.id, l.sub_section_id, l.character_id, l.text, l.ordering, l.created_at
		FROM lines l
		JOIN sub_sections ss ON ss.id = l.sub_section_id
		WHERE ss.section_id = ? AND ss.ordering = (SELECT MIN(ordering) FROM sub_sections WHERE section_id = ?)
		ORDER BY l.ordering LIMIT 1`, sectionID, sectionID)
	var l Line
	err := row.Scan(&l.ID, &l.SubSectionID, &l.CharacterID, &l.Text, &l.Ordering, &l.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// SetPrimingLine sets (replacing if present) the seed line in a section's
// first sub_section — the line the first writer builds on.
func (s *Store) SetPrimingLine(sectionID, characterID int64, text string) error {
	var ssID int64
	err := s.DB.QueryRow(
		"SELECT id FROM sub_sections WHERE section_id = ? ORDER BY ordering LIMIT 1", sectionID).Scan(&ssID)
	if err != nil {
		return fmt.Errorf("section %d has no first sub_section: %w", sectionID, err)
	}
	if _, err := s.DB.Exec("DELETE FROM lines WHERE sub_section_id = ?", ssID); err != nil {
		return err
	}
	_, err = s.DB.Exec(
		"INSERT INTO lines (sub_section_id, character_id, text, ordering, created_at) VALUES (?, ?, ?, 1, ?)",
		ssID, characterID, text, Now())
	return err
}
