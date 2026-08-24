// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — TimeCat, weight handling, polymeter.

package core

// TimeCat stacks patterns weighted by durations (like polymeter slowcat variant).
// Usage: TimeCat([1, pat1], [2, pat2]) etc via varargs: TimeCat(dur1, pat1, dur2, pat2...)
func TimeCatWeighted(pairs ...any) Pattern {
	if len(pairs)%2 != 0 {
		// odd, last is pat without duration? treat as 1
		pairs = append(pairs, FractionFromInt(1))
	}
	type weighted struct {
		dur Fraction
		pat Pattern
	}
	var w []weighted
	total := FractionFromInt(0)
	for i := 0; i < len(pairs); i += 2 {
		var dur Fraction
		switch v := pairs[i].(type) {
		case Fraction:
			dur = v
		case *Fraction:
			dur = *v
		case int:
			dur = FractionFromInt(int64(v))
		case float64:
			dur = FractionFromFloat(v)
		default:
			dur = toFraction(v)
		}
		var pat Pattern
		switch v := pairs[i+1].(type) {
		case Pattern:
			pat = v
		case *Pattern:
			pat = *v
		default:
			pat = Reify(v)
		}
		w = append(w, weighted{dur, pat})
		total = total.Add(dur)
	}
	if total.Equals(FractionFromInt(0)) {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			base := sub.Begin.Sam()
			acc := FractionFromInt(0)
			for _, wp := range w {
				segBegin := base.Add(acc.Div(total))
				segEnd := segBegin.Add(wp.dur.Div(total))
				seg := NewTimeSpan(segBegin, segEnd)
				inter := seg.Intersection(sub)
				if inter == nil {
					acc = acc.Add(wp.dur)
					continue
				}
				// Map inter to pat time: (t - segBegin) * total/dur
				factor := total.Div(wp.dur)
				mapped := inter.WithTime(func(t Fraction) Fraction { return t.Sub(segBegin).Mul(factor) })
				haps := wp.pat.Query(state.SetSpan(mapped))
				for _, h := range haps {
					nh := h.WithSpan(func(s TimeSpan) TimeSpan {
						return s.WithTime(func(t Fraction) Fraction { return t.Div(factor).Add(segBegin) })
					})
					if nh.Part.Intersection(sub) != nil {
						out = append(out, nh)
					}
				}
				acc = acc.Add(wp.dur)
			}
		}
		return out
	}, nil)
}

// StepCat concatenates patterns stepwise weighted by durations, mirroring JS stepcat/timeCat.
// Usage: StepCat([dur, pat], pat2, ...) or StepCat(dur, pat, dur, pat) via TimeCatWeighted.
// JS: stepcat([3,"e3"],[1,"g3"]) or timecat("bd","sd")
func StepCat(args ...any) Pattern {
	if len(args) == 0 {
		return Silence()
	}
	type pair struct {
		dur Fraction
		pat Pattern
	}
	var pairs []pair
	for _, arg := range args {
		switch v := arg.(type) {
		case []any:
			if len(v) == 2 {
				dur := toFraction(v[0])
				var pat Pattern
				switch pp := v[1].(type) {
				case Pattern:
					pat = pp
				case *Pattern:
					pat = *pp
				default:
					pat = Reify(pp)
				}
				pairs = append(pairs, pair{dur, pat})
			} else if len(v) == 1 {
				// single pat without duration -> 1
				pat := Reify(v[0])
				dur := FractionFromInt(1)
				if pat.Steps != nil {
					dur = *pat.Steps
				}
				pairs = append(pairs, pair{dur, pat})
			}
		case Pattern:
			dur := FractionFromInt(1)
			if v.Steps != nil {
				dur = *v.Steps
			}
			pairs = append(pairs, pair{dur, v})
		default:
			pat := Reify(v)
			dur := FractionFromInt(1)
			if pat.Steps != nil {
				dur = *pat.Steps
			}
			pairs = append(pairs, pair{dur, pat})
		}
	}
	// Normalize undefined durations: if any dur is zero and others non-zero, replace with average
	hasZero := false
	hasNonZero := false
	for _, pr := range pairs {
		if pr.dur.Equals(FractionFromInt(0)) {
			hasZero = true
		} else {
			hasNonZero = true
		}
	}
	if hasZero && hasNonZero {
		// compute average of non-zero
		sum := FractionFromInt(0)
		count := 0
		for _, pr := range pairs {
			if !pr.dur.Equals(FractionFromInt(0)) {
				sum = sum.Add(pr.dur)
				count++
			}
		}
		avg := sum.Div(FractionFromInt(int64(count)))
		for i := range pairs {
			if pairs[i].dur.Equals(FractionFromInt(0)) {
				pairs[i].dur = avg
			}
		}
	} else if !hasNonZero {
		// all zero or none — fallback to fastcat
		var pats []Pattern
		for _, pr := range pairs {
			pats = append(pats, pr.pat)
		}
		return FastCat(pats...)
	}
	if len(pairs) == 1 {
		result := pairs[0].pat
		d := pairs[0].dur
		result.Steps = &d
		return result
	}
	total := FractionFromInt(0)
	for _, pr := range pairs {
		total = total.Add(pr.dur)
	}
	if total.Equals(FractionFromInt(0)) {
		return Silence()
	}
	begin := FractionFromInt(0)
	var compressed []Pattern
	for _, pr := range pairs {
		if pr.dur.Equals(FractionFromInt(0)) {
			continue
		}
		end := begin.Add(pr.dur)
		bNorm := begin.Div(total)
		eNorm := end.Div(total)
		compressed = append(compressed, pr.pat.Compress(bNorm, eNorm))
		begin = end
	}
	result := Stack(compressed...)
	result.Steps = &total
	return result
}

// PolymeterSlowcat stacks patterns slowed by their weight (like stack with _slow)
func PolymeterSlowcat(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	slowed := append([]Pattern(nil), pats...)
	return Stack(slowed...)
}
