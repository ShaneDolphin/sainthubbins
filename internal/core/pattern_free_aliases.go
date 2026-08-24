// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Free-function aliases for JS top-level exports (JS is lowercase, Go is capitalized)
// This file provides lowercase-mirroring free functions where Go already has method/capitalized free.
// These are thin wrappers so JS tunes ported via goja can find symbols.

package core

// Core free functions (JS: pure, silence, stack, cat, fastcat, slowcat, arrange, etc)
// Already have Pure, Silence, Stack, Cat, FastCat, SlowCat, Arrange, Sequence, Polymeter etc as capitalized.
// Provide lowercase aliases as exported funcs with same name lowercase? Go can't export lowercase, but we provide capitalized wrappers that match JS naming via map for goja.
// For now, provide capitalized aliases that JS expects via Evaluate's scope map: add entries there.
// This file just ensures free functions exist for direct Go calls with JS-style names via wrapper map.

// JS `pure` -> Pure, `silence` -> Silence, etc are already available capitalized.
// Provide additional free wrappers for missing JS top-level that are not yet free in Go:

func GapF(steps int) Pattern { return Gap(steps) }
func SilenceF() Pattern { return Silence() }
func PureF(v any) Pattern { return Pure(v) }
func StackF(pats ...Pattern) Pattern { return Stack(pats...) }
func CatF(pats ...Pattern) Pattern { return Cat(pats...) }
func FastCatF(pats ...Pattern) Pattern { return FastCat(pats...) }
func SlowCatF(pats ...Pattern) Pattern { return SlowCat(pats...) }
func SequenceF(pats ...Pattern) Pattern { return Sequence(pats...) }
func ArrangeF(sections ...any) Pattern { return Arrange(sections...) }
func PolymeterF(pats ...Pattern) Pattern { return Polymeter(pats...) }
func TimeCatF(pats ...any) Pattern { return TimeCat(pats...) }

// Additional helpers: ensure all JS `register` free functions have Go free wrappers
// Many already exist as methods; provide free wrappers that take pat as last arg (like JS curry)
func AddF(a, b any) Pattern { return Reify(b).Add(a) }
func SubF(a, b any) Pattern { return Reify(b).Sub(a) }
func MulF(a, b any) Pattern { return Reify(b).Mul(a) }
func DivF(a, b any) Pattern { return Reify(b).Div(a) }
func ModF(a, b any) Pattern { return Reify(b).Mod(a) }
func PowF(a, b any) Pattern { return Reify(b).Pow(a) }
