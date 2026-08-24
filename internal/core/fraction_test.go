// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestFractionCreation(t *testing.T) {
	f := NewFraction(2, 4)
	if f.Num != 1 || f.Den != 2 {
		t.Fatalf("expected 1/2, got %s", f.Show())
	}
	f = NewFraction(-2, 4)
	if f.Num != -1 || f.Den != 2 {
		t.Fatalf("expected -1/2, got %s", f.Show())
	}
	f = NewFraction(0, 5)
	if f.Num != 0 || f.Den != 1 {
		t.Fatalf("expected 0/1, got %s", f.Show())
	}
}

func TestFractionArithmetic(t *testing.T) {
	a := MustParseFraction("1/3")
	b := MustParseFraction("1/6")
	if !a.Add(b).Equals(MustParseFraction("1/2")) {
		t.Fatalf("1/3+1/6 != 1/2, got %s", a.Add(b).Show())
	}
	if !a.Sub(b).Equals(MustParseFraction("1/6")) {
		t.Fatalf("1/3-1/6 != 1/6")
	}
	if !a.Mul(b).Equals(MustParseFraction("1/18")) {
		t.Fatalf("1/3*1/6 != 1/18")
	}
	if !a.Div(b).Equals(MustParseFraction("2/1")) {
		t.Fatalf("1/3 / 1/6 != 2")
	}
}

func TestFractionCompare(t *testing.T) {
	a := FractionFromInt(1)
	b := FractionFromInt(2)
	if !a.Lt(b) {
		t.Fatalf("1 < 2 failed")
	}
	if !b.Gt(a) {
		t.Fatalf("2 > 1 failed")
	}
	if !a.Lte(a) {
		t.Fatalf("1 <=1 failed")
	}
	if !a.Equals(FractionFromInt(1)) {
		t.Fatalf("equals failed")
	}
}

func TestFractionFloorSam(t *testing.T) {
	f := MustParseFraction("3/2")
	if !f.Floor().Equals(FractionFromInt(1)) {
		t.Fatalf("floor 3/2 !=1, got %s", f.Floor().Show())
	}
	if !f.Sam().Equals(FractionFromInt(1)) {
		t.Fatalf("sam 3/2 !=1")
	}
	f = MustParseFraction("-3/2")
	if !f.Floor().Equals(FractionFromInt(-2)) {
		t.Fatalf("floor -3/2 != -2, got %s", f.Floor().Show())
	}
	if !MustParseFraction("1/3").CyclePos().Equals(MustParseFraction("1/3")) {
		t.Fatalf("cyclePos failed")
	}
}

func TestFractionMod(t *testing.T) {
	a := MustParseFraction("5/2")
	m := FractionFromInt(2)
	if !a.Mod(m).Equals(MustParseFraction("1/2")) {
		t.Fatalf("5/2 mod 2 != 1/2, got %s", a.Mod(m).Show())
	}
}

func TestFractionParse(t *testing.T) {
	cases := map[string]Fraction{
		"2":   FractionFromInt(2),
		"1/2": MustParseFraction("1/2"),
		"3/4": MustParseFraction("3/4"),
		"0":   FractionFromInt(0),
	}
	for s, expected := range cases {
		f, err := ParseFraction(s)
		if err != nil || !f.Equals(expected) {
			t.Fatalf("parse %q failed: %v got %v", s, err, f.Show())
		}
	}
}

func TestFractionGcdLcm(t *testing.T) {
	a := FractionFromInt(4)
	b := FractionFromInt(6)
	g := GcdFraction(a, b)
	if !g.Equals(FractionFromInt(2)) {
		t.Fatalf("gcd 4,6 !=2 got %s", g.Show())
	}
	l := LcmFraction(a, b)
	if !l.Equals(FractionFromInt(12)) {
		t.Fatalf("lcm 4,6 !=12 got %s", l.Show())
	}
}
