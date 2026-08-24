// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

// Speak is speech synthesis stub
func Speak(text string) Pattern {
	return Pure(map[string]any{"speak": text})
}
