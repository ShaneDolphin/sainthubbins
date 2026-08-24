package core

import "testing"

func TestMJS_Port594_JuxSuperimposeLayerPalFourth(t *testing.T) {
	p := Pure("bd").Jux(func(q Pattern) Pattern { return q.Rev() })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 { t.Fatalf("Jux") }
	q := Pure("a").Superimpose(func(x Pattern) Pattern { return x.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) < 2 { t.Fatalf("Superimpose") }
	r := Stack(Pure("a"), Pure("b"))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Stack") }
	s := FastCat(Pure("a"), Pure("b"), Pure("c")).Palindrome()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Pal") }
}
func TestMJS_Port594_IterChunkSqueezeJoinFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).IterBack(2)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("IterBack") }
	q := Pure("a b c d").Chunk(2, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("Chunk") }
	r := Pure(Pure("a")).SqueezeJoin()
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("SqueezeJoin") }
	inner := Pure(Pure("x")).InnerJoin()
	if len(inner.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 { t.Fatalf("InnerJoin") }
}
func TestMJS_Port594_FilterMaskStructFourth(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).FilterValues(func(v any) bool { return v != "b" })
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Filter") }
	q := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d")).Struct(FastCat(Pure(true), Pure(true), Pure(false), Pure(false)))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Struct") }
	r := FastCat(Pure("a"), Pure("b")).Mask(Pure(true))
	if len(r.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 { t.Fatalf("Mask") }
}
