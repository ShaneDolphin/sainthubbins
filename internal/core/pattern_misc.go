// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Misc Pattern methods: cpm, ratio, revv, bypass, ribbon, hsl/hsla, compressSpan/focusSpan, density/sparsity, etc.

package core

import "fmt"

// Cpm multiplies cps by cpm/60
func (p Pattern) Cpm(cpm any) Pattern {
	f := toFloat(Reify(cpm).FirstCycleValue())
	if f == 0 {
		f = 60
	}
	return p.Fast(Pure(f / 60.0 / 1.0))
}

// Ratio computes ratio from array [base, ...ratios]
func (p Pattern) Ratio() Pattern {
	return p.Fmap(func(v any) any {
		arr, ok := v.([]any)
		if !ok {
			return v
		}
		if len(arr) == 0 {
			return v
		}
		base := toFloat(arr[0])
		acc := base
		for _, n := range arr[1:] {
			nf := toFloat(n)
			if nf == 0 {
				nf = 1
			}
			acc /= nf
		}
		return acc
	})
}

// Revv reverses whole pattern via negating span
func (p Pattern) Revv() Pattern {
	negate := func(ts TimeSpan) TimeSpan {
		return NewTimeSpan(FractionFromInt(0).Sub(ts.End), FractionFromInt(0).Sub(ts.Begin))
	}
	return p.WithQuerySpan(negate).WithHapSpan(negate)
}

// Bypass silences if on true
func (p Pattern) Bypass(on any) Pattern {
	b := Reify(on).FirstCycleValue()
	onBool := false
	switch v := b.(type) {
	case bool:
		onBool = v
	case int:
		onBool = v != 0
	case float64:
		onBool = v != 0
	case string:
		onBool = v != "" && v != "0" && v != "false"
	default:
		if s := fmt.Sprintf("%v", v); s == "1" || s == "true" {
			onBool = true
		}
	}
	if onBool {
		return Silence()
	}
	return p
}

// Ribbon loops offset for cycles
func (p Pattern) Ribbon(offset, cycles any) Pattern {
	off := toFloat(Reify(offset).FirstCycleValue())
	return p.Early(Pure(off)).WithQuerySpan(func(ts TimeSpan) TimeSpan {
		// restart every cycles: simplified as just early
		return ts
	})
}
func (p Pattern) Rib(offset, cycles any) Pattern { return p.Ribbon(offset, cycles) }

// Hsla/Hsl color helpers (set color field)
func (p Pattern) Hsla(h, s, l, a any) Pattern {
	hv := toFloat(Reify(h).FirstCycleValue())
	sv := toFloat(Reify(s).FirstCycleValue())
	lv := toFloat(Reify(l).FirstCycleValue())
	av := toFloat(Reify(a).FirstCycleValue())
	return p.Fmap(func(v any) any {
		color := fmt.Sprintf("hsla(%vturn,%v%%,%v%%,%v)", hv, sv*100, lv*100, av)
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["color"] = color
			return m2
		}
		return map[string]any{"value": v, "color": color}
	})
}
func (p Pattern) Hsl(h, s, l any) Pattern { return p.Hsla(h, s, l, Pure(1)) }

// CompressSpan synonym (FocusSpan already in pattern_poly.go)
func (p Pattern) CompressSpan(span any) Pattern {
	s := Reify(span).FirstCycleValue()
	if ts, ok := s.(TimeSpan); ok {
		return p.Compress(Pure(ts.Begin.Float64()), Pure(ts.End.Float64()))
	}
	return p
}

// Focus is like compress but with gap handling? Simplified as Compress
func (p Pattern) Focus(b, e any) Pattern { return p.Compress(b, e) }

// Density/Sparsity synonyms for Fast/Slow
func (p Pattern) Density(factor any) Pattern {
	if f := Reify(factor).FirstCycleValue(); f != nil {
		return p.Fast(Pure(toFloat(f)))
	}
	return p
}
func (p Pattern) Sparsity(factor any) Pattern {
	if f := Reify(factor).FirstCycleValue(); f != nil {
		return p.SlowF(FractionFromFloat(toFloat(f)))
	}
	return p
}

// Gap creates silence with steps
func Gap(steps int) Pattern { return NewPattern(func(State) []Hap { return []Hap{} }, func() *Fraction { f := FractionFromInt(int64(steps)); return &f }()) }

// Silence / Nothing free functions (already Silence() exists)
func Nothing() Pattern { return Gap(0) }

// Pure wrapper already exists; Cat/Slowcat etc free already

// StackBy alias already in pattern_poly.go

// FirstOf / LastOf helpers
func (p Pattern) FirstOf(n int, fn func(Pattern) Pattern) Pattern {
	return p.Every(n, fn)
}
func (p Pattern) LastOf(n int, fn func(Pattern) Pattern) Pattern {
	cycle := n - 1
	if cycle < 0 {
		cycle = 0
	}
	// SplitQueries for the same reason as Every: the decision is per cycle.
	return NewPattern(func(state State) []Hap {
		c := state.Span.Begin.Sam().Floor().Float64()
		if int(c)%n == cycle {
			return fn(p).Query(state)
		}
		return p.Query(state)
	}, p.Steps).SplitQueries()
}

// Apply etc (simplified as Fmap)
func (p Pattern) Apply(fn func(any) any) Pattern { return p.Fmap(fn) }

// Cat etc are already free functions in pattern.go (Cat, SlowCat, FastCat, Polymeter)
// Provide only missing aliases
func SlowcatPrime(pats ...Pattern) Pattern { return SlowCat(pats...) }
func SequenceP(pats ...Pattern) Pattern { return FastCat(pats...) }
func Pm(pats ...Pattern) Pattern       { return Polymeter(pats...) }
func Pr(pats ...Pattern) Pattern       { return Stack(pats...) }
