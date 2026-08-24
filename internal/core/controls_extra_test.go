package core

import "testing"

func TestControlsS2(t *testing.T) {
	p := S("bd")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("S bd expected 1 got %d", len(haps))
	}
}

func TestControlsSN(t *testing.T) {
	p := S([]any{"bd", 2})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("S with n expected 1 got %d", len(haps))
	}
	if m, ok := haps[0].Value.(map[string]any); ok {
		if m["n"] == nil {
			t.Fatalf("n missing in S bd 2")
		}
	}
}

func TestControlsCutoff2(t *testing.T) {
	p := Cutoff(500)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Cutoff expected 1")
	}
}

func TestControlsStack2(t *testing.T) {
	p := Stack(S("bd"), S("sd"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Stack S bd/sd expected 2 got %d", len(haps))
	}
}
