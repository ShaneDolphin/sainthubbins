// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestHushSilence(t *testing.T) {
	p := Pure("a").Hush()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush expected 0")
	}
	s := Silence()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence expected 0")
	}
	g := Gap(2)
	if len(g.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Gap 2 expected 0")
	}
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(2)) {
		t.Fatalf("Gap steps 2")
	}
}

func TestPureHush(t *testing.T) {
	p := Pure("a")
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("Pure a expected 1")
	}
	hushed := p.Hush()
	if len(hushed.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Pure Hush expected 0")
	}
}

func TestSilenceP(t *testing.T) {
	p := Pure("a").SilenceP()
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("SilenceP expected 0")
	}
}
