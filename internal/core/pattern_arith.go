// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Arithmetic that understands control bags.

package core

import "fmt"

// numericControls are the fields arithmetic targets, in priority order. A bag
// carries several values and only the pitch-like one should move when a
// pattern is transposed.
var numericControls = []string{"note", "n", "freq"}

// addValues adds b into a. Bare numbers add directly. When either side is a
// control bag, the addition lands on the bag's numeric field and every other
// control is preserved — transposing a note must not discard its gain.
func addValues(a, b any) any {
	am, aIsBag := a.(map[string]any)
	bm, bIsBag := b.(map[string]any)

	switch {
	case !aIsBag && !bIsBag:
		return toFloat(a) + toFloat(b)

	case aIsBag && !bIsBag:
		out := cloneBag(am)
		addIntoField(out, primaryNumeric(am), toFloat(b))
		return out

	case !aIsBag && bIsBag:
		out := cloneBag(bm)
		addIntoField(out, primaryNumeric(bm), toFloat(a))
		return out

	default:
		out := cloneBag(am)
		// The primary field (note/n/freq) gets the same no-op-on-non-numeric
		// protection as the bag+number direction, routed through
		// addIntoField, instead of the generic right-hand-wins merge below.
		// Right-hand-wins is correct for a non-primary collision (it is what
		// protects s:"bd" vs s:"sd"), but on the primary field it would do
		// exactly what addIntoField exists to prevent: Note("c3").Add(Note(60))
		// must leave "c3" untouched, not overwrite it with 60, and
		// Note(60).Add(Note("c3")) must leave 60 untouched rather than
		// destroy it with a string that has nothing numeric to add.
		primary := primaryNumeric(am)
		for k, v := range bm {
			if k == primary {
				if vf, ok := numericValue(v); ok {
					addIntoField(out, k, vf)
				}
				// else: the right-hand side's primary field isn't genuinely
				// numeric (e.g. a named note) — there's nothing sensible to
				// add, so the left-hand value already cloned into out is
				// left exactly as it was.
				continue
			}
			existing, ok := out[k]
			if !ok {
				out[k] = v
				continue
			}
			ef, eok := numericValue(existing)
			vf, vok := numericValue(v)
			if eok && vok {
				out[k] = ef + vf
			} else {
				// A matching key that isn't genuinely numeric on both sides
				// (e.g. two sound names both under "s") must not be summed —
				// toFloat on a non-numeric string is 0, which would silently
				// wipe it out. Let the right-hand value win instead, matching
				// the override semantics Set already uses for control bags.
				out[k] = v
			}
		}
		return out
	}
}

// addIntoField adds delta into bag[key] in place. If the field is absent,
// it is created with value delta (the bag had no pitch yet — Add gives it
// one). If the field is present but not genuinely numeric — a named note
// such as "c3", kept as a string until the sound engine parses it — the
// field is left untouched rather than summed: toFloat on a non-numeric
// string is 0, and computing 0+delta would silently collapse every distinct
// named note in a pattern to the same wrong value. A no-op is the safer
// failure; see docs/tutorial/05-transformations.md.
func addIntoField(bag map[string]any, key string, delta float64) {
	cur, present := bag[key]
	if !present {
		bag[key] = delta
		return
	}
	if nf, ok := numericValue(cur); ok {
		bag[key] = nf + delta
	}
	// else: non-numeric current value — leave it exactly as cloned.
}

// numericValue reports whether v is a value toFloat treats as genuinely
// numeric (as opposed to a non-numeric string that toFloat would silently
// coerce to 0).
func numericValue(v any) (float64, bool) {
	switch x := v.(type) {
	case float64, float32, int, int64, int32, uint, uint64, Fraction, *Fraction:
		return toFloat(x), true
	case string:
		var f float64
		if _, err := fmt.Sscanf(x, "%f", &f); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// primaryNumeric picks the field arithmetic should target, defaulting to note
// so that adding to a bag without one still produces something musical.
func primaryNumeric(m map[string]any) string {
	for _, k := range numericControls {
		if _, ok := m[k]; ok {
			return k
		}
	}
	return "note"
}

func cloneBag(m map[string]any) map[string]any {
	out := make(map[string]any, len(m)+1)
	for k, v := range m {
		out[k] = v
	}
	return out
}
