// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package tonal

import "testing"

func TestTonalScaleChordTranspose(t *testing.T) {
	sc := Scale("C:major")
	if len(sc) < 5 {
		t.Fatalf("C:major len %d", len(sc))
	}
	ch := Chord("C:maj7")
	if len(ch) < 3 {
		t.Fatalf("Cmaj7 len %d", len(ch))
	}
	tr := Transpose("C4", "4P")
	if tr == "" {
		t.Fatalf("transpose empty")
	}
	// Check that major scale contains C E G
	foundCEG := 0
	for _, n := range sc {
		if n == "C" || n == "E" || n == "G" {
			foundCEG++
		}
	}
	if foundCEG < 3 {
		t.Fatalf("C major missing CEG %v", sc)
	}
}

func TestTonalMidiNote(t *testing.T) {
	// midiToNote and noteToMidi are internal, test via Transpose
	// C4 up 12 semitones should be C5
	if tr := Transpose("C4", 12); tr != "C5" {
		t.Fatalf("C4 +12 expected C5 got %q", tr)
	}
	if tr := Transpose("C4", "8P"); tr != "C5" {
		t.Fatalf("C4 8P expected C5 got %q", tr)
	}
}
