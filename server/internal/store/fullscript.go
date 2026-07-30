package store

// FullSection is a section with its complete text, assembled for script
// views and printing.
type FullSection struct {
	Section
	Props   []Prop     `json:"props"`
	Lines   []LineView `json:"lines"`
	Writers []string   `json:"writers"`
}

type FullScript struct {
	Script     Script        `json:"script"`
	Characters []Character   `json:"characters"`
	Actors     []Actor       `json:"actors"`
	Sections   []FullSection `json:"sections"`
}

// FullScriptView assembles the whole script: every section's lines in
// written order, plus cast and credited writers.
func (s *Store) FullScriptView(scriptID int64) (FullScript, error) {
	sc, err := s.GetScript(scriptID)
	if err != nil {
		return FullScript{}, err
	}
	chars, err := s.ListCharacters(scriptID)
	if err != nil {
		return FullScript{}, err
	}
	actors, err := s.ListActors()
	if err != nil {
		return FullScript{}, err
	}
	secs, err := s.ListSections(scriptID)
	if err != nil {
		return FullScript{}, err
	}
	out := FullScript{Script: sc, Characters: chars, Actors: actors}
	for _, sec := range secs {
		lines, err := s.scanLineViews(`
			SELECT l.id, l.text, l.ordering, c.id, c.name, c.role
			FROM lines l
			JOIN characters c ON c.id = l.character_id
			JOIN sub_sections ss ON ss.id = l.sub_section_id
			WHERE ss.section_id = ?
			ORDER BY ss.ordering, l.ordering`, sec.ID)
		if err != nil {
			return FullScript{}, err
		}
		props, err := s.ListProps(sec.ID)
		if err != nil {
			return FullScript{}, err
		}
		writers, err := s.sectionWriters(sec.ID)
		if err != nil {
			return FullScript{}, err
		}
		out.Sections = append(out.Sections, FullSection{
			Section: sec, Props: props, Lines: lines, Writers: writers,
		})
	}
	return out, nil
}

func (s *Store) sectionWriters(sectionID int64) ([]string, error) {
	rows, err := s.DB.Query(`
		SELECT w.name FROM sub_sections ss
		JOIN writers w ON w.id = ss.writer_id
		WHERE ss.section_id = ? AND w.name != ''
		ORDER BY ss.ordering`, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// WritingStats summarizes a show for the post-mortem: turns and duration per
// section.
type SectionStats struct {
	SectionID   int64   `json:"section_id"`
	SectionName string  `json:"section_name"`
	Turns       int     `json:"turns"`
	Lines       int     `json:"lines"`
	AvgSeconds  float64 `json:"avg_turn_seconds"`
}

func (s *Store) ScriptStats(scriptID int64) ([]SectionStats, error) {
	rows, err := s.DB.Query(`
		SELECT sec.id, sec.name,
			COUNT(DISTINCT CASE WHEN ss.writer_id IS NOT NULL THEN ss.id END),
			COUNT(l.id),
			COALESCE(AVG(CASE WHEN ss.started_at IS NOT NULL AND ss.completed_at IS NOT NULL
				THEN (julianday(ss.completed_at) - julianday(ss.started_at)) * 86400 END), 0)
		FROM sections sec
		LEFT JOIN sub_sections ss ON ss.section_id = sec.id
		LEFT JOIN lines l ON l.sub_section_id = ss.id
		WHERE sec.script_id = ?
		GROUP BY sec.id ORDER BY sec.ordering`, scriptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []SectionStats{}
	for rows.Next() {
		var st SectionStats
		if err := rows.Scan(&st.SectionID, &st.SectionName, &st.Turns, &st.Lines, &st.AvgSeconds); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}
