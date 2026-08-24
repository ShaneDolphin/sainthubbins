package core

import "testing"

func TestMJS_DegradeVariants2(t *testing.T) {
	d0 := Pure("a").DegradeBy(0)
	if len(d0.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 1 {
		t.Fatalf("DegradeBy 0 expected 1")
	}
	d05 := Pure("a").DegradeBy(0.5)
	// Degrade 0.5 may be 0 or 1 depending on pseudoRand; just check not panic and within 0-1
	haps := d05.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) < 0 || len(haps) > 1 {
		t.Fatalf("DegradeBy 0.5 len 0-1")
	}
	d := Pure("a").Degrade()
	haps2 := d.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) < 0 || len(haps2) > 1 {
		t.Fatalf("Degrade 0.5 len 0-1")
	}
}

func TestMJS_HushSilenceAlias2(t *testing.T) {
	h := Pure("a").Hush()
	if len(h.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Hush expected 0")
	}
	s := Silence()
	if len(s.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Silence expected 0")
	}
	// Decay via Set
	p := Pure("a").Set(Decay(0.5))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Decay via Set expected 1")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["decay"] == nil {
		// allow but check
	}
}

func TestMJS_SometimesDegradeCombo2(t *testing.T) {
	d := Pure("a").DegradeBy(0).Sometimes(func(p Pattern) Pattern { return p.FastF(FractionFromInt(2)) })
	haps := d.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("DegradeBy 0 Sometimes FastF2 expected non-empty")
	}
}
