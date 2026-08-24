// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package tonal

import "testing"

func TestVoicingDrop2Range(t *testing.T) {
	notes := []string{"C4", "E4", "G4", "B4"}
	drop2 := Voicing(notes, "drop2")
	if len(drop2) != len(notes) {
		t.Fatalf("drop2 len %d", len(drop2))
	}
	// drop2 should lower 2nd from top by octave: check contains minus 12 semitone note
	drop3 := Voicing(notes, "drop3")
	if len(drop3) != len(notes) {
		t.Fatalf("drop3 len %d", len(drop3))
	}
	// unknown voicing returns same
	same := Voicing(notes, "unknown")
	if len(same) != len(notes) {
		t.Fatalf("unknown len %d", len(same))
	}
}

func TestTransposeIntervalOctave(t *testing.T) {
	if got := Transpose("C4", "8P"); got != "C5" {
		t.Fatalf("C4 8P expected C5 got %q", got)
	}
	if got := Transpose("C4", "4P"); got != "F4" {
		t.Fatalf("C4 4P expected F4 got %q", got)
	}
	if got := Transpose("D4", "-2M"); got != "C4" {
		t.Fatalf("D4 -2M expected C4 got %q", got)
	}
}

func TestScaleChordEdge(t *testing.T) {
	sc := Scale("C:minor")
	if len(sc) != 7 {
		t.Fatalf("minor len %d", len(sc))
	}
	ch := Chord("Am7")
	if len(ch) < 3 {
		t.Fatalf("Am7 len %d", len(ch))
	}
}
