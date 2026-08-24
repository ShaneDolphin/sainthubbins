package core

import "testing"

func TestMJS_SqueezeJoinInnerJoin2(t *testing.T) {
	p := Pure(Pure("a")).SqueezeJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("SqueezeJoin expected non-empty")
	}
	ij := Pure(Pure("a")).InnerJoin()
	if len(ij.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("InnerJoin expected non-empty")
	}
}

func TestMJS_OutSqueeze2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd"})
	// Set via OpOut alias
	out := p.Set(map[string]any{"gain": 0.5})
	haps := out.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Set gain 0.5")
	}
}

func TestMJS_PureFmapLog2(t *testing.T) {
	p := Pure(3).Fmap(func(v any) any { return v.(int) + 1 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value.(int) != 4 {
		t.Fatalf("Fmap 3+1=4 got %v", haps[0].Value)
	}
}
