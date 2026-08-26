// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func countHaps(t *testing.T, code string) int {
	t.Helper()
	p, err := Evaluate(code)
	if err != nil {
		t.Fatalf("Evaluate(%q): %v", code, err)
	}
	return len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1)))
}

func TestFastDoublesEvents(t *testing.T) {
	if got := countHaps(t, `s("bd sd").fast(2)`); got != 4 {
		t.Errorf("got %d haps, want 4", got)
	}
}

func TestChainedCalls(t *testing.T) {
	p, err := Evaluate(`s("bd").gain(0.5).cutoff(800)`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	m := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0].Value.(map[string]any)
	if m["s"] != "bd" || m["gain"] != 0.5 || m["cutoff"] != 800.0 {
		t.Errorf("got %v, want s:bd gain:0.5 cutoff:800", m)
	}
}

func TestUnaryTransform(t *testing.T) {
	p, err := Evaluate(`s("bd sd").rev()`)
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	first := haps[0].Value.(map[string]any)
	if first["s"] != "bd" || haps[0].Part.Begin.Float64() != 0.5 {
		t.Errorf("rev did not move bd to the second half: %v @ %v",
			first, haps[0].Part.Begin.Float64())
	}
}

func TestEuclidFromJS(t *testing.T) {
	if got := countHaps(t, `s("bd").euclid(3, 8)`); got != 3 {
		t.Errorf("got %d haps, want 3", got)
	}
}
