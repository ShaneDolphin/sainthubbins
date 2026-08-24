package core

import "testing"

func TestSpanCyclesBasic(t *testing.T) {
	ts := NewTimeSpan(FractionFromInt(0), FractionFromInt(2))
	parts := ts.SpanCycles()
	if len(parts) != 2 {
		t.Fatalf("SpanCycles 0-2 expected 2 got %d", len(parts))
	}
	if !parts[0].Begin.Equals(FractionFromInt(0)) || !parts[0].End.Equals(FractionFromInt(1)) {
		t.Fatalf("SpanCycles first expected 0-1 got %v", parts[0])
	}
}

func TestHapHasOnsetExtra(t *testing.T) {
	whole := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	part := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	h := Hap{Whole: &whole, Part: part, Value: "a"}
	if !h.HasOnset() {
		t.Fatalf("HasOnset expected true")
	}
	part2 := NewTimeSpan(FractionFromFloat(0.5), FractionFromInt(1))
	h2 := Hap{Whole: &whole, Part: part2, Value: "a"}
	if h2.HasOnset() {
		t.Fatalf("HasOnset expected false")
	}
}

func TestHapWholeOrPartExtra(t *testing.T) {
	tsWhole := NewTimeSpan(FractionFromInt(0), FractionFromInt(1))
	tsPart := NewTimeSpan(FractionFromFloat(0.25), FractionFromFloat(0.75))
	h := Hap{Whole: &tsWhole, Part: tsPart, Value: "a"}
	if !h.WholeOrPart().Begin.Equals(tsWhole.Begin) {
		t.Fatalf("WholeOrPart expected Whole")
	}
	h2 := Hap{Whole: nil, Part: tsPart, Value: "a"}
	if !h2.WholeOrPart().Begin.Equals(tsPart.Begin) {
		t.Fatalf("WholeOrPart expected Part")
	}
}

func TestFractionOpsExtra(t *testing.T) {
	a := FractionFromInt(1).Div(FractionFromInt(3))
	b := FractionFromInt(2).Div(FractionFromInt(3))
	sum := a.Add(b)
	if !sum.Equals(FractionFromInt(1)) {
		t.Fatalf("1/3+2/3 expected 1 got %v", sum)
	}
	prod := a.Mul(FractionFromInt(3))
	if !prod.Equals(FractionFromInt(1)) {
		t.Fatalf("1/3*3 expected 1 got %v", prod)
	}
}
