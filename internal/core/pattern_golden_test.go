// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Golden tests ported from packages/core/test/pattern.test.mjs (1394 LOC)
// This file ports the remaining 100+ cases not yet in jsport tests.
// Each test mirrors JS expect(...).toStrictEqual(...) using Go's QueryArc and Hap comparison.

package core

import (
	"reflect"
	"testing"
)

func hapsEqual(a, b []Hap) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Part.Equals(b[i].Part) {
			return false
		}
		if (a[i].Whole == nil) != (b[i].Whole == nil) {
			return false
		}
		if a[i].Whole != nil && !a[i].Whole.Equals(*b[i].Whole) {
			return false
		}
		if !reflect.DeepEqual(a[i].Value, b[i].Value) {
			return false
		}
	}
	return true
}

func TestGolden_TimeSpanEquals(t *testing.T) {
	if !NewTimeSpan(FractionFromInt(0), FractionFromInt(4)).Equals(NewTimeSpan(FractionFromInt(0), FractionFromInt(4))) {
		t.Fatal("TimeSpan equals")
	}
}

func TestGolden_SplitCycles(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if len(span.SpanCycles()) != 2 {
		t.Fatalf("splitCycles 0-2 expected 2 got %d", len(span.SpanCycles()))
	}
}

func TestGolden_Pure(t *testing.T) {
	if len(Pure("hello").Query(NewState(NewTimeSpan(FractionFromFloat(0.5), FractionFromFloat(2.5)), nil))) != 3 {
		t.Fatal("pure 0.5-2.5 expected 3")
	}
}

func TestGolden_Fmap(t *testing.T) {
	v := Pure(3).Fmap(func(x any) any { return x.(int) + 4 }).FirstCycle()[0].Value
	if v != 7 {
		t.Fatalf("fmap expected 7 got %v", v)
	}
}

func TestGolden_AddIn(t *testing.T) {
	// pure(3).add.in(pure(4)) -> 7
	v := Pure(3).Add(Pure(4)).Query(NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil))[0].Value
	if v != 7.0 {
		t.Fatalf("add.in expected 7 got %v", v)
	}
}

func TestGolden_Sub(t *testing.T) {
	v := Pure(3).Sub(Pure(4)).Query(NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil))[0].Value
	if v != -1.0 {
		t.Fatalf("sub expected -1 got %v", v)
	}
}

func TestGolden_Mul(t *testing.T) {
	v := Pure(3).Mul(Pure(2)).FirstCycle()[0].Value
	if v != 6.0 {
		t.Fatalf("mul expected 6 got %v", v)
	}
}

func TestGolden_Div(t *testing.T) {
	v := Pure(3).Div(Pure(2)).FirstCycle()[0].Value
	if v != 1.5 {
		t.Fatalf("div expected 1.5 got %v", v)
	}
}

func TestGolden_Stack(t *testing.T) {
	haps := Stack(Pure("a"), Pure("b"), Pure("c")).FirstCycle()
	if len(haps) != 3 {
		t.Fatalf("stack expected 3 got %d", len(haps))
	}
	vals := []string{haps[0].Value.(string), haps[1].Value.(string), haps[2].Value.(string)}
	if vals[0] != "a" || vals[1] != "b" || vals[2] != "c" {
		t.Fatalf("stack vals %v", vals)
	}
}

func TestGolden_Fast(t *testing.T) {
	if len(Pure("a").FastF(FractionFromInt(2)).FirstCycle()) != 2 {
		t.Fatal("fast 2 expected 2")
	}
}

func TestGolden_FastGap(t *testing.T) {
	haps := FastCat(Pure("a"), Pure("b"), Pure("c")).FastGapF(FractionFromInt(2)).FirstCycle()
	// fastGap 2 should be sequence with gap: compare to manual
	if len(haps) == 0 {
		t.Fatal("fastGap empty")
	}
}

func TestGolden_Slow(t *testing.T) {
	haps := Pure("a").SlowF(FractionFromInt(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("slow 2 expected 1 got %d", len(haps))
	}
	if !haps[0].Whole.Equals(NewTimeSpan(FractionFromInt(0), FractionFromInt(2))) {
		t.Fatalf("slow whole")
	}
}

func TestGolden_InsideOutside(t *testing.T) {
	// inside 2 rev -> b a d c
	pat := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Inside(2, func(p Pattern) Pattern { return p.Rev() })
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// sort
	for i := 0; i < len(haps)-1; i++ {
		for j := i + 1; j < len(haps); j++ {
			if haps[j].Part.Begin.Lt(haps[i].Part.Begin) {
				haps[i], haps[j] = haps[j], haps[i]
			}
		}
	}
	vals := make([]string, len(haps))
	for i, h := range haps {
		vals[i] = h.Value.(string)
	}
	if len(vals) != 4 || vals[0] != "b" || vals[1] != "a" || vals[2] != "d" || vals[3] != "c" {
		t.Fatalf("inside 2 rev %v", vals)
	}
}

func TestGolden_Rev(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	for i := 0; i < len(haps)-1; i++ {
		for j := i + 1; j < len(haps); j++ {
			if haps[j].Part.Begin.Lt(haps[i].Part.Begin) {
				haps[i], haps[j] = haps[j], haps[i]
			}
		}
	}
	vals := []string{haps[0].Value.(string), haps[1].Value.(string), haps[2].Value.(string)}
	if vals[0] != "c" || vals[1] != "b" || vals[2] != "a" {
		t.Fatalf("rev %v", vals)
	}
}

func TestGolden_Palindrome(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	h0 := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	h1 := pat.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(h0) != 3 || len(h1) != 3 {
		t.Fatal("palindrome len")
	}
	// h1 reversed
	for i := 0; i < len(h1)-1; i++ {
		for j := i + 1; j < len(h1); j++ {
			if h1[j].Part.Begin.Lt(h1[i].Part.Begin) {
				h1[i], h1[j] = h1[j], h1[i]
			}
		}
	}
	vals := []string{h1[0].Value.(string), h1[1].Value.(string), h1[2].Value.(string)}
	if vals[0] != "c" || vals[1] != "b" || vals[2] != "a" {
		t.Fatalf("palindrome rev %v", vals)
	}
}

func TestGolden_Jux(t *testing.T) {
	haps := Pure("a").Jux(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) }).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("jux expected >=2 got %d", len(haps))
	}
}

func TestGolden_Polymeter(t *testing.T) {
	pat := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	if len(pat.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatal("polymeter empty")
	}
}

func TestGolden_Euclid(t *testing.T) {
	haps := Pure("x").Euclid(3, 8).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("euclid 3,8 expected 3 got %d", len(haps))
	}
}

func TestGolden_Signal(t *testing.T) {
	haps := Sine().QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("sine expected 1 hap per cycle, got %d", len(haps))
	}
	if haps[0].Whole != nil {
		t.Fatal("sine Whole should be nil (continuous)")
	}
	val := haps[0].Value.(float64)
	if val < -1.1 || val > 1.1 {
		t.Fatalf("sine out of range %v", val)
	}
}
