// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package audio

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestOfflinePanGain(t *testing.T) {
	r := NewOfflineRenderer(48000)
	p := core.Pure(map[string]any{"s": "bd", "pan": 0.5, "gain": 0.8})
	buf, err := r.RenderOffline(p, 1)
	if err != nil {
		t.Fatalf("render %v", err)
	}
	if len(buf) == 0 {
		t.Fatalf("empty")
	}
	// Check non-zero
	hasSound := false
	for _, v := range buf {
		if v != 0 {
			hasSound = true
			break
		}
	}
	if !hasSound {
		t.Fatalf("silence")
	}
}

func TestNoteFreqHH(t *testing.T) {
	// noteToFreq is for notes like c5, a4 — hh is a sample, not a note, returns 220 fallback
	if f := noteToFreq("a4"); f < 430 || f > 450 {
		t.Fatalf("a4 freq %v", f)
	}
	if f := noteToFreq("c5"); f < 500 || f > 550 {
		t.Fatalf("c5 freq %v", f)
	}
	if f := noteToFreq("hh"); f != 220 {
		t.Fatalf("hh as note should fallback 220 got %v", f)
	}
}
