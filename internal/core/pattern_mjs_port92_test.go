package core

import "testing"

func TestMJS_PickModSqueeze2(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b"})
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose")
	}
	// SqueezeJoin
	sq := Pure(Pure("a")).SqueezeJoin()
	if len(sq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SqueezeJoin")
	}
	ij := Pure(Pure("a")).InnerJoin()
	if len(ij.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("InnerJoin")
	}
}

func TestMJS_CatSequence2(t *testing.T) {
	c := Cat(Pure("a"), Pure("b"), Pure("c"))
	if len(c.QueryArc(FractionFromInt(0), FractionFromInt(3))) == 0 {
		t.Fatalf("Cat 3 cycles")
	}
	seq := Sequence(Pure("a"), Pure("b"))
	if len(seq.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Sequence 2")
	}
}

func TestMJS_SetWithValue2(t *testing.T) {
	p := Pure(map[string]any{"a": 1})
	p2 := p.Set(map[string]any{"b": 2})
	haps := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Set a1 b2")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["a"] != 1 || m["b"] != 2 {
		t.Fatalf("Set merge got %v", haps[0].Value)
	}
	wv := Pure(1).WithValue(func(v any) any { return v.(int) * 2 })
	haps2 := wv.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2)==0 || haps2[0].Value.(int)!=2 {
		t.Fatalf("WithValue *2")
	}
}
