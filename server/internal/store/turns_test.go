package store

import (
	"testing"
	"time"
)

// setupShow returns a store with one script/section, the VOSD, a character,
// and a priming line.
func setupShow(t *testing.T) (*Store, Script, Section, Character, Character) {
	t.Helper()
	s := testStore(t)
	sc, _ := s.CreateScript("Show", "", "", 60, "")
	sec, _ := s.CreateSection(sc.ID, "Act 1", "", nil)
	det, _ := s.CreateCharacter(sc.ID, nil, "Detective", "")
	chars, _ := s.ListCharacters(sc.ID)
	var vosd Character
	for _, c := range chars {
		if c.Role == "vosd" {
			vosd = c
		}
	}
	if err := s.SetCharacterSection(det.ID, sec.ID, true, true); err != nil {
		t.Fatal(err)
	}
	if err := s.SetPrimingLine(sec.ID, vosd.ID, "It was raining."); err != nil {
		t.Fatal(err)
	}
	return s, sc, sec, vosd, det
}

func TestTurnLifecycle(t *testing.T) {
	s, _, sec, vosd, det := setupShow(t)

	// First writer sees the priming line.
	turn1, err := s.StartTurn(sec.ID, "Alice", nil, 60)
	if err != nil {
		t.Fatal(err)
	}
	last, err := s.LastVisibleLine(sec.ID, turn1.Ordering)
	if err != nil || last == nil {
		t.Fatalf("LastVisibleLine: %v %v", last, err)
	}
	if last.Text != "It was raining." || last.CharacterRole != "vosd" {
		t.Errorf("expected priming line, got %+v", last)
	}

	// Only one active turn per section.
	if _, err := s.StartTurn(sec.ID, "Bob", nil, 60); err != ErrTurnActive {
		t.Errorf("second StartTurn: want ErrTurnActive, got %v", err)
	}

	// Alice writes two lines.
	if _, err := s.AddLine(turn1.ID, det.ID, "Somebody knocked."); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddLine(turn1.ID, vosd.ID, "The door creaks open."); err != nil {
		t.Fatal(err)
	}
	if err := s.CompleteTurn(turn1.ID); err != nil {
		t.Fatal(err)
	}

	// Lines rejected after completion.
	if _, err := s.AddLine(turn1.ID, det.ID, "too late"); err != ErrTurnOver {
		t.Errorf("line after completion: want ErrTurnOver, got %v", err)
	}

	// Second writer sees ONLY Alice's last line.
	turn2, err := s.StartTurn(sec.ID, "Bob", nil, 60)
	if err != nil {
		t.Fatal(err)
	}
	last, _ = s.LastVisibleLine(sec.ID, turn2.Ordering)
	if last == nil || last.Text != "The door creaks open." {
		t.Errorf("expected Alice's last line, got %+v", last)
	}

	// Turn counts exclude the priming sub_section.
	if n, _ := s.TurnCount(sec.ID); n != 2 {
		t.Errorf("TurnCount: got %d, want 2", n)
	}
}

func TestAddLineDeadline(t *testing.T) {
	s, _, sec, _, det := setupShow(t)
	turn, err := s.StartTurn(sec.ID, "Alice", nil, 60)
	if err != nil {
		t.Fatal(err)
	}
	// Force the deadline into the past, beyond grace.
	past := time.Now().UTC().Add(-(GraceSeconds + 1) * time.Second).Format(time.RFC3339)
	if _, err := s.DB.Exec("UPDATE sub_sections SET ends_at = ? WHERE id = ?", past, turn.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AddLine(turn.ID, det.ID, "too late"); err != ErrTurnOver {
		t.Errorf("want ErrTurnOver, got %v", err)
	}

	// Within grace, the line still lands.
	inGrace := time.Now().UTC().Add(-(GraceSeconds - 5) * time.Second).Format(time.RFC3339)
	s.DB.Exec("UPDATE sub_sections SET ends_at = ? WHERE id = ?", inGrace, turn.ID)
	if _, err := s.AddLine(turn.ID, det.ID, "finishing my sentence"); err != nil {
		t.Errorf("line within grace: %v", err)
	}
}

func TestOnStageCharacters(t *testing.T) {
	s, sc, sec, _, det := setupShow(t)
	// An attached-but-offstage character should not be pickable.
	extra, _ := s.CreateCharacter(sc.ID, nil, "Bartender", "")
	s.SetCharacterSection(extra.ID, sec.ID, true, false)

	chars, err := s.OnStageCharacters(sec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(chars) != 2 {
		t.Fatalf("expected VOSD + Detective, got %+v", chars)
	}
	if chars[0].Role != "vosd" {
		t.Errorf("VOSD should sort first, got %+v", chars[0])
	}
	if chars[1].ID != det.ID {
		t.Errorf("expected Detective, got %+v", chars[1])
	}
}
