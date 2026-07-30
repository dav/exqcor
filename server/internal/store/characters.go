package store

const characterCols = "id, script_id, actor_id, name, description, role"

func scanCharacter(row interface{ Scan(...any) error }) (Character, error) {
	var c Character
	err := row.Scan(&c.ID, &c.ScriptID, &c.ActorID, &c.Name, &c.Description, &c.Role)
	return c, err
}

func (s *Store) ListCharacters(scriptID int64) ([]Character, error) {
	rows, err := s.DB.Query(
		"SELECT "+characterCols+" FROM characters WHERE script_id = ? ORDER BY role DESC, name", scriptID)
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

func (s *Store) GetCharacter(id int64) (Character, error) {
	return scanCharacter(s.DB.QueryRow("SELECT "+characterCols+" FROM characters WHERE id = ?", id))
}

func (s *Store) CreateCharacter(scriptID int64, actorID *int64, name, description string) (Character, error) {
	now := Now()
	res, err := s.DB.Exec(
		"INSERT INTO characters (script_id, actor_id, name, description, role, created_at, updated_at) VALUES (?, ?, ?, ?, 'character', ?, ?)",
		scriptID, actorID, name, description, now, now)
	if err != nil {
		return Character{}, err
	}
	id, _ := res.LastInsertId()
	return s.GetCharacter(id)
}

func (s *Store) UpdateCharacter(id int64, actorID *int64, name, description string) (Character, error) {
	cur, err := s.GetCharacter(id)
	if err != nil {
		return Character{}, err
	}
	if cur.Role == "vosd" && name != cur.Name {
		return Character{}, ErrVOSDProtected
	}
	_, err = s.DB.Exec(
		"UPDATE characters SET actor_id = ?, name = ?, description = ?, updated_at = ? WHERE id = ?",
		actorID, name, description, Now(), id)
	if err != nil {
		return Character{}, err
	}
	return s.GetCharacter(id)
}

func (s *Store) DeleteCharacter(id int64) error {
	c, err := s.GetCharacter(id)
	if err != nil {
		return err
	}
	if c.Role == "vosd" {
		return ErrVOSDProtected
	}
	_, err = s.DB.Exec("DELETE FROM characters WHERE id = ?", id)
	return err
}
