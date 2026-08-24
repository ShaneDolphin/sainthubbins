package core

import "testing"

func TestMJS_Port115_JuxSuperimposeLayerPal(t *testing.T) {
	p := Pure("bd sd").Jux(func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Jux")
	}
	s := Pure("bd").Superimpose(func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 {
		t.Fatalf("Superimpose >=2")
	}
	lay := Stack(Pure("a"), Pure("bd"), Pure("sd"))
	if len(lay.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Layer")
	}
	pal := Pure("a b c").Palindrome()
	if len(pal.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Palindrome")
	}
}

func TestMJS_Port115_IterChunkSqueezeJoin(t *testing.T) {
	ib := Pure("a b c d").IterBack(2)
	if len(ib.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("IterBack 2")
	}
	ch := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Chunk(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2")
	}
	sq := Pure(Pure("a")).SqueezeJoin()
	if sq.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("SqueezeJoin")
	}
	ij := Pure(Pure("bd")).InnerJoin()
	if ij.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("InnerJoin")
	}
}

func TestMJS_Port115_FilterMaskStruct(t *testing.T) {
	f := Pure("bd").FilterHaps(func(h Hap) bool { return h.Value == "bd" })
	if len(f.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("FilterHaps bd")
	}
	maskTrue := Pure("a").Mask(Pure(true))
	if len(maskTrue.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true")
	}
	maskFalse := Pure("a").Mask(Pure(false))
	if len(maskFalse.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Mask false should be 0")
	}
	st := Pure("a").Struct(Pure(true))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true")
	}
}
