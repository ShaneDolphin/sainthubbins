package tonal

import "testing"

func TestScale(t *testing.T) {
	notes := Scale("C:major")
	if len(notes) != 7 {
		t.Fatalf("scale expected 7 got %d", len(notes))
	}
	if notes[0] != "C" {
		t.Fatalf("scale first %v", notes[0])
	}
	notes2 := Scale("A:minor")
	if len(notes2) != 7 || notes2[0] != "A" {
		t.Fatalf("A:minor scale %v", notes2)
	}
	if notes2[1] != "B" || notes2[2] != "C" {
		t.Fatalf("A:minor intervals %v", notes2)
	}
}

func TestChord(t *testing.T) {
	notes := Chord("Cmaj7")
	if len(notes) != 4 {
		t.Fatalf("chord expected 4 got %d", len(notes))
	}
	notes2 := Chord("C")
	if len(notes2) != 3 {
		t.Fatalf("chord C expected 3 got %d", len(notes2))
	}
}
