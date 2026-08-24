// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// SlowCat must give each argument a full cycle, not a zero-width span.

package core

import "testing"

// fourOnTheFloor has events at sub-cycle positions, which is what exposes the
// collapsed-span bug — a single whole-cycle event does not.
func fourOnTheFloor() Pattern {
	return FastCat(Pure("bd"), Pure("bd"), Pure("bd"), Pure("bd"))
}

func TestSlowCatKeepsMultiEventArguments(t *testing.T) {
	p := SlowCat(Silence(), fourOnTheFloor())
	for cycle := int64(0); cycle < 6; cycle++ {
		got := len(p.QueryArc(FractionFromInt(cycle), FractionFromInt(cycle+1)))
		want := 0
		if cycle%2 == 1 {
			want = 4
		}
		if got != want {
			t.Errorf("cycle %d: got %d haps, want %d", cycle, got, want)
		}
	}
}

func TestSlowCatWideQueryMatchesPerCycle(t *testing.T) {
	p := SlowCat(fourOnTheFloor(), Silence(), fourOnTheFloor(), Silence())
	whole := len(p.QueryArc(FractionFromInt(0), FractionFromInt(8)))
	split := 0
	for c := int64(0); c < 8; c++ {
		split += len(p.QueryArc(FractionFromInt(c), FractionFromInt(c+1)))
	}
	if whole != split {
		t.Errorf("one 8-cycle query gave %d haps, eight 1-cycle queries gave %d", whole, split)
	}
	if whole != 16 {
		// pats has period 4: cycles 0 and 2 (of every 4) carry four events
		// each, cycles 1 and 3 are silent. Over 8 cycles (two periods) that
		// is 4 four-event cycles * 4 events = 16.
		t.Errorf("got %d haps over 8 cycles, want 16 (four events on two of every four bars, times two periods)", whole)
	}
}

// Events must land at the right offset inside their cycle, not just be present.
func TestSlowCatPreservesEventPositions(t *testing.T) {
	p := SlowCat(Silence(), fourOnTheFloor())
	haps := p.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(haps) != 4 {
		t.Fatalf("got %d haps, want 4", len(haps))
	}
	want := []float64{1.0, 1.25, 1.5, 1.75}
	for i, w := range want {
		if got := haps[i].Part.Begin.Float64(); got != w {
			t.Errorf("hap %d begins at %v, want %v", i, got, w)
		}
	}
}
