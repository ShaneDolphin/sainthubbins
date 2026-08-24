// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package tonal

import (
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// ScalePattern returns pattern of scale notes (mirrors JS scale register)
func ScalePattern(scaleName string) core.Pattern {
	notes := Scale(scaleName)
	pats := make([]core.Pattern, len(notes))
	for i, n := range notes {
		pats[i] = core.Pure(n)
	}
	if len(pats) == 0 {
		return core.Silence()
	}
	return core.FastCat(pats...)
}

// ChordPattern returns pattern of chord notes
func ChordPattern(chordName string) core.Pattern {
	notes := Chord(chordName)
	pats := make([]core.Pattern, len(notes))
	for i, n := range notes {
		pats[i] = core.Pure(n)
	}
	if len(pats) == 0 {
		return core.Silence()
	}
	return core.FastCat(pats...)
}

// TransposePattern transposes pattern notes
func TransposePattern(pat core.Pattern, interval any) core.Pattern {
	return pat.Fmap(func(v any) any {
		var note string
		switch x := v.(type) {
		case string:
			note = x
		case map[string]any:
			if n, ok := x["note"]; ok {
				if s, ok := n.(string); ok {
					note = s
				}
			} else if n, ok := x["n"]; ok {
				if s, ok := n.(string); ok {
					note = s
				}
			}
		}
		if note == "" {
			return v
		}
		transposed := Transpose(note, interval)
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["note"] = transposed
			return m2
		}
		return transposed
	})
}
