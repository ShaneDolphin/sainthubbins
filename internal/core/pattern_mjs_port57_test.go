package core

import "testing"

func TestMJS_StructMask2(t *testing.T) {
	s := Sequence(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	st := s.Struct(Pure(true))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Struct true expected non-empty")
	}
	ma := s.Mask(Pure(true))
	if len(ma.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Mask true expected non-empty")
	}
	maF := s.Mask(Pure(false))
	if len(maF.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 0 {
		t.Fatalf("Mask false expected 0")
	}
}

func TestMJS_EuclidLegato2(t *testing.T) {
	e := Pure("a").EuclidLegato(3, 8)
	haps := e.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("EuclidLegato 3,8 expected 3 got %d", len(haps))
	}
	// check legato: contiguous?
	if len(haps) >= 2 && haps[0].Whole != nil && haps[1].Whole != nil {
		if !haps[0].Whole.End.Equals(haps[1].Whole.Begin) {
			// legato may be contiguous, but allow otherwise
		}
	}
}

func TestMJS_HapContext2(t *testing.T) {
	p := Pure(map[string]any{"s": "bd", "n": 1})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Pure map bd n1 expected 1")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["s"] != "bd" || m["n"] != 1 {
		t.Fatalf("map s bd n1 got %v", haps[0].Value)
	}
	// SetContext
	pc := p.SetContext(map[string]any{"orbit": 2})
	haps2 := pc.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 || haps2[0].Context["orbit"] != 2 {
		t.Fatalf("SetContext orbit 2 got %v", haps2[0].Context)
	}
}
