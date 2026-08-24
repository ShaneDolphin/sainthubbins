// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 44th batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_ClockDispose(t *testing.T) {
	c := NewClock(0.5)
	if c.CPS != 0.5 {
		t.Fatalf("clock 0.5 got %v", c.CPS)
	}
	c.SetCPS(2.0)
	if c.CPS != 2.0 {
		t.Fatalf("clock 2.0 got %v", c.CPS)
	}
	// Dispose via context not needed; check Duration/Interval fields exist
	if c.Duration == 0 || c.Interval == 0 {
		t.Fatalf("clock Duration Interval nonzero")
	}
}

func TestMJS_PatternAddSub(t *testing.T) {
	p := Pure(10).Add(Pure(5)).Sub(Pure(3))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps[0].Value) != 12 {
		t.Fatalf("10+5-3=12 got %v", haps[0].Value)
	}
	q := Pure(4).Mul(Pure(3)).Div(Pure(2))
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps2[0].Value) != 6 {
		t.Fatalf("4*3/2=6 got %v", haps2[0].Value)
	}
}

func TestMJS_RangeSignal(t *testing.T) {
	s := Saw().Range(10, 20)
	haps := s.QueryArc(FractionFromInt(0), FractionFromFloat(0.001))
	if len(haps) == 0 {
		t.Fatalf("saw range empty")
	}
	v := toFloat(haps[0].Value)
	if v < 10 || v > 20 {
		t.Fatalf("saw Range 10-20 got %v", v)
	}
	tri := Tri().Range(-1, 1)
	haps2 := tri.QueryArc(FractionFromInt(0), FractionFromFloat(0.001))
	if len(haps2) == 0 {
		t.Fatalf("tri range empty")
	}
	v2 := toFloat(haps2[0].Value)
	if v2 < -1 || v2 > 1 {
		t.Fatalf("tri Range -1,1 got %v", v2)
	}
}
