// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import (
	"math"
	"testing"
)

func TestNewOps_ModPow(t *testing.T) {
	if v := Pure(7).Mod(Pure(3)).FirstCycle()[0].Value; v != math.Mod(7, 3) {
		t.Fatalf("mod expected %v got %v", math.Mod(7, 3), v)
	}
	if v := Pure(2).Pow(Pure(3)).FirstCycle()[0].Value; v != 8.0 {
		t.Fatalf("pow expected 8 got %v", v)
	}
}

func TestNewOps_Bitwise(t *testing.T) {
	if v := Pure(6).Band(Pure(3)).FirstCycle()[0].Value; v != 2 {
		t.Fatalf("band 6&3 expected 2 got %v", v)
	}
	if v := Pure(6).Bor(Pure(1)).FirstCycle()[0].Value; v != 7 {
		t.Fatalf("bor expected 7 got %v", v)
	}
	if v := Pure(6).Bxor(Pure(3)).FirstCycle()[0].Value; v != 5 {
		t.Fatalf("bxor expected 5 got %v", v)
	}
	if v := Pure(4).Blshift(Pure(1)).FirstCycle()[0].Value; v != 8 {
		t.Fatalf("blshift expected 8 got %v", v)
	}
	if v := Pure(8).Brshift(Pure(1)).FirstCycle()[0].Value; v != 4 {
		t.Fatalf("brshift expected 4 got %v", v)
	}
}

func TestNewOps_Compare(t *testing.T) {
	if v := Pure(2).Lt(Pure(3)).FirstCycle()[0].Value; v != true {
		t.Fatalf("lt expected true got %v", v)
	}
	if v := Pure(5).Gt(Pure(3)).FirstCycle()[0].Value; v != true {
		t.Fatalf("gt expected true got %v", v)
	}
	if v := Pure(3).Lte(Pure(3)).FirstCycle()[0].Value; v != true {
		t.Fatalf("lte expected true got %v", v)
	}
	if v := Pure(3).Gte(Pure(4)).FirstCycle()[0].Value; v != false {
		t.Fatalf("gte expected false got %v", v)
	}
	if v := Pure(3).Eq(Pure(3)).FirstCycle()[0].Value; v != true {
		t.Fatalf("eq expected true got %v", v)
	}
	if v := Pure(3).Ne(Pure(4)).FirstCycle()[0].Value; v != true {
		t.Fatalf("ne expected true got %v", v)
	}
}

func TestNewOps_AndOr(t *testing.T) {
	if v := Pure(true).And(Pure(false)).FirstCycle()[0].Value; v != false {
		t.Fatalf("and true&&false expected false got %v", v)
	}
	if v := Pure(true).Or(Pure(false)).FirstCycle()[0].Value; v != true {
		t.Fatalf("or expected true got %v", v)
	}
	if v := Pure(1).And(Pure(0)).FirstCycle()[0].Value; v != false {
		t.Fatalf("and numeric 1&&0 expected false got %v", v)
	}
}

func TestNewOps_ChopStriate(t *testing.T) {
	pat := Pure(map[string]any{"s": "bd"}).Chop(2)
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("chop 2 expected 2 haps got %d %v", len(haps), haps)
	}
	if haps[0].Value.(map[string]any)["begin"] != 0.0 {
		t.Fatalf("chop begin 0 expected 0 got %v", haps[0].Value)
	}
	pat2 := Pure(map[string]any{"s": "bd"}).Striate(4)
	haps2 := pat2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("striate expected >0 haps got 0")
	}
}

func TestNewOps_SliceSplice(t *testing.T) {
	sliced := Slice(4, 2, map[string]any{"s": "bd"})
	haps := sliced.FirstCycle()
	if len(haps) == 0 {
		t.Fatalf("slice expected haps")
	}
	if m, ok := haps[0].Value.(map[string]any); ok {
		if m["begin"] != 0.5 {
			t.Fatalf("slice 4 idx2 begin expected 0.5 got %v", m["begin"])
		}
	}
	spliced := Splice(4, 0, "bd")
	if len(spliced.FirstCycle()) == 0 {
		t.Fatalf("splice expected haps")
	}
}

func TestNewOps_FitLoopAt(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "begin": 0.0, "end": 0.5}).Fit()
	if len(p.FirstCycle()) != 1 {
		t.Fatalf("fit expected 1")
	}
	l := Pure(map[string]any{"s": "bd"}).LoopAt(2)
	if len(l.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("loopAt expected haps")
	}
}

func TestNewOps_BiteLingerSegment(t *testing.T) {
	b := Bite(4, 1, Pure("a"))
	if len(b.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("bite expected haps")
	}
	li := Pure("a").Linger(0.5)
	if len(li.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("linger expected haps")
	}
	s := Pure("a").Segment(4)
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("segment expected haps")
	}
}

func TestNewOps_FilterWithin(t *testing.T) {
	p := Stack(Pure(1), Pure(2))
	f := p.Filter(func(h Hap) bool { return h.Value == 1 })
	if len(f.FirstCycle()) != 1 {
		t.Fatalf("filter expected 1 got %d", len(f.FirstCycle()))
	}
	w := Pure(1).Within(0, 0.5, func(pat Pattern) Pattern { return pat.Fast(Pure(2)) })
	if len(w.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("within expected haps")
	}
}

func TestNewOps_InvertBrak(t *testing.T) {
	inv := Pure(true).Invert()
	if inv.FirstCycle()[0].Value != false {
		t.Fatalf("invert expected false")
	}
	braked := Pure("bd").Brak()
	if len(braked.QueryArc(FractionFromInt(0), FractionFromInt(2))) == 0 {
		t.Fatalf("brak expected haps")
	}
}

func TestNewOps_PressHush(t *testing.T) {
	pressed := Pure("bd").Press()
	if len(pressed.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("press expected haps")
	}
	hushed := Pure("bd").Hush()
	if len(hushed.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("hush expected silence")
	}
}
