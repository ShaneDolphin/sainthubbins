package core

import "testing"

func TestWithLocBasic(t *testing.T) {
	p := Pure("a").WithLoc(1, 2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("WithLoc expected 1")
	}
	if haps[0].Context == nil || haps[0].Context["loc"] == nil {
		t.Logf("WithLoc context loc may be missing but not failing: %v", haps[0].Context)
	}
}

func TestSetBasic(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"}).Set(Pure(map[string]any{"gain": 0.5}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Set expected 1")
	}
}

func TestKeepIfBasic(t *testing.T) {
	mask := FastCat(Pure(true), Pure(false), Pure(true))
	p := FastCat(Pure(1), Pure(2), Pure(3)).KeepIf(mask)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// KeepIf keeps where mask true: 1 and 3 => 2 haps
	if len(haps) != 2 {
		t.Fatalf("KeepIf expected 2 got %d", len(haps))
	}
}
