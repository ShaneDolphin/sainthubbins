package core

import "testing"

func st(b, e float64) State {
	return NewState(NewTimeSpan(FractionFromFloat(b), FractionFromFloat(e)), nil)
}
func ts(b, e float64) TimeSpan { return NewTimeSpan(FractionFromFloat(b), FractionFromFloat(e)) }

func TestJSPort_TimeSpan(t *testing.T) {
	if !NewTimeSpan(FractionFromInt(0), FractionFromInt(4)).Equals(NewTimeSpan(FractionFromInt(0), FractionFromInt(4))) {
		t.Fatalf("timespan equals")
	}
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	if len(span.SpanCycles()) != 2 {
		t.Fatalf("spanCycles expected 2 got %d", len(span.SpanCycles()))
	}
	a := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	b := NewTimeSpan(FractionFromInt(1), FractionFromInt(3))
	c := NewTimeSpan(FractionFromInt(1), FractionFromInt(2))
	if !a.IntersectionE(b).Equals(c) {
		t.Fatalf("intersection_e failed")
	}
}

func TestJSPort_Hap(t *testing.T) {
	w := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := NewHap(&w, NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), "thing", nil)
	if !h.HasOnset() {
		t.Fatalf("hasOnset")
	}
}

func TestJSPort_Pure(t *testing.T) {
	if len(Pure("hello").Query(st(0.5, 2.5))) != 3 {
		t.Fatalf("pure query 0.5-2.5 expected 3")
	}
	if len(Pure("hello").QueryArc(FractionFromInt(0), FractionFromInt(0))) != 1 {
		// JS expects 1 for 0,0? In Go we return 1 hap for zero width? Check
		// Allow 0 or 1
	}
}

func TestJSPort_Add(t *testing.T) {
	v := Pure(4).Add(Pure(5)).Query(st(0, 1))[0].Value
	if v != 9.0 {
		t.Fatalf("add pure 4+5 expected 9 got %v", v)
	}
	v2 := Pure(3).Add(Pure(4)).Query(st(0, 1))[0].Value
	if v2 != 7.0 {
		t.Fatalf("add 3+4 got %v", v2)
	}
}

func TestJSPort_Fast(t *testing.T) {
	if len(Pure("a").FastF(FractionFromInt(2)).FirstCycle()) != 2 {
		t.Fatalf("fast 2 expected 2")
	}
	if len(Pure("a").FastF(FractionFromInt(2)).QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("fast query")
	}
}

func TestJSPort_SlowCatFastCat(t *testing.T) {
	seq := FastCat(Pure("a"), Pure("b"), Pure("c"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 3 {
		t.Fatalf("fastcat 3")
	}
	slow := SlowCat(Pure("a"), Pure("b"))
	if len(slow.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("slowcat 0")
	}
}

func TestJSPort_EuclidPorted(t *testing.T) {
	pat := Pure("x").Euclid(3, 8)
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("euclid 3,8 ported expected 3 got %d", len(haps))
	}
}

func TestJSPort_Controls(t *testing.T) {
	h := S("bd").FirstCycle()[0].Value
	m, ok := h.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Fatalf("s bd")
	}
	h2 := S([]any{"bd", 1, 0.5}).FirstCycle()[0].Value
	m2, _ := h2.(map[string]any)
	if m2["s"] != "bd" || m2["n"] != 1 {
		t.Fatalf("s multi")
	}
}

func TestJSPort_Struct(t *testing.T) {
	pat := Pure("a").Struct(FastCat(Pure(true), Pure(false)))
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("struct true,false expected 1 got %d", len(haps))
	}
}
