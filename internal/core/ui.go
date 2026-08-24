// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

// UI helpers stub (slider, etc.)
func Slider(name string, min, max, value float64) Pattern {
	return Pure(map[string]any{"slider": name, "min": min, "max": max, "value": value})
}
