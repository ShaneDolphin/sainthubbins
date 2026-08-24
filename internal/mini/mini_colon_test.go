// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestMiniColonSample(t *testing.T) {
	p := Mini("bd:1")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("bd:1 expected 1 got %d", len(haps))
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok {
		t.Fatalf("bd:1 value not map %T %v", haps[0].Value, haps[0].Value)
	}
	if m["s"] != "bd" || m["n"] != 1 {
		t.Fatalf("bd:1 expected s=bd n=1 got %v", m)
	}
	p2 := Mini("bd:2 sd:3")
	haps2 := p2.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("bd:2 sd:3 expected 2 got %d", len(haps2))
	}
	if m, ok := haps2[0].Value.(map[string]any); !ok || m["n"] != 2 {
		t.Fatalf("bd:2 n=2 got %v", haps2[0].Value)
	}
}

func TestMiniColonFast(t *testing.T) {
	p := Mini("bd:1*2")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("bd:1*2 expected 2 got %d", len(haps))
	}
	for _, h := range haps {
		m, ok := h.Value.(map[string]any)
		if !ok || m["s"] != "bd" || m["n"] != 1 {
			t.Fatalf("bd:1*2 haps wrong %v", h.Value)
		}
	}
}
