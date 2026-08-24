package core

import "testing"

func TestMJS_TimeSpanHap3(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	if span.Duration().Float64() != 1 {
		t.Fatalf("Duration 1")
	}
	cycles := span.SpanCycles()
	if len(cycles) != 1 {
		t.Fatalf("SpanCycles 1")
	}
	h := NewHap(&span, span, "a", nil)
	if h.Value != "a" {
		t.Fatalf("Hap a")
	}
}

func TestMJS_PatternPureFmap3(t *testing.T) {
	p := Pure(3)
	q := p.Fmap(func(v any) any { switch x:=v.(type){case int: return x+4; case float64: return x+4; default: return v} })
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 || func()bool{v:=haps[0].Value; switch x:=v.(type){case int: return x!=7; case float64: return x!=7; default: return true}}() {
		t.Fatalf("Fmap 3+4=7")
	}
}

func TestMJS_AddWithStructure4(t *testing.T) {
	p := FastCat(Pure(1), Pure(2))
	added := p.Add(Pure(10))
	haps := added.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("FastCat Add 10 expected 2 got %d", len(haps))
	}
	if func()bool{ v0:=haps[0].Value; v1:=haps[1].Value; f0:=0.0; f1:=0.0; switch x:=v0.(type){case int: f0=float64(x); case float64: f0=x}; switch x:=v1.(type){case int: f1=float64(x); case float64: f1=x}; return f0!=11 || f1!=12 }() {
		t.Fatalf("Add 11,12 got %v %v", haps[0].Value, haps[1].Value)
	}
}
