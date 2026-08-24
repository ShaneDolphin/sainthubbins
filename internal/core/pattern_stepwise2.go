// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Stepwise aliases and small missing free functions.

package core

// Steps alias for Pace (JS: steps = pace)
func (p Pattern) Steps2(target any) Pattern { return p.Pace(target) }

// ShrinkList/GrowList stepwise (simplified wrappers)
func (p Pattern) ShrinkList(amount any) Pattern {
	amt := int(toFloat(Reify(amount).FirstCycleValue()))
	if amt <= 0 {
		amt = 1
	}
	return p.Shrink(amt)
}
func (p Pattern) GrowList(amount any) Pattern {
	amt := int(toFloat(Reify(amount).FirstCycleValue()))
	if amt <= 0 {
		amt = 1
	}
	return p.Grow(amt)
}

// S aliases for stepcat variants (JS: s_cat = stepcat, s_add = take, etc)
func SCat(pats ...Pattern) Pattern      { return SlowCat(pats...) }
func SAlt(groups ...[]Pattern) Pattern {
	// stepalt: alternate groups stepwise half? Simplified as SlowCat of each group cat
	var pats []Pattern
	for _, g := range groups {
		pats = append(pats, FastCat(g...))
	}
	return SlowCat(pats...)
}
func SPolymeter(pats ...Pattern) Pattern { return Polymeter(pats...) }
func STaper(amount any, pat Pattern) Pattern { return pat.Shrink(amount) }
func STaperList(amount any, pat Pattern) Pattern { return pat.ShrinkList(amount) }
func SAdd(n any, pat Pattern) Pattern {
	f := toFloat(Reify(n).FirstCycleValue())
	if f <= 0 {
		f = 1
	}
	return pat.Shrink(f)
}
func SSub(n any, pat Pattern) Pattern {
	f := toFloat(Reify(n).FirstCycleValue())
	if f <= 0 {
		f = 1
	}
	return pat.Grow(f)
}
func SExpand(factor any, pat Pattern) Pattern {
	f := toFloat(Reify(factor).FirstCycleValue())
	if f == 0 {
		f = 1
	}
	return pat.Fast(Pure(f))
}
func SExtend(factor any, pat Pattern) Pattern {
	f := toFloat(Reify(factor).FirstCycleValue())
	if f == 0 {
		f = 1
	}
	return pat.Slow(FractionFromFloat(f))
}
func SContract(factor any, pat Pattern) Pattern {
	f := toFloat(Reify(factor).FirstCycleValue())
	if f == 0 {
		f = 1
	}
	return pat.Compress(Pure(0), Pure(1.0/f))
}
func STour(pat Pattern, many ...Pattern) Pattern { return pat.Tour(many...) }
func SZip(pats ...Pattern) Pattern { return Zip(pats...) }

// Gap alias free functions (lowercase aliases for JS compat)
func GapPat(steps int) Pattern { return Gap(steps) }

// Polyrhythm aliases
func Polyrhythm(pats ...Pattern) Pattern { return Stack(pats...) }
func Pm2(pats ...Pattern) Pattern       { return Polymeter(pats...) }
func Pr2(pats ...Pattern) Pattern       { return Stack(pats...) }

// Silence alias (lowercase)
func SilencePat() Pattern { return Silence() }

// Pure alias already exists as Pure

// StackBy already exists in pattern_poly.go; provide lowercase alias wrapper
func StackCentre(pats ...Pattern) Pattern { return Stack(pats...) }
func StackLeft(pats ...Pattern) Pattern   { return Stack(pats...) }
func StackRight(pats ...Pattern) Pattern  { return Stack(pats...) }

// Arrange alias already exists as Arrange; provide pure free

// Ratio already exists as method; also free Ratio? Already method Ratio; free wrapper
func RatioPat(pat Pattern) Pattern { return pat.Ratio() }

// Cpm free wrapper
func CpmPat(cpm any, pat Pattern) Pattern { return pat.Cpm(cpm) }

// Revv free wrapper
func RevvPat(pat Pattern) Pattern { return pat.Revv() }

// Bypass free wrapper
func BypassPat(on any, pat Pattern) Pattern { return pat.Bypass(on) }

// Extra free wrappers for missing exports that are just functions returning Patterns
func GapFunc(steps int) Pattern { return Gap(steps) }
func PaceFree(target any, pat Pattern) Pattern { return pat.Pace(target) }
