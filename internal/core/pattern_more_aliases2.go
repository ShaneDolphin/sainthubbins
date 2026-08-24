// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Additional free wrappers for completeness (JS top-level curried functions)

package core

func ApplyN2(n int, fn func(any) any, pat Pattern) Pattern {
	result := pat
	for i := 0; i < n; i++ {
		result = result.Fmap(fn)
	}
	return result
}

func PlyPat(factor any, pat Pattern) Pattern {
	f := int(toFloat(Reify(factor).FirstCycleValue()))
	return pat.Ply(f)
}

func WhenPat(test any, fn func(Pattern) Pattern, pat Pattern) Pattern {
	return pat.When(test, fn)
}

func OffPat(timePat any, fn func(Pattern) Pattern, pat Pattern) Pattern {
	return pat.Off(timePat, fn)
}

func EveryPat(n int, fn func(Pattern) Pattern, pat Pattern) Pattern {
	return pat.Every(n, fn)
}
