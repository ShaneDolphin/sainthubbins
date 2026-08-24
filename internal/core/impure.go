// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

// Impure helpers for side-effectful patterns (console, etc.)
func LogPat(p Pattern) Pattern {
	return p.Fmap(func(v any) any {
		Logger(v)
		return v
	})
}
