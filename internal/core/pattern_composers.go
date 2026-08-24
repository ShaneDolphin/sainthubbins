// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs COMPOSERS alignment matrix.

package core

// Alignment controls how value combination preserves structure.
type Alignment int

const (
	AlignIn Alignment = iota
	AlignOut
	AlignMix
	AlignSqueeze
	AlignSqueezeOut
	AlignReset
	AlignRestart
	AlignPoly
)

// _opIn etc — mirror JS _opIn/_opOut etc
func (p Pattern) OpIn(other Pattern, fn func(any) func(any) any) Pattern {
	return p.Fmap(func(a any) any { return fn(a) }).AppLeft(other.Fmap(func(b any) any { return b }))
}
func (p Pattern) OpOut(other Pattern, fn func(any) func(any) any) Pattern {
	return p.Fmap(func(a any) any { return fn(a) }).AppRight(other.Fmap(func(b any) any { return b }))
}
func (p Pattern) OpMix(other Pattern, fn func(any) func(any) any) Pattern {
	return p.Fmap(func(a any) any { return fn(a) }).AppBoth(other.Fmap(func(b any) any { return b }))
}
func (p Pattern) OpSqueeze(other Pattern, fn func(any) func(any) any) Pattern {
	otherPat := other
	return p.Fmap(func(a any) any {
		return func(b any) any { return fn(a)(b) }
	}).Fmap(func(v any) any {
		// Wrap as pattern of patterns for squeezeJoin
		fn2 := v.(func(any) any)
		return otherPat.Fmap(func(b any) any { return fn2(b) })
	}).SqueezeJoin()
}
func (p Pattern) OpSqueezeOut(other Pattern, fn func(any) func(any) any) Pattern {
	thisPat := p
	return other.Fmap(func(a any) any {
		return func(b any) any { return fn(b)(a) }
	}).Fmap(func(v any) any {
		fn2 := v.(func(any) any)
		return thisPat.Fmap(func(b any) any { return fn2(b) })
	}).SqueezeJoin()
}
func (p Pattern) OpReset(other Pattern, fn func(any) func(any) any) Pattern {
	otherPat := other
	return otherPat.Fmap(func(b any) any {
		return p.Fmap(func(a any) any { return fn(a)(b) })
	}).ResetJoin()
}
func (p Pattern) OpRestart(other Pattern, fn func(any) func(any) any) Pattern {
	otherPat := other
	return otherPat.Fmap(func(b any) any {
		return p.Fmap(func(a any) any { return fn(a)(b) })
	}).RestartJoin()
}
func (p Pattern) OpPoly(other Pattern, fn func(any) func(any) any) Pattern {
	otherPat := other
	return p.Fmap(func(b any) any {
		return otherPat.Fmap(func(a any) any { return fn(a)(b) })
	}).PolyJoin()
}

// Dispatch helper
func (p Pattern) WithAlignment(other Pattern, fn func(any) func(any) any, align Alignment) Pattern {
	switch align {
	case AlignIn:
		return p.OpIn(other, fn)
	case AlignOut:
		return p.OpOut(other, fn)
	case AlignMix:
		return p.OpMix(other, fn)
	case AlignSqueeze:
		return p.OpSqueeze(other, fn)
	case AlignSqueezeOut:
		return p.OpSqueezeOut(other, fn)
	case AlignReset:
		return p.OpReset(other, fn)
	case AlignRestart:
		return p.OpRestart(other, fn)
	case AlignPoly:
		return p.OpPoly(other, fn)
	default:
		return p.OpIn(other, fn)
	}
}
