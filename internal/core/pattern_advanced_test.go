package core

import "testing"

func TestPatternWithValueAdvanced(t *testing.T) {
	p := Pure(2).Fmap(func(v any) any { return toFloat(v) + 3 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if toFloat(haps[0].Value) != 5 {
		t.Fatalf("WithValueAdvanced")
	}
}

func TestPatternOffWhen(t *testing.T) {
	p := Pure("a")
	off := p.Off(MustParseFraction("1/8"), func(pp Pattern) Pattern { return pp.Fmap(func(v any) any { return "b" }) })
	haps := off.QueryArc(FractionFromInt(0), FractionFromInt(1))
	// Stack original + off should give 2?
	if len(haps) == 0 {
		t.Fatalf("Off no haps")
	}
	t.Logf("Off haps %d", len(haps))
}

func TestPatternPly(t *testing.T) {
	p := Pure("a").Ply(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Ply 2 expected 2 got %d", len(haps))
	}
}
