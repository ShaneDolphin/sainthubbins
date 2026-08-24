// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Extra: remaining free wrappers for JS `register` curry helpers that are Pattern methods in Go but also free in JS.

package core

// JS free functions that are curry wrappers around Pattern methods: add, sub, mul, div, etc are already AddF etc in pattern_free_aliases.go.
// Provide additional ones for completeness that were still missing in free2:

func BandF(a, b any) Pattern { return Reify(b).Band(a) }
func BorF(a, b any) Pattern { return Reify(b).Bor(a) }
func BxorF(a, b any) Pattern { return Reify(b).Bxor(a) }
func BlshiftF(a, b any) Pattern { return Reify(b).Blshift(a) }
func BrshiftF(a, b any) Pattern { return Reify(b).Brshift(a) }
func LtF(a, b any) Pattern { return Reify(b).Lt(a) }
func GtF(a, b any) Pattern { return Reify(b).Gt(a) }
func LteF(a, b any) Pattern { return Reify(b).Lte(a) }
func GteF(a, b any) Pattern { return Reify(b).Gte(a) }
func EqF(a, b any) Pattern { return Reify(b).Eq(a) }
func NeF(a, b any) Pattern { return Reify(b).Ne(a) }
func AndF(a, b any) Pattern { return Reify(b).And(a) }
func OrF(a, b any) Pattern { return Reify(b).Or(a) }
func KeepF(a, b any) Pattern { return Reify(b).Keep(a) }
func KeepIfF(a, b any) Pattern { return Reify(b).KeepIf(a) }
func SetF(a, b any) Pattern { return Reify(b).Set(a) }
