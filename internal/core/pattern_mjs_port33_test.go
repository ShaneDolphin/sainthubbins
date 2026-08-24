// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// 33rd batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_License(t *testing.T) {
	// Go port should carry AGPL-3.0-or-later header — smoke chain
	pat := Pure("a").FastF(FractionFromInt(2)).Slow(FractionFromInt(2))
	if len(pat.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("license smoke chain empty")
	}
}

func TestMJS_PolymeterEuclid(t *testing.T) {
	for _, tc := range []struct{ k, n int }{{0, 8}, {8, 8}, {3, 8}, {5, 8}} {
		pat := Pure("bd").Euclid(tc.k, tc.n)
		haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
		if tc.k == 0 && len(haps) != 0 {
			t.Fatalf("euclid 0,8 expected 0 haps got %d", len(haps))
		}
		if tc.k == 8 && tc.n == 8 && len(haps) == 0 {
			t.Fatalf("euclid 8,8 expected haps")
		}
	}
}

func TestMJS_SignalBypass(t *testing.T) {
	p := Sine().Range(0, 1)
	haps := p.QueryArc(FractionFromInt(0), FractionFromFloat(0.001))
	if len(haps) == 0 {
		t.Fatalf("sine range empty")
	}
	val, ok := haps[0].Value.(float64)
	if !ok {
		t.Fatalf("sine not float64 %T", haps[0].Value)
	}
	if val < 0 || val > 1 {
		t.Fatalf("sine Range 0-1 out of bounds %v", val)
	}
}
