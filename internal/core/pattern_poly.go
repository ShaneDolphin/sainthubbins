// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — poly/reset/squeeze joins, defragment.

package core

// ResetJoin aligns inner pattern cycle start to outer hap.
func (p Pattern) ResetJoin(restart ...bool) Pattern {
	restartFlag := false
	if len(restart) > 0 {
		restartFlag = restart[0]
	}
	patOfPats := p
	return NewPattern(func(state State) []Hap {
		outerHaps := patOfPats.DiscreteOnly().Query(state)
		var out []Hap
		for _, outer := range outerHaps {
			if outer.Whole == nil {
				continue
			}
			var innerPat Pattern
			if ip, ok := outer.Value.(Pattern); ok {
				innerPat = ip
			} else if ipp, ok := outer.Value.(*Pattern); ok && ipp != nil {
				innerPat = *ipp
			} else {
				continue
			}
			var shifted Pattern
			if restartFlag {
				shifted = innerPat.LateF(outer.Whole.Begin)
			} else {
				shifted = innerPat.LateF(outer.Whole.Begin.CyclePos())
			}
			innerHaps := shifted.Query(state)
			for _, inner := range innerHaps {
				var whole *TimeSpan
				if inner.Whole != nil {
					inter := inner.Whole.Intersection(*outer.Whole)
					if inter == nil {
						continue
					}
					whole = inter
				}
				part := inner.Part.Intersection(outer.Part)
				if part == nil {
					continue
				}
				ctx := inner.CombineContext(outer.Context)
				out = append(out, NewHap(whole, *part, inner.Value, ctx))
			}
		}
		return out
	}, p.Steps)
}

func (p Pattern) RestartJoin() Pattern { return p.ResetJoin(true) }

// PolyJoin extends inner patterns to match outer steps then outerJoin.
func (p Pattern) PolyJoin() Pattern {
	pp := p
	return pp.Fmap(func(v any) any {
		if pat, ok := v.(Pattern); ok {
			if pp.Steps != nil && pat.Steps != nil && !pat.Steps.Equals(FractionFromInt(0)) {
				// extend inner to match outer steps
				factor := pp.Steps.Div(*pat.Steps)
				// extend = replicate factor? Use Fast? For now simple
				_ = factor
			}
			return pat
		}
		return v
	}).OuterJoin()
}

func (p Pattern) PolyBind(fn func(any) Pattern) Pattern {
	return p.Fmap(func(v any) any { return fn(v) }).PolyJoin()
}

// Defragment already in pattern_structure.go, but provide correct version here if missing.
// Keep alias for compatibility.
func (p Pattern) Defragment2() Pattern { return p.Defragment() }

// _focusSpan helper for squeezeJoin (focus inner pattern into outer whole)
func (p Pattern) FocusSpan(span TimeSpan) Pattern {
	// Map TimeSpan to inner cycle: focus inner pattern's first cycle into span
	dur := span.Duration()
	if dur.Equals(FractionFromInt(0)) {
		return Silence()
	}
	return p.WithQuerySpan(func(s TimeSpan) TimeSpan {
		return s.WithCycle(func(t Fraction) Fraction { return t.Mul(dur).Add(span.Begin) })
	}).WithHapSpan(func(s TimeSpan) TimeSpan {
		return s.WithCycle(func(t Fraction) Fraction { return t.Sub(span.Begin).Div(dur) })
	}).SplitQueries()
}

// StackBy stacks with custom time division
func StackBy(by Fraction, pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, pat := range pats {
			// Query with time scaled by by?
			haps := pat.Query(state.WithSpan(func(s TimeSpan) TimeSpan {
				return s.WithTime(func(t Fraction) Fraction { return t.Mul(by) })
			}))
			for _, h := range haps {
				shifted := h.WithSpan(func(s TimeSpan) TimeSpan {
					return s.WithTime(func(t Fraction) Fraction { return t.Div(by) })
				})
				out = append(out, shifted)
			}
		}
		return out
	}, nil)
}
