package core

import "testing"

func TestMJS_FilterEvents2(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	f := p.Filter(func(h Hap) bool { return h.Value == "a" })
	haps := f.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value != "a" {
		t.Fatalf("Filter a expected 1 a got %v", haps)
	}
}

func TestMJS_OnsetsOnly2(t *testing.T) {
	p := Pure("a")
	o := p.OnsetsOnly()
	haps := o.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("OnsetsOnly expected non-empty")
	}
	for _, h := range haps {
		if !h.HasOnset() {
			t.Fatalf("OnsetsOnly has onset false")
		}
	}
}

func TestMJS_SpanConversions2(t *testing.T) {
	span := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	if span.Duration().Float64() != 1 {
		t.Fatalf("Duration 1")
	}
	span2 := NewTimeSpan(FractionFromInt(0), FractionFromFloat(0.5))
	if span2.Duration().Float64() != 0.5 {
		t.Fatalf("Duration 0.5")
	}
}
