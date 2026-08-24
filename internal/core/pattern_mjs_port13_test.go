// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Thirteenth batch of js/packages/core/test/pattern.test.mjs port.

package core

import "testing"

func TestMJS_BjorklundEuclid(t *testing.T) {
	bj := Bjorklund(3, 8)
	if len(bj) != 8 {
		t.Fatalf("bjorklund 3,8 len %d", len(bj))
	}
	p := Pure("a").Euclid(3, 8)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("euclid empty")
	}
	p2 := Pure("a").EuclidRot(3, 8, 2)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("euclidRot empty")
	}
}

func TestMJS_SignalPerlin(t *testing.T) {
	s := Sine().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(s) == 0 {
		t.Fatalf("sine empty")
	}
	p := Perlin().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(p) == 0 {
		t.Fatalf("perlin empty")
	}
	r := Rand().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(r) == 0 {
		t.Fatalf("rand empty")
	}
}

func TestMJS_ControlsN(t *testing.T) {
	p := S("bd")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("s bd empty")
	}
	p2 := N(2)
	if len(p2.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("n 2 empty")
	}
	p3 := S("bd").Add(N(1))
	if len(p3.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("s bd + n1 empty")
	}
}
