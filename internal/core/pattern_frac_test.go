package core

import "testing"

func TestFractionShow(t *testing.T) {
	f := FractionFromInt(3).Div(FractionFromInt(4))
	if f.Show() != "3/4" {
		t.Fatalf("Show expected 3/4 got %s", f.Show())
	}
	if f.Float64() != 0.75 {
		t.Fatalf("Float64 expected 0.75 got %v", f.Float64())
	}
}

func TestFractionSam(t *testing.T) {
	f := FractionFromFloat(2.7)
	sam := f.Sam()
	if !sam.Equals(FractionFromInt(2)) {
		t.Fatalf("Sam 2.7 expected 2 got %v", sam)
	}
	pos := f.CyclePos().Float64()
	if pos < 0.69 || pos > 0.71 {
		t.Fatalf("CyclePos 2.7 expected ~0.7 got %v", pos)
	}
}

func TestFractionNextSam(t *testing.T) {
	f := FractionFromInt(1)
	next := f.NextSam()
	if !next.Equals(FractionFromInt(2)) {
		t.Fatalf("NextSam 1 expected 2 got %v", next)
	}
}

func TestFractionComparison(t *testing.T) {
	a := FractionFromInt(1).Div(FractionFromInt(2))
	b := FractionFromInt(2).Div(FractionFromInt(4))
	if !a.Equals(b) {
		t.Fatalf("1/2 equals 2/4 failed")
	}
	if !a.Lt(FractionFromInt(1)) {
		t.Fatalf("1/2 < 1 failed")
	}
	if !a.Gt(FractionFromInt(0)) {
		t.Fatalf("1/2 > 0 failed")
	}
}
