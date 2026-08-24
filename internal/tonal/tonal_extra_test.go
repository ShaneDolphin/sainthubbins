package tonal

import "testing"

func TestScaleMajor(t *testing.T) {
	s := Scale("C:major")
	if len(s) != 7 {
		t.Fatalf("C:major expected 7 got %d %v", len(s), s)
	}
	if s[0] != "C" || s[1] != "D" || s[2] != "E" {
		t.Fatalf("C:major expected C D E got %v", s)
	}
}

func TestScaleMinor(t *testing.T) {
	s := Scale("A:minor")
	if len(s) != 7 {
		t.Fatalf("A:minor expected 7 got %v", s)
	}
}

func TestScaleLocrian(t *testing.T) {
	s := Scale("C:locrian")
	if len(s) != 7 {
		t.Fatalf("locrian expected 7 got %v", s)
	}
}

func TestScalePentatonic(t *testing.T) {
	s := Scale("C:pentatonic")
	if len(s) != 5 {
		t.Fatalf("pentatonic expected 5 got %v", s)
	}
}

func TestChordMaj7(t *testing.T) {
	c := Chord("Cmaj7")
	if len(c) != 4 {
		t.Fatalf("Cmaj7 expected 4 got %v", c)
	}
}

func TestChordMin7(t *testing.T) {
	c := Chord("Cm7")
	if len(c) != 4 {
		t.Fatalf("Cm7 expected 4 got %v", c)
	}
}

func TestChordDim(t *testing.T) {
	c := Chord("Cdim")
	if len(c) != 3 {
		t.Fatalf("Cdim expected 3 got %v", c)
	}
}

func TestVoicingDrop2(t *testing.T) {
	c := Chord("Cmaj7")
	v := Voicing(c, "drop2")
	if len(v) != 4 {
		t.Fatalf("drop2 expected 4 got %v", v)
	}
	if v[0] == c[0] {
		t.Logf("drop2 may rotate")
	}
}
