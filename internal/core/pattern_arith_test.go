// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package core

import (
	"reflect"
	"testing"
)

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

// TestAddOnNamedNotesIsANoOp guards against a second corruption reported by
// the coordinator: named notes ("c3", "e3", "g3" — the strings mini-notation
// produces, kept as-is until the sound engine parses them) are not numeric,
// so toFloat("c3") is 0 and the old arithmetic computed 0+12 for every one
// of them, collapsing three distinct pitches to the identical wrong value
// (note:12, note:12, note:12). Add must leave a non-numeric primary field
// untouched instead — a no-op, not a destroyer.
func TestAddOnNamedNotesIsANoOp(t *testing.T) {
	// FastCat(Pure("c3"), ...) is exactly the pattern mini.Mini("c3 e3 g3")
	// produces; internal/core cannot import internal/mini (mini imports
	// core), so the raw strings are built directly here.
	p := Note(FastCat(Pure("c3"), Pure("e3"), Pure("g3"))).Add(12)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("got %d haps, want 3", len(haps))
	}
	want := []string{"c3", "e3", "g3"}
	for i, w := range want {
		m := bagOf(t, haps[i].Value)
		if got := m["note"]; got != w {
			t.Errorf("hap %d note = %v, want %q — named notes must survive Add unchanged, not collapse to a number", i, got, w)
		}
	}
}

// TestAddBagPlusBagPreservesAbsentPrimaryField guards against a third
// corruption in the same bag+bag branch: the no-op-on-non-numeric guard
// added for TestAddOnNamedNotesIsANoOp skipped the primary field on a
// non-numeric right-hand value without checking whether the left-hand bag
// even had that key. For a bag that lacks a pitch field entirely —
// S("bd"), Gain(0.5) — primaryNumeric defaults to "note", so the guard's
// "leave it as cloned" fell through to nothing being cloned at all, and
// the note vanished instead of being copied across. Every case here
// compares the whole resulting bag, not one key, because the bug is a
// missing key that a single-key assertion would not catch.
func TestAddBagPlusBagPreservesAbsentPrimaryField(t *testing.T) {
	cases := []struct {
		name string
		p    Pattern
		want map[string]any
	}{
		{
			name: `S("bd").Add(Note("c3")) keeps both s and note`,
			p:    S("bd").Add(Note("c3")),
			want: map[string]any{"s": "bd", "note": "c3"},
		},
		{
			name: `Gain(0.5).Add(Note("c3")) keeps both gain and note`,
			p:    Gain(0.5).Add(Note("c3")),
			want: map[string]any{"gain": 0.5, "note": "c3"},
		},
		{
			name: `Note("c3").Add(Note(60)) leaves the name intact (no-op)`,
			p:    Note("c3").Add(Note(60)),
			want: map[string]any{"note": "c3"},
		},
		{
			name: `Note(60).Add(Note("c3")) keeps the number`,
			p:    Note(60).Add(Note("c3")),
			want: map[string]any{"note": 60},
		},
		{
			name: `Note(60).Add(Note(12)) sums to 72`,
			p:    Note(60).Add(Note(12)),
			want: map[string]any{"note": float64(72)},
		},
		{
			name: `S("bd").Add(S("sd")) right-hand-wins for non-primary keys`,
			p:    S("bd").Add(S("sd")),
			want: map[string]any{"s": "sd"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			haps := c.p.QueryArc(FractionFromInt(0), FractionFromInt(1))
			if len(haps) != 1 {
				t.Fatalf("got %d haps, want 1", len(haps))
			}
			m := bagOf(t, haps[0].Value)
			if !reflect.DeepEqual(m, c.want) {
				t.Errorf("bag = %#v, want %#v", m, c.want)
			}
		})
	}
}
