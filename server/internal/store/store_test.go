package store

import (
	"path/filepath"
	"testing"
)

func testStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestCreateScriptAutoCreatesVOSD(t *testing.T) {
	s := testStore(t)
	sc, err := s.CreateScript("Noir Night", "", "Film Noir", 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if sc.WritingSeconds != 300 || sc.StationMode != "station" {
		t.Errorf("defaults not applied: %+v", sc)
	}
	chars, err := s.ListCharacters(sc.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(chars) != 1 || chars[0].Role != "vosd" || chars[0].Name != "VOSD" {
		t.Fatalf("expected exactly one VOSD, got %+v", chars)
	}
}

func TestVOSDProtection(t *testing.T) {
	s := testStore(t)
	sc, _ := s.CreateScript("Show", "", "", 0, "")
	chars, _ := s.ListCharacters(sc.ID)
	vosd := chars[0]

	if err := s.DeleteCharacter(vosd.ID); err != ErrVOSDProtected {
		t.Errorf("delete VOSD: want ErrVOSDProtected, got %v", err)
	}
	if _, err := s.UpdateCharacter(vosd.ID, nil, "Renamed", ""); err != ErrVOSDProtected {
		t.Errorf("rename VOSD: want ErrVOSDProtected, got %v", err)
	}
	// Updating the description while keeping the name is allowed.
	if _, err := s.UpdateCharacter(vosd.ID, nil, vosd.Name, "new desc"); err != nil {
		t.Errorf("update VOSD description: %v", err)
	}
	// Detaching VOSD from a section is not allowed.
	sec, _ := s.CreateSection(sc.ID, "Act 1", "", nil)
	if err := s.SetCharacterSection(vosd.ID, sec.ID, false, false); err != ErrVOSDProtected {
		t.Errorf("detach VOSD: want ErrVOSDProtected, got %v", err)
	}
}

func TestOnlyOneVOSDPerScript(t *testing.T) {
	s := testStore(t)
	sc, _ := s.CreateScript("Show", "", "", 0, "")
	_, err := s.DB.Exec(
		"INSERT INTO characters (script_id, name, role, created_at, updated_at) VALUES (?, 'VOSD2', 'vosd', '', '')", sc.ID)
	if err == nil {
		t.Fatal("second vosd insert should violate the partial unique index")
	}
}

func TestCreateSectionBootstraps(t *testing.T) {
	s := testStore(t)
	sc, _ := s.CreateScript("Show", "", "", 0, "")
	sec, err := s.CreateSection(sc.ID, "Act 1", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	var ssCount int
	s.DB.QueryRow("SELECT COUNT(*) FROM sub_sections WHERE section_id = ?", sec.ID).Scan(&ssCount)
	if ssCount != 1 {
		t.Errorf("expected 1 bootstrap sub_section, got %d", ssCount)
	}

	css, _ := s.ListCharacterSections(sec.ID)
	if len(css) != 1 {
		t.Fatalf("expected VOSD attached to new section, got %+v", css)
	}

	// Appended sections get increasing ordering.
	sec2, _ := s.CreateSection(sc.ID, "Act 2", "", nil)
	if sec2.Ordering != sec.Ordering+1 {
		t.Errorf("ordering: got %d, want %d", sec2.Ordering, sec.Ordering+1)
	}
}

func TestPrimingLine(t *testing.T) {
	s := testStore(t)
	sc, _ := s.CreateScript("Show", "", "", 0, "")
	sec, _ := s.CreateSection(sc.ID, "Act 1", "", nil)
	chars, _ := s.ListCharacters(sc.ID)
	vosd := chars[0]

	if err := s.SetPrimingLine(sec.ID, vosd.ID, "A dark and stormy night."); err != nil {
		t.Fatal(err)
	}
	// Setting again replaces rather than appends.
	if err := s.SetPrimingLine(sec.ID, vosd.ID, "A bright and calm morning."); err != nil {
		t.Fatal(err)
	}
	line, err := s.PrimingLine(sec.ID)
	if err != nil || line == nil {
		t.Fatalf("primingLine: %v, %v", line, err)
	}
	if line.Text != "A bright and calm morning." {
		t.Errorf("got %q", line.Text)
	}
	var count int
	s.DB.QueryRow("SELECT COUNT(*) FROM lines").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 line after replace, got %d", count)
	}
}

func TestDuplicateScript(t *testing.T) {
	s := testStore(t)
	sc, _ := s.CreateScript("Template", "desc", "SciFi", 240, "byod")
	actor, _ := s.CreateActor("Sam", "")
	det, _ := s.CreateCharacter(sc.ID, &actor.ID, "Detective", "hard-boiled")
	_, _ = s.CreateCharacter(sc.ID, nil, "Femme Fatale", "")
	sec, _ := s.CreateSection(sc.ID, "Act 1", "opening", nil)
	chars, _ := s.ListCharacters(sc.ID)
	var vosd Character
	for _, c := range chars {
		if c.Role == "vosd" {
			vosd = c
		}
	}
	s.SetCharacterSection(det.ID, sec.ID, true, true)
	s.SetPrimingLine(sec.ID, vosd.ID, "It was raining.")

	dup, err := s.DuplicateScript(sc.ID, "")
	if err != nil {
		t.Fatal(err)
	}
	if dup.Title != "Template (copy)" || dup.Theme != "SciFi" || dup.WritingSeconds != 240 || dup.StationMode != "byod" {
		t.Errorf("script fields not copied: %+v", dup)
	}

	dchars, _ := s.ListCharacters(dup.ID)
	if len(dchars) != 3 { // VOSD + 2
		t.Fatalf("expected 3 characters, got %+v", dchars)
	}
	for _, c := range dchars {
		if c.Name == "Detective" && (c.ActorID == nil || *c.ActorID != actor.ID) {
			t.Errorf("actor link not carried: %+v", c)
		}
	}

	dsecs, _ := s.ListSections(dup.ID)
	if len(dsecs) != 1 || dsecs[0].Name != "Act 1" {
		t.Fatalf("sections not copied: %+v", dsecs)
	}
	css, _ := s.ListCharacterSections(dsecs[0].ID)
	onStage := 0
	for _, cs := range css {
		if cs.OnStage {
			onStage++
		}
	}
	if onStage != 1 {
		t.Errorf("on_stage flags not copied: %+v", css)
	}
	line, _ := s.PrimingLine(dsecs[0].ID)
	if line == nil || line.Text != "It was raining." {
		t.Errorf("priming line not copied: %+v", line)
	}
	// Priming line must be re-attributed to the duplicate's own VOSD.
	var role string
	var scriptID int64
	s.DB.QueryRow("SELECT role, script_id FROM characters WHERE id = ?", line.CharacterID).Scan(&role, &scriptID)
	if role != "vosd" || scriptID != dup.ID {
		t.Errorf("priming line attributed to character role=%s script=%d", role, scriptID)
	}
}
