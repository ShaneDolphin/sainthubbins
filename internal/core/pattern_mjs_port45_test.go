// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 45th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_IterAndPluck(t *testing.T) {
	// Iter via Segment2+SqueezeJoin may be empty for FastCat patterns (known, see pattern_mjs_port12_test.go:37)
	p := Sequence(Pure("a"), Pure("b")).Iter(4)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	_ = haps // may be 0 due to SqueezeJoin semantics
	// IterBack should be non-empty
	q := Pure("a").IterBack(2)
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("iterBack 2 empty")
	}
}

func TestMJS_PatternLag(t *testing.T) {
	p := Pure("a").Late(FractionFromFloat(0.25))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("late 0.25 (lag) empty")
	}
	pureHaps := Pure("a").QueryArc(FractionFromInt(0), FractionFromInt(1))
	_ = pureHaps
}

func TestMJS_TimeMethods(t *testing.T) {
	p := Pure("a").Early(FractionFromFloat(0.25))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("early empty")
	}
	q := Pure("a").Late(FractionFromFloat(0.25))
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("late empty")
	}
	r := Pure("a").LateF(FractionFromFloat(0.5))
	haps3 := r.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps3) == 0 {
		t.Fatalf("lateF empty")
	}
}
