package core

import "testing"

func TestControlsS(t *testing.T) {
	h := S("bd").FirstCycle()[0].Value
	m, ok := h.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Fatalf("S bd failed %v", h)
	}
	h2 := S([]any{"bd", 2, 0.5}).FirstCycle()[0].Value
	m2, _ := h2.(map[string]any)
	if m2["s"] != "bd" || m2["n"] != 2 || m2["gain"] != 0.5 {
		t.Fatalf("S multi failed %v", m2)
	}
}

func TestControlsAlias(t *testing.T) {
	h1 := Cutoff(500).FirstCycle()[0].Value
	h2 := Lpf(500).FirstCycle()[0].Value
	m1, _ := h1.(map[string]any)
	m2, _ := h2.(map[string]any)
	if m1["cutoff"] != m2["cutoff"] {
		t.Fatalf("Cutoff/Lpf alias mismatch")
	}
}

func TestControlsGenerated(t *testing.T) {
	// Check a few generated controls exist
	if Pan == nil {
		t.Fatalf("Pan nil")
	}
	if Gain == nil {
		t.Fatalf("Gain nil")
	}
	h := Pan(0.5).FirstCycle()[0].Value
	if m, ok := h.(map[string]any); !ok || m["pan"] != 0.5 {
		t.Fatalf("Pan 0.5 failed")
	}
}
