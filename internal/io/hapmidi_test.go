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
