package core

import "testing"

func TestStepwiseAliases(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"))
	f := FractionFromInt(2)
	p.Steps = &f
	if len(SAdd(1, p).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SAdd")
	}
	if len(SSub(1, p).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SSub")
	}
	if len(SExpand(2, p).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SExpand")
	}
	if len(SExtend(2, p).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SExtend")
	}
	if len(SContract(2, p).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SContract")
	}
	if len(p.Gap(2).QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Gap expected silence")
	}
}
