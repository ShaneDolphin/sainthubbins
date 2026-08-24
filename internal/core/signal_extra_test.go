package core

import "testing"

func TestSignalSine(t *testing.T) {
	p := Sine()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Sine expected >=1 hap")
	}
	if _, ok := haps[0].Value.(float64); !ok {
		t.Fatalf("Sine value expected float64 got %T", haps[0].Value)
	}
}

func TestSignalTriExtra(t *testing.T) {
	p := Tri()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Tri expected >=1")
	}
}

func TestSignalSquareExtra(t *testing.T) {
	p := Square()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Square expected >=1")
	}
}

func TestSignalRandExtra(t *testing.T) {
	p := Rand()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Rand expected >=1")
	}
}

func TestSignalRange(t *testing.T) {
	p := Sine().Range(0, 1)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Range expected >=1")
	}
	val := haps[0].Value.(float64)
	if val < 0 || val > 1 {
		t.Fatalf("Range Sine 0-1 expected 0..1 got %v", val)
	}
}
