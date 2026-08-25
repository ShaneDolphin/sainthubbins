// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func hapWith(v any) core.Hap {
	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	return core.Hap{Whole: &span, Part: span, Value: v}
}

func TestHapToNoteFromNoteName(t *testing.T) {
	note, _, _, ok := HapToNote(hapWith(map[string]any{"note": "c4"}))
	if !ok {
		t.Fatal("c4 should resolve to a note")
	}
	if note != 60 {
		t.Errorf("c4 = %d, want 60", note)
	}
}

func TestHapToNoteNumericWins(t *testing.T) {
	note, _, _, ok := HapToNote(hapWith(map[string]any{"n": 72, "note": "c4"}))
	if !ok || note != 72 {
		t.Errorf("n should take precedence: got %d, ok=%v", note, ok)
	}
}

func TestHapToNoteVelocityFromGain(t *testing.T) {
	_, vel, _, _ := HapToNote(hapWith(map[string]any{"note": 60, "gain": 0.5}))
	if vel != 63 {
		t.Errorf("velocity = %d, want 63 (gain 0.5 of 127)", vel)
	}
}

func TestHapToNoteSkipsPitchlessEvents(t *testing.T) {
	if _, _, _, ok := HapToNote(hapWith(map[string]any{"gain": 0.5})); ok {
		t.Error("an event with no pitch should be skipped, not defaulted")
	}
}

func TestRenderMIDIProducesPairedNoteEvents(t *testing.T) {
	pat := core.Note(core.FastCat(core.Pure(60), core.Pure(64)))
	f := RenderMIDI(pat, 1, 480)
	if len(f.events) != 4 {
		t.Fatalf("got %d events, want 4 (two notes on and off)", len(f.events))
	}
}

// TestHapToNoteBareNoteNames guards against a private note parser silently
// dropping notes that the audio renderer accepts: bare pitch classes with no
// octave digit must resolve, defaulting to octave 4 like the audio paths.
func TestHapToNoteBareNoteNames(t *testing.T) {
	cases := map[string]int{"c": 60, "e": 64, "g": 67}
	for name, want := range cases {
		note, _, _, ok := HapToNote(hapWith(map[string]any{"note": name}))
		if !ok {
			t.Errorf("%q should resolve to a note, got ok=false", name)
			continue
		}
		if note != want {
			t.Errorf("%q = %d, want %d", name, note, want)
		}
	}
}

// TestHapToNoteSharpFlatSuffixNames guards against a private note parser that
// only understands "#"/"b" accidentals: "s"/"f" suffix spellings (as used
// elsewhere in the engine, e.g. core.NoteToMidi) must also resolve.
func TestHapToNoteSharpFlatSuffixNames(t *testing.T) {
	cases := map[string]int{"cs4": 61, "ef4": 63}
	for name, want := range cases {
		note, _, _, ok := HapToNote(hapWith(map[string]any{"note": name}))
		if !ok {
			t.Errorf("%q should resolve to a note, got ok=false", name)
			continue
		}
		if note != want {
			t.Errorf("%q = %d, want %d", name, note, want)
		}
	}
}

// TestHapToNoteC4StillWorks is the regression guard for the already-working
// explicit-octave path once note-name parsing routes through core.NoteToMidi.
func TestHapToNoteC4StillWorks(t *testing.T) {
	note, _, _, ok := HapToNote(hapWith(map[string]any{"note": "c4"}))
	if !ok || note != 60 {
		t.Errorf("c4 = %d, ok=%v, want 60, true", note, ok)
	}
}

// TestHapToNoteDrumNameResolvesViaDrumMapNotNoteParser confirms drum names
// still resolve via drumNotes and not core.NoteToMidi, and land on channel 9.
func TestHapToNoteDrumNameResolvesViaDrumMapNotNoteParser(t *testing.T) {
	note, _, ch, ok := HapToNote(hapWith("bd"))
	if !ok {
		t.Fatal("bd should resolve via the drum map")
	}
	if note != 36 {
		t.Errorf("bd note = %d, want 36 (drum map, not note parser)", note)
	}
	if ch != 9 {
		t.Errorf("bd channel = %d, want 9 (percussion channel)", ch)
	}
}
