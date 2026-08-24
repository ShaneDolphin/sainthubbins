// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 34th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ChordVoicing(t *testing.T) {
	p := Pure("c3").Off(0.25, func(q Pattern) Pattern { return q.Add(Pure(4)) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Off empty")
	}
	p2 := Pure("c").Scale("minor")
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("scale minor empty")
	}
}

func TestMJS_ScaleQuantize(t *testing.T) {
	p := Pure("c4").Add(Pure(2)).Scale("major")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("scale major empty")
	}
}

func TestMJS_EuclidWithOff(t *testing.T) {
	p := Pure("bd").Euclid(3, 8).Off(0.25, func(q Pattern) Pattern { return q.Fmap(func(v any) any { return "sd" }) })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("euclid off empty")
	}
	plain := Pure("bd").Euclid(3, 8).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) <= len(plain) {
		t.Fatalf("Off expected more haps than plain %d <= %d", len(haps), len(plain))
	}
}
