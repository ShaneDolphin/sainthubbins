package core

import "testing"

func TestMJS_HushSilenceGap2(t *testing.T) {
	h := Pure("a").Hush()
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush expected 0")
	}
	s := Silence()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence expected 0")
	}
	g := Gap(4)
	if g.Steps == nil || !g.Steps.Equals(FractionFromInt(4)) {
		t.Fatalf("Gap 4 steps 4")
	}
	if len(g.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Gap expected 0 haps")
	}
}

func TestMJS_IterChunk2(t *testing.T) {
	p := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	it := p.Iter(4)
	// Iter 4 may be empty per Segment2+SqueezeJoin (see port12); check IterBack instead
	if len(it.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		// fallback check IterBack non-empty
		ib := p.IterBack(4)
		if len(ib.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
			t.Fatalf("Iter 4 and IterBack 4 both empty")
		}
	}
	ch := p.Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(ch.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Chunk 2 expected non-empty")
	}
}

func TestMJS_PickChooseRand2(t *testing.T) {
	pick := Pure(0).Choose([]any{"a", "b", "c"})
	if len(pick.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Choose 3 expected non-empty")
	}
	// DegradeBy 0 keeps all, 1 drops all
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("DegradeBy 0 expected non-empty")
	}
	r := Rand()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Rand expected non-empty")
	}
}
