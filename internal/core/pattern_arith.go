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
		key := primaryNumeric(am)
		out[key] = toFloat(am[key]) + toFloat(b)
		return out

	case !aIsBag && bIsBag:
		out := cloneBag(bm)
		key := primaryNumeric(bm)
		out[key] = toFloat(a) + toFloat(bm[key])
		return out

	default:
		out := cloneBag(am)
		for k, v := range bm {
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
