// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — temporal/structure operations.

package core

// ---- Time transformations ----

// FastF is the non-patternified fast (factor is Fraction).
// Already defined in pattern.go as FastF; this file adds patternified variants and other time ops.

// Early shifts pattern early by offset (negative late).
func (p Pattern) Early(offset any) Pattern {
	frac := toFraction(Reify(offset).FirstCycleValue())
	res := p.WithQueryTime(func(t Fraction) Fraction { return t.Add(frac) }).WithHapTime(func(t Fraction) Fraction { return t.Sub(frac) })
	if p.Steps != nil {
		return res.SetSteps(*p.Steps)
	}
	return res
}

// EarlyF is helper with Fraction.
func (p Pattern) EarlyF(offset Fraction) Pattern {
	res := p.WithQueryTime(func(t Fraction) Fraction { return t.Add(offset) }).WithHapTime(func(t Fraction) Fraction { return t.Sub(offset) })
	if p.Steps != nil {
		return res.SetSteps(*p.Steps)
	}
	return res
}

// Late shifts pattern late.
func (p Pattern) Late(offset any) Pattern {
	frac := toFraction(Reify(offset).FirstCycleValue())
	return p.EarlyF(FractionFromInt(0).Sub(frac))
}

func (p Pattern) LateF(offset Fraction) Pattern {
	return p.EarlyF(FractionFromInt(0).Sub(offset))
}

// Compress compresses pattern into [b,e) sub-cycle.
func (p Pattern) Compress(b, e any) Pattern {
	bf := toFraction(Reify(b).FirstCycleValue())
	ef := toFraction(Reify(e).FirstCycleValue())
	if bf.Gt(ef) || bf.Gt(FractionFromInt(1)) || ef.Gt(FractionFromInt(1)) || bf.Lt(FractionFromInt(0)) || ef.Lt(FractionFromInt(0)) {
		return Silence()
	}
	dur := ef.Sub(bf)
	if dur.Equals(FractionFromInt(0)) {
		return Silence()
	}
	factor := FractionFromInt(1).Div(dur)
	return p.FastGapF(factor).LateF(bf)
}

// FastGapF speeds up leaving gap.
func (p Pattern) FastGapF(factor Fraction) Pattern {
	if factor.Equals(FractionFromInt(0)) {
		return Silence()
	}
	qf := func(span TimeSpan) *TimeSpan {
		cycle := span.Begin.Sam()
		bpos := span.Begin.Sub(cycle).Mul(factor)
		if bpos.Gt(FractionFromInt(1)) {
			bpos = FractionFromInt(1)
		}
		epos := span.End.Sub(cycle).Mul(factor)
		if epos.Gt(FractionFromInt(1)) {
			epos = FractionFromInt(1)
		}
		if bpos.Gte(FractionFromInt(1)) {
			return nil
		}
		return &TimeSpan{Begin: cycle.Add(bpos), End: cycle.Add(epos)}
	}
	ef := func(h Hap) Hap {
		begin := h.Part.Begin
		end := h.Part.End
		cycle := begin.Sam()
		beginPos := begin.Sub(cycle).Div(factor)
		if beginPos.Gt(FractionFromInt(1)) {
			beginPos = FractionFromInt(1)
		}
		endPos := end.Sub(cycle).Div(factor)
		if endPos.Gt(FractionFromInt(1)) {
			endPos = FractionFromInt(1)
		}
		newPart := NewTimeSpan(cycle.Add(beginPos), cycle.Add(endPos))
		var newWhole *TimeSpan
		if h.Whole != nil {
			newWholeBegin := newPart.Begin.Sub(begin.Sub(h.Whole.Begin).Div(factor))
			newWholeEnd := newPart.End.Add(h.Whole.End.Sub(end).Div(factor))
			w := NewTimeSpan(newWholeBegin, newWholeEnd)
			newWhole = &w
		}
		return NewHap(newWhole, newPart, h.Value, h.Context)
	}
	return p.WithQuerySpanMaybe(qf).WithHap(ef).SplitQueries()
}

// FastGap patternified.
func (p Pattern) FastGap(factor any) Pattern {
	pat := Reify(factor)
	if pat.PureValue != nil {
		return p.FastGapF(toFraction(pat.PureValue))
	}
	// patternified: use first value approximation
	frac := toFraction(pat.FirstCycleValue())
	return p.FastGapF(frac)
}

// Zoom selects [s,e) and stretches to whole cycle.
func (p Pattern) Zoom(s, e any) Pattern {
	sf := toFraction(Reify(s).FirstCycleValue())
	ef := toFraction(Reify(e).FirstCycleValue())
	if sf.Gte(ef) {
		return Silence()
	}
	d := ef.Sub(sf)
	steps := p.Steps
	if steps != nil {
		mul := d
		newSteps := steps.Mul(mul)
		steps = &newSteps
	}
	res := p.WithQuerySpan(func(span TimeSpan) TimeSpan {
		return span.WithCycle(func(t Fraction) Fraction { return t.Mul(d).Add(sf) })
	}).WithHapSpan(func(span TimeSpan) TimeSpan {
		return span.WithCycle(func(t Fraction) Fraction { return t.Sub(sf).Div(d) })
	}).SplitQueries()
	if steps != nil {
		return res.SetSteps(*steps)
	}
	return res
}

func stepsIfNotNil(s *Fraction) *Fraction { return s }

// Ply repeats each event factor times via squeezeJoin.
func (p Pattern) Ply(factor any) Pattern {
	frac := toFraction(Reify(factor).FirstCycleValue())
	if frac.Equals(FractionFromInt(0)) {
		return Silence()
	}
	result := p.Fmap(func(v any) any { return Pure(v).FastF(frac) }).SqueezeJoin()
	if stepsEnabled && p.Steps != nil {
		newSteps := frac.Mul(*p.Steps)
		result.Steps = &newSteps
	}
	return result
}

// Off superimposes function result offset in time.
func (p Pattern) Off(timePat any, fn func(Pattern) Pattern) Pattern {
	offPat := Reify(timePat)
	frac := toFraction(offPat.FirstCycleValue())
	return Stack(p, fn(p.LateF(frac)))
}

// When applies function when test is true (patternified boolean).
func (p Pattern) When(test any, fn func(Pattern) Pattern) Pattern {
	b := Reify(test).FirstCycleValue()
	apply := false
	switch v := b.(type) {
	case bool:
		apply = v
	case int:
		apply = v != 0
	case float64:
		apply = v != 0
	case Fraction:
		apply = !v.Equals(FractionFromInt(0))
	default:
		apply = b != nil && b != false
	}
	if apply {
		return fn(p)
	}
	return p
}

// Every applies function every n cycles.
//
// SplitQueries is required: the decision depends on the cycle number, so a
// query spanning several cycles must be broken into one query per cycle.
// Without it a render of many cycles reads the cycle from the span start and
// applies fn — or fails to — across the whole render.
func (p Pattern) Every(n int, fn func(Pattern) Pattern) Pattern {
	if n <= 0 {
		return p
	}
	return NewPattern(func(state State) []Hap {
		cycle := state.Span.Begin.Sam().Floor().Float64()
		if Mod(int(cycle), n) == 0 {
			return fn(p).Query(state)
		}
		return p.Query(state)
	}, p.Steps).SplitQueries()
}

// Inside carries operation inside cycle.
func (p Pattern) Inside(factor any, fn func(Pattern) Pattern) Pattern {
	frac := toFraction(Reify(factor).FirstCycleValue())
	return fn(p.SlowF(frac)).FastF(frac)
}

// Outside carries operation outside cycle.
func (p Pattern) Outside(factor any, fn func(Pattern) Pattern) Pattern {
	frac := toFraction(Reify(factor).FirstCycleValue())
	return fn(p.FastF(frac)).SlowF(frac)
}

// SlowF already in pattern.go as Slow, but add SlowAny for patternified.
func (p Pattern) SlowAny(factor any) Pattern {
	frac := toFraction(Reify(factor).FirstCycleValue())
	if frac.Equals(FractionFromInt(0)) {
		return Silence()
	}
	return p.FastF(FractionFromInt(1).Div(frac))
}

// Rev reverses each cycle.
func (p Pattern) Rev() Pattern {
	query := func(state State) []Hap {
		span := state.Span
		cycle := span.Begin.Sam()
		nextCycle := span.Begin.NextSam()
		reflect := func(ts TimeSpan) TimeSpan {
			// reflect = cycle + nextCycle - time
			b := cycle.Add(nextCycle.Sub(ts.Begin))
			e := cycle.Add(nextCycle.Sub(ts.End))
			// swap
			return NewTimeSpan(e, b)
		}
		haps := p.Query(state.SetSpan(reflect(span)))
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = h.WithSpan(reflect)
		}
		return out
	}
	return NewPattern(query, p.Steps).SplitQueries()
}

// SlowF helper for Inside/Outside if not already defined via Slow.
func (p Pattern) SlowF(frac Fraction) Pattern {
	return p.FastF(FractionFromInt(1).Div(frac))
}

// Helper to get first cycle value.
func (p Pattern) FirstCycleValue() any {
	haps := p.FirstCycle()
	if len(haps) == 0 {
		return nil
	}
	return haps[0].Value
}

// WithQuerySpanMaybe filters spans where fn returns nil.
func (p Pattern) WithQuerySpanMaybe(fn func(TimeSpan) *TimeSpan) Pattern {
	return NewPattern(func(state State) []Hap {
		newSpan := fn(state.Span)
		if newSpan == nil {
			return []Hap{}
		}
		return p.Query(state.SetSpan(*newSpan))
	}, p.Steps)
}

// Hurry speeds both pattern and speed control.
func (p Pattern) Hurry(r any) Pattern {
	frac := toFraction(Reify(r).FirstCycleValue())
	return p.FastF(frac).Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, val := range m {
				m2[k] = val
			}
			m2["speed"] = frac
			return m2
		}
		return map[string]any{"value": v, "speed": frac}
	})
}

// Helper for pure value steps preservation.
func (p Pattern) SetStepsCopy(s *Fraction) Pattern {
	if s == nil {
		return NewPattern(p.Query, nil)
	}
	cp := *s
	return NewPattern(p.Query, &cp)
}
