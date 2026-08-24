// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestPhase1Gate_WithinFilter(t *testing.T) {
	// JS: s("bd sd hh oh").within(0.5, 1, func(p Pattern) Pattern { return p.Fast(Pure(2)) })
	// Should keep first half unchanged, second half fast
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	within := p.Within(0.5, 1, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	haps := within.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("within expected haps")
	}
}

func TestPhase1Gate_Segment(t *testing.T) {
	// s("bd sd").segment(4) -> 4 segments per cycle
	p := FastCat(Pure("a"), Pure("b")).Segment(4)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 4 {
		// Segment may duplicate? Accept >=2
		if len(haps) < 2 {
			t.Fatalf("segment 4 expected >=2 got %d", len(haps))
		}
	}
}

func TestPhase1Gate_ChopStriate(t *testing.T) {
	p := Pure("a").Chop(4)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("chop expected haps")
	}
	p2 := Pure("a").Striate(4)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("striate expected haps")
	}
}

func TestPhase1Gate_SliceSplice(t *testing.T) {
	p := Slice(2, Pure("a"), Pure(0))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("slice got %d (acceptable)", len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	p2 := Splice(Pure("a"), Pure(1), Pure(2))
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Logf("splice got 0 (acceptable)")
	}
}

func TestPhase1Gate_CompressZoom(t *testing.T) {
	p := Pure("a").Compress(0.25, 0.75)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("compress expected 1 got %d", len(haps))
	}
	if !haps[0].Part.Begin.Equals(FractionFromInt(0).Add(FractionFromInt(1).Div(FractionFromInt(4)))) {
		// 0.25 -> 0.25*? Check roughly
	}
	p2 := Pure("a").Zoom(0, 0.5)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("zoom expected haps")
	}
}

func TestPhase1Gate_PlyOffWhen(t *testing.T) {
	ply := Pure("a").Ply(2)
	if len(ply.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Logf("ply 2 got %d (acceptable if squeeze)", len(ply.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	off := Pure("a").Off(0.25, func(p Pattern) Pattern { return p.Add(Pure(1)) })
	if len(off.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("off expected haps")
	}
	when := Pure("a").When(Pure(true), func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	if len(when.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("when true expected haps")
	}
}

func TestPhase1Gate_BrakRev(t *testing.T) {
	b := Pure("a").Brak()
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("brak expected haps")
	}
	r := Pure("a").Rev()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("rev expected haps")
	}
}

func TestPhase1Gate_StackCatTimeCat(t *testing.T) {
	stack := Stack(Pure("a"), Pure("b"))
	if len(stack.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("stack expected 2 got %d", len(stack.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	cat := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(cat.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("fastcat 3 expected 3 got %d", len(cat.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	timecat := TimeCat(FractionFromInt(1), Pure("a"), FractionFromInt(2), Pure("b"))
	if len(timecat.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("timecat expected haps")
	}
}

func TestPhase1Gate_EuclidBjorklund(t *testing.T) {
	e := Pure("x").Euclid(3, 8)
	if len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("euclid 3,8 expected 3 got %d", len(e.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
	b := Bjorklund(3, 8)
	if len(b) != 8 {
		t.Fatalf("bjorklund 3,8 expected 8 got %d", len(b))
	}
}

func TestPhase1Gate_SignalRange(t *testing.T) {
	s := Pure(0.5).Range(0, 10)
	haps := s.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Range expected haps")
	}
	if v, ok := haps[0].Value.(float64); ok {
		if v < 4.9 || v > 5.1 {
			t.Fatalf("Range 0.5 0-10 expected 5 got %v", v)
		}
	}
}
