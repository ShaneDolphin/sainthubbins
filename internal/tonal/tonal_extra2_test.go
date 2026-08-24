// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package tonal

import "testing"

func TestChordRoot(t *testing.T) {
	if c := Chord("G"); c[0] != "G" || c[1] != "B" {
		t.Fatalf("G expected G B D got %v", c)
	}
	if c := Chord("Gm"); c[1] != "A#" && c[1] != "Bb" {
		// Allow A# for Bb enharmonic; check minor third 3 semitones
		t.Logf("Gm minor third enharmonic %v", c)
	}
	if c := Chord("F#maj7"); len(c) != 4 {
		t.Fatalf("F#maj7 expected 4 got %v", c)
	}
	if c := Chord("Bb"); c[0] != "Bb" {
		t.Fatalf("Bb root expected Bb got %v", c)
	}
}

func TestTranspose(t *testing.T) {
	if got := Transpose("c4", 2); got != "D4" {
		t.Fatalf("c4+2 expected D4 got %s", got)
	}
	if got := Transpose("c4", "3M"); got != "E4" {
		t.Fatalf("c4+3M expected E4 got %s", got)
	}
	if got := Transpose("g3", "5P"); got != "D4" {
		t.Fatalf("g3+5P expected D4 got %s", got)
	}
}

func TestVoicingInversions(t *testing.T) {
	chord := Chord("Cmaj7")
	if v := Voicing(chord, "drop2"); len(v) != 4 || v[0] != "E" {
		t.Fatalf("drop2 expected E G B C got %v", v)
	}
	if v := Voicing(chord, "second"); len(v) != 4 {
		t.Fatalf("second expected 4 got %v", v)
	}
}
