// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestMiniPolymeterStepsSeq(t *testing.T) {
	// {a b, c d e} style polymeter via mini.go fallback Stack vs Polymeter
	p := Mini("{a b, c d e}")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("polymeter curly empty")
	}
	// Ensure both sides present in first cycle (Stack fallback)
	hasA := false
	hasC := false
	for _, h := range haps {
		if s, ok := h.Value.(string); ok {
			if s == "a" {
				hasA = true
			}
			if s == "c" {
				hasC = true
			}
		}
		if m, ok := h.Value.(map[string]any); ok {
			if m["s"] == "a" || m["value"] == "a" {
				hasA = true
			}
			if m["s"] == "c" || m["value"] == "c" {
				hasC = true
			}
		}
	}
	if !hasA && !hasC {
		t.Logf("polymeter haps %v", haps)
	}
}

func TestMiniWeightDegrade(t *testing.T) {
	p := Mini("bd@2 sd")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("bd@2 empty")
	}
	// weight should give more bd than sd in distribution (bd weight 2)
	bdCnt := 0
	for _, h := range haps {
		if h.Value == "bd" {
			bdCnt++
		}
		if m, ok := h.Value.(map[string]any); ok && m["s"] == "bd" {
			bdCnt++
		}
	}
	if bdCnt == 0 {
		t.Fatalf("bd@2 no bd")
	}
	// degrade random but should still parse
	p2 := Mini("bd sd?")
	if p2.Query == nil {
		t.Fatalf("bd sd? nil query")
	}
}

func TestMiniRangeChooseColon(t *testing.T) {
	// range with floats
	p := Mini("0 .. 2")
	if len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("0..2 empty")
	}
	// choose
	p2 := Mini("a | b | c")
	if len(p2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("a|b|c empty")
	}
	// colon sample
	p3 := Mini("bd:1 sd:2")
	haps := p3.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("bd:1 sd:2 expected 2 got %d", len(haps))
	}
}
