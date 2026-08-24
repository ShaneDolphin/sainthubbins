// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pick.mjs — pick pattern helpers.

package core

import (
	"fmt"
	"math"
)

// Pick picks patterns/values from lookup by index or key (innerJoin).
func Pick(lookup any, pat Pattern) Pattern {
	// Handle lookup as []any or map
	switch lk := lookup.(type) {
	case []any:
		return pickFromSlice(lk, pat, false)
	case []string:
		anyList := make([]any, len(lk))
		for i, v := range lk {
			anyList[i] = v
		}
		return pickFromSlice(anyList, pat, false)
	case map[string]any:
		return pickFromMap(lk, pat, false)
	default:
		// try to handle map[string]string etc via any
		return pat.Fmap(func(v any) any { return lookup }).InnerJoin()
	}
}

func PickMod(lookup any, pat Pattern) Pattern {
	switch lk := lookup.(type) {
	case []any:
		return pickFromSlice(lk, pat, true)
	case []string:
		anyList := make([]any, len(lk))
		for i, v := range lk {
			anyList[i] = v
		}
		return pickFromSlice(anyList, pat, true)
	default:
		return Pick(lookup, pat)
	}
}

func pickFromSlice(list []any, pat Pattern, mod bool) Pattern {
	if len(list) == 0 {
		return Silence()
	}
	// Reify list elements? For now keep as is
	return pat.Fmap(func(v any) any {
		var idx int
		switch x := v.(type) {
		case int:
			idx = x
		case float64:
			idx = int(mathRound(x))
		case Fraction:
			idx = int(mathRound(x.Float64()))
		default:
			idx = 0
		}
		if mod {
			idx = Mod(idx, len(list))
		} else {
			if idx < 0 {
				idx = 0
			}
			if idx >= len(list) {
				idx = len(list) - 1
			}
		}
		return list[idx]
	}).InnerJoin()
}

func pickFromMap(m map[string]any, pat Pattern, mod bool) Pattern {
	if len(m) == 0 {
		return Silence()
	}
	return pat.Fmap(func(v any) any {
		key := fmt.Sprintf("%v", v)
		if val, ok := m[key]; ok {
			return val
		}
		return nil
	}).InnerJoin()
}

func mathRound(f float64) float64 { return math.Round(f) }
