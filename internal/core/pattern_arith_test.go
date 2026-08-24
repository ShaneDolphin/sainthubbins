// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package core

import "testing"

func bagOf(t *testing.T, v any) map[string]any {
	t.Helper()
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("value is %T (%v), want a control bag", v, v)
	}
	return m
}

func TestAddTransposesANotePattern(t *testing.T) {
	p := Note(FastCat(Pure(0), Pure(4), Pure(7))).Add(12)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("got %d haps, want 3", len(haps))
	}
	want := []float64{12, 16, 19}
	for i, w := range want {
		m := bagOf(t, haps[i].Value)
		if got := toFloat(m["note"]); got != w {
			t.Errorf("hap %d note = %v, want %v", i, got, w)
		}
	}
}

func TestAddKeepsOtherControls(t *testing.T) {
	p := Note(Pure(60)).Set(Gain(0.5)).Add(12)
	m := bagOf(t, p.QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value)
	if toFloat(m["note"]) != 72 {
		t.Errorf("note = %v, want 72", m["note"])
	}
	if toFloat(m["gain"]) != 0.5 {
		t.Errorf("gain = %v, want 0.5 — Add must not disturb other controls", m["gain"])
	}
}

func TestAddOnBareNumbersStillWorks(t *testing.T) {
	haps := Pure(60).Add(12).QueryArc(FractionFromInt(0), FractionFromInt(1))
	if got := toFloat(haps[0].Value); got != 72 {
		t.Errorf("got %v, want 72", got)
	}
}

// TestAddBagPlusBagPreservesNonNumericFields guards against a corruption
// found while implementing this change: bag+bag addition summed *every*
// matching key via toFloat, including non-numeric ones. Two bags that both
// carry a sound name ("s") collided into s:0, silently destroying it. Only
// genuinely numeric matching fields should sum; a non-numeric collision
// should keep the right-hand (b) value, matching Set's override semantics.
func TestAddBagPlusBagPreservesNonNumericFields(t *testing.T) {
	a := Note(Pure(60)).Set(S(Pure("bd")))
	b := Note(Pure(3)).Set(S(Pure("sd")))
	m := bagOf(t, a.Add(b).QueryArc(FractionFromInt(0), FractionFromInt(1))[0].Value)
	if toFloat(m["note"]) != 63 {
		t.Errorf("note = %v, want 63 (60+3)", m["note"])
	}
	if m["s"] != "sd" {
		t.Errorf("s = %v, want %q — non-numeric collisions must not become 0", m["s"], "sd")
	}
}
