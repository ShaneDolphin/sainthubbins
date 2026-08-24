package core

import "testing"

func TestMJS_ControlsValue2(t *testing.T) {
	p := S("bd")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("S bd")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["s"] != "bd" {
		t.Fatalf("S bd map")
	}
	// S with n
	p2 := Stack(S("bd"), Pure(map[string]any{"n": 2}))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("Stack S bd + n2 expected 2")
	}
}

func TestMJS_UtilsMod2(t *testing.T) {
	if Mod(5, 3) != 2 {
		t.Fatalf("Mod 5,3=2")
	}
	if Mod(-1, 3) != 2 {
		t.Fatalf("Mod -1,3=2")
	}
}

func TestMJS_LogValues2(t *testing.T) {
	p := Pure("a").Fmap(func(v any) any { return v.(string) + "b" })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 || haps[0].Value != "ab" {
		t.Fatalf("Fmap ab")
	}
	st := Stack(Pure("a"), Pure("b"))
	if len(st.QueryArc(FractionFromInt(0), FractionFromInt(1))) != 2 {
		t.Fatalf("Stack 2")
	}
}
