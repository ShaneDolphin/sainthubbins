// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Extra Pattern methods: echo/stut, plyWith/plyForEach, arrange, mask/structAll, weave stub, distort stubs.

package core

import "math"

// Echo repeats with offset and gain feedback
func (p Pattern) Echo(times, time, feedback any) Pattern {
	nt := int(toFloat(Reify(times).FirstCycleValue()))
	tf := toFloat(Reify(time).FirstCycleValue())
	fb := toFloat(Reify(feedback).FirstCycleValue())
	if nt <= 0 {
		return p
	}
	pats := make([]Pattern, 0, nt)
	for i := 0; i < nt; i++ {
		gain := math.Pow(fb, float64(i))
		offset := FractionFromFloat(tf * float64(i))
		pats = append(pats, p.Late(offset).Fmap(func(v any) any {
			if m, ok := v.(map[string]any); ok {
				m2 := map[string]any{}
				for k, vv := range m {
					m2[k] = vv
				}
				m2["gain"] = gain
				if g, ok := m["gain"]; ok {
					m2["gain"] = toFloat(g) * gain
				}
				return m2
			}
			return map[string]any{"value": v, "gain": gain}
		}))
	}
	return Stack(pats...)
}

// Stut like echo but flipped params
func (p Pattern) Stut(times, feedback, td any) Pattern {
	return p.Echo(times, td, feedback)
}

// EchoWith generic
func (p Pattern) EchoWith(times, tm any, fn func(Pattern, int) Pattern) Pattern {
	nt := int(toFloat(Reify(times).FirstCycleValue()))
	tf := toFloat(Reify(tm).FirstCycleValue())
	pats := make([]Pattern, 0, nt)
	for i := 0; i < nt; i++ {
		offset := FractionFromFloat(tf * float64(i))
		pats = append(pats, fn(p.Late(offset), i))
	}
	return Stack(pats...)
}

// PlyWith repeats with function
func (p Pattern) PlyWith(factor any, fn func(any) any) Pattern {
	f := int(toFloat(Reify(factor).FirstCycleValue()))
	if f <= 0 {
		return p
	}
	return p.Fmap(func(v any) any {
		pats := make([]Pattern, 0, f)
		for i := 0; i < f; i++ {
			val := v
			if i > 0 {
				// apply fn with index i? simplified: call fn on value
				if fn != nil {
					val = fn(v)
				}
			}
			pats = append(pats, Pure(val))
		}
		return FastCat(pats...).FastF(FractionFromInt(int64(f)))
	}).SqueezeJoin()
}

// PlyForEach like plyWith but passes index
func (p Pattern) PlyForEach(factor any, fn func(Pattern, int) Pattern) Pattern {
	f := int(toFloat(Reify(factor).FirstCycleValue()))
	if f <= 0 {
		return p
	}
	return p.Fmap(func(v any) any {
		pats := make([]Pattern, 0, f)
		pats = append(pats, Pure(v))
		for i := 1; i < f; i++ {
			pats = append(pats, fn(Pure(v), i))
		}
		return FastCat(pats...).FastF(FractionFromInt(int64(f)))
	}).SqueezeJoin()
}

// Arrange cycles through sections (simplified as SlowCat of pats weighted by cycles)
func Arrange2(sections ...any) Pattern {
	if len(sections) == 0 {
		return Silence()
	}
	// sections are [cycles, pat, cycles, pat ...] already handled in pattern_structure.go Arrange
	return Arrange(sections...)
}

// Mask keeps when binary true
func (p Pattern) Mask(maskPat any) Pattern {
	mp := Reify(maskPat)
	return p.KeepIf(mp)
}
func (p Pattern) MaskAll(maskPat any) Pattern { return p.Keep(Reify(maskPat)) }

// StructAll keeps all (alias)
func (p Pattern) StructAll(pat any) Pattern { return p.Keep(Reify(pat)) }

// Weave distributes patterns over time — JS weaveWith: stack(...funcs.map((func,i)=> pat.inside(t,func).early(Fraction(i).div(l))))._slow(t)
// and weave: this.weaveWith(t, ...pats.map((x)=> set.out(x)))
func (p Pattern) Weave(t any, pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	// Convert pats to funcs via set.out equivalent: patch pattern onto base
	fns := make([]func(Pattern) Pattern, len(pats))
	for i, pat := range pats {
		captured := pat
		fns[i] = func(base Pattern) Pattern {
			// set.out: combine base and captured such that captured overrides base's controls?
			// Simplified as captured's value bags merged onto base via Set; approximate as Stack/Cat?
			// For tests, use base.Set(captured) if Set exists, else just captured
			return base.Set(captured)
		}
		_ = captured
	}
	return p.WeaveWith(t, fns...)
}
func (p Pattern) WeaveWith(t any, fns ...func(Pattern) Pattern) Pattern {
	if len(fns) == 0 {
		return Silence()
	}
	tFrac := FractionFromInt(1)
	switch v := t.(type) {
	case int:
		tFrac = FractionFromInt(int64(v))
	case int64:
		tFrac = FractionFromInt(v)
	case float64:
		tFrac = FractionFromFloat(v)
	case Fraction:
		tFrac = v
	case *Fraction:
		if v != nil {
			tFrac = *v
		}
	default:
		if f, ok := v.(Fraction); ok {
			tFrac = f
		} else {
			tFrac = FractionFromFloat(toFloat(v))
		}
	}
	l := len(fns)
	pats := make([]Pattern, 0, l)
	for i, fn := range fns {
		inside := p.Inside(tFrac, fn)
		earlyFrac := FractionFromInt(int64(i)).Div(FractionFromInt(int64(l)))
		pats = append(pats, inside.Early(earlyFrac))
	}
	stacked := Stack(pats...)
	// _slow(t) => slow by t
	return stacked.Slow(tFrac)
}

// Distort stubs (superdough) — just set distort field
func (p Pattern) Distort(args any) Pattern {
	return p.Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["distort"] = args
			return m2
		}
		return map[string]any{"value": v, "distort": args}
	})
}
func distortWith(name string) func(any, Pattern) Pattern {
	return func(args any, pat Pattern) Pattern {
		return pat.Distort([]any{args, 1, name})
	}
}
func (p Pattern) Soft(args any) Pattern      { return p.Distort([]any{args, 1, "soft"}) }
func (p Pattern) Hard(args any) Pattern      { return p.Distort([]any{args, 1, "hard"}) }
func (p Pattern) Cubic(args any) Pattern     { return p.Distort([]any{args, 1, "cubic"}) }
func (p Pattern) Diode(args any) Pattern     { return p.Distort([]any{args, 1, "diode"}) }
func (p Pattern) Asym(args any) Pattern      { return p.Distort([]any{args, 1, "asym"}) }
func (p Pattern) Fold(args any) Pattern      { return p.Distort([]any{args, 1, "fold"}) }
func (p Pattern) Sinefold(args any) Pattern  { return p.Distort([]any{args, 1, "sinefold"}) }
func (p Pattern) Chebyshev(args any) Pattern { return p.Distort([]any{args, 1, "chebyshev"}) }

// Pace alias for steps
func (p Pattern) Pace(targetSteps any) Pattern {
	ts := toFloat(Reify(targetSteps).FirstCycleValue())
	if ts == 0 {
		return p
	}
	// Simplified: fast by ratio
	return p.Fast(Pure(ts))
}
func Pace2(targetSteps any, pat Pattern) Pattern { return pat.Pace(targetSteps) }

// Contract/Expand etc already in pattern.go via stepRegister; provide simple wrappers
// Already have Expand/Contract via step helpers, but ensure they exist
