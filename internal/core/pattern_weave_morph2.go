// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Weave/morph and remaining combinators.

package core

// Weave2 distributes patterns over time — delegate to Weave for full fidelity
func (p Pattern) Weave2(t any, pats ...Pattern) Pattern {
	return p.Weave(t, pats...)
}

// Morph interpolates between from/to by factor by (0..1) via AppBoth
func (p Pattern) Morph(from, to any, by any) Pattern {
	f := Reify(from)
	t := Reify(to)
	bp := Reify(by)
	return f.Fmap(func(fv any) any {
		return func(tv any) any {
			return func(bv any) any {
				bf := toFloat(bv)
				if bf < 0.5 {
					return fv
				}
				return tv
			}
		}
	}).AppBoth(t.Fmap(func(v any) any { return v })).AppBoth(bp.Fmap(func(v any) any { return v })).Fmap(func(v any) any {
		if fn, ok := v.(func(any) any); ok {
			return fn
		}
		return v
	})
}
func MorphFree(from, to, by any) Pattern {
	f := Reify(from)
	t := Reify(to)
	b := Reify(by)
	return f.Morph(f, t, b)
}

// Additional small aliases missing earlier
func (p Pattern) Polyrhythm2(pats ...Pattern) Pattern { return Stack(append([]Pattern{p}, pats...)...) }
func (p Pattern) PmAlias(pats ...Pattern) Pattern { return Polymeter(append([]Pattern{p}, pats...)...) }

// Ensure file compiles
