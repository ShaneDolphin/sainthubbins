// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestControlsSMap(t *testing.T) {
	p := S(Pure("bd"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("S bd empty")
	}
	// Check context has s or value
	found := false
	for _, h := range haps {
		if m, ok := h.Value.(map[string]any); ok && m["s"] == "bd" {
			found = true
		}
		if h.Value == "bd" {
			found = true
		}
	}
	if !found {
		t.Fatalf("S bd not found %v", haps[0].Value)
	}
}

func TestControlsNoteGain(t *testing.T) {
	p := Stack(Note("c4"), Gain(0.9))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("stack note+gain expected 2 got %d", len(haps))
	}
	foundNote := false
	foundGain := false
	for _, h := range haps {
		if m, ok := h.Value.(map[string]any); ok {
			if m["note"] == "c4" {
				foundNote = true
			}
			if m["gain"] == 0.9 || m["gain"] == float64(0.9) {
				foundGain = true
			}
		}
	}
	if !foundNote || !foundGain {
		t.Fatalf("stack note+gain missing note=%v gain=%v haps %v", foundNote, foundGain, haps)
	}
	// Single control patterns each non-empty
	if len(Note("c4").QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Note c4 expected 1")
	}
	if len(Gain(0.9).QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Gain 0.9 expected 1")
	}
}

func TestControlsCutoffStack(t *testing.T) {
	p := Stack(Cutoff(800), Cutoff(1200))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("cutoff stack expected 2 got %d", len(haps))
	}
	// Combined cutoff + s via Add
	p2 := S("bd").Add(Cutoff(800))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("s+cutoff empty")
	}
}
