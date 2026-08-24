// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — sample slicing, chop/striate/slice/splice/fit/loopAt, bite/linger/segment, etc.

package core

// Chop cuts each sample into n parts via squeezeBind, merging slice begin/end.
func (p Pattern) Chop(n int) Pattern {
	if n <= 0 {
		return p
	}
	return p.SqueezeBind(func(v any) Pattern {
		// v is value bag; create n slices
		base, isMap := v.(map[string]any)
		if !isMap {
			base = map[string]any{"value": v}
		}
		var pats []Pattern
		for i := 0; i < n; i++ {
			b := float64(i) / float64(n)
			e := float64(i+1) / float64(n)
			slice := map[string]any{}
			for k, val := range base {
				slice[k] = val
			}
			slice["begin"] = b
			slice["end"] = e
			pats = append(pats, Pure(slice))
		}
		return Sequence(pats...)
	}).WithSteps(func(f Fraction) Fraction { return f.Mul(FractionFromInt(int64(n))) })
}

// Striate cuts into n parts but keeps timing (fast + set slicePat).
func (p Pattern) Striate(n int) Pattern {
	if n <= 0 {
		return p
	}
	sliceObjs := make([]any, n)
	for i := 0; i < n; i++ {
		sliceObjs[i] = map[string]any{"begin": float64(i) / float64(n), "end": float64(i+1) / float64(n)}
	}
	// slowcat of slice objects as patterns
	slicePats := make([]Pattern, n)
	for i, o := range sliceObjs {
		slicePats[i] = Pure(o)
	}
	slicePat := SlowCat(slicePats...)
	return p.Set(slicePat).FastF(FractionFromInt(int64(n)))
}

// Slice triggers slices of n with index pattern ipat on object pat (opat).
// Mirrors JS slice(npat, ipat, opat) outerBind chain.
func Slice(npat, ipat, opat any) Pattern {
	np := Reify(npat)
	ip := Reify(ipat)
	op := Reify(opat)
	return np.InnerBind(func(nv any) Pattern {
		n := toFloat(nv)
		ni := int(n)
		if ni <= 0 {
			ni = 1
		}
		return ip.OuterBind(func(iv any) Pattern {
			idx := int(toFloat(iv))
			return op.OuterBind(func(ov any) Pattern {
				m, isMap := ov.(map[string]any)
				if !isMap {
					if s, ok := ov.(string); ok {
						m = map[string]any{"s": s}
					} else {
						m = map[string]any{"value": ov}
					}
				}
				var b, e float64
				// if n is array (slice), not supported fully; fallback to i/n
				b = float64(idx) / float64(ni)
				e = float64(idx+1) / float64(ni)
				out := map[string]any{}
				for k, v := range m {
					out[k] = v
				}
				out["begin"] = b
				out["end"] = e
				out["_slices"] = ni
				return Pure(out)
			})
		})
	})
}

// Splice like slice but sets speed to fit duration.
func Splice(npat, ipat, opat any) Pattern {
	sliced := Slice(npat, ipat, opat)
	return NewPattern(func(state State) []Hap {
		cps := 1.0
		if v, ok := state.Controls["_cps"]; ok {
			cps = toFloat(v)
			if cps == 0 {
				cps = 1
			}
		}
		haps := sliced.Query(state)
		for i, h := range haps {
			if m, ok := h.Value.(map[string]any); ok {
				nSlices := 1.0
				if ns, ok := m["_slices"]; ok {
					nSlices = toFloat(ns)
				}
				dur := FractionFromInt(1)
				if h.Whole != nil {
					dur = h.Whole.End.Sub(h.Whole.Begin)
				}
				speed := (cps / nSlices / dur.Float64())
				if s, ok := m["speed"]; ok {
					speed *= toFloat(s)
				}
				m2 := map[string]any{}
				for k, v := range m {
					m2[k] = v
				}
				m2["speed"] = speed
				m2["unit"] = "c"
				haps[i] = h.WithValue(func(any) any { return m2 })
			}
		}
		return haps
	}, sliced.Steps)
}

// Fit makes sample fit its event duration.
func (p Pattern) Fit() Pattern {
	return p.WithHaps(func(haps []Hap, state State) []Hap {
		cps := 1.0
		if v, ok := state.Controls["_cps"]; ok {
			cps = toFloat(v)
		}
		if cps == 0 {
			cps = 1
		}
		for i, h := range haps {
			if m, ok := h.Value.(map[string]any); ok {
				b := 0.0
				e := 1.0
				if bv, ok := m["begin"]; ok {
					b = toFloat(bv)
				}
				if ev, ok := m["end"]; ok {
					e = toFloat(ev)
				}
				slicedur := e - b
				dur := 1.0
				if h.Whole != nil {
					dur = h.Whole.End.Sub(h.Whole.Begin).Float64()
				}
				if dur == 0 {
					dur = 1
				}
				speed := (cps / dur) * slicedur
				if s, ok := m["speed"]; ok {
					speed *= toFloat(s)
				}
				m2 := map[string]any{}
				for k, v := range m {
					m2[k] = v
				}
				m2["speed"] = speed
				m2["unit"] = "c"
				haps[i] = h.WithValue(func(any) any { return m2 })
			}
		}
		return haps
	})
}

// LoopAt makes sample fit factor cycles.
func (p Pattern) LoopAt(factor any) Pattern {
	f := toFloat(Reify(factor).FirstCycleValue())
	if f == 0 {
		f = 1
	}
	return p.Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["speed"] = (1.0 / f) * 0.5
			if s, ok := m["speed"]; ok {
				m2["speed"] = toFloat(s) * (1.0 / f) * 0.5
			}
			m2["unit"] = "c"
			return m2
		}
		return v
	}).SlowF(FractionFromFloat(f))
}

// Bite selects fraction of pattern via zoom (npat controls division, ipat indexes).
func Bite(npat, ipat, pat any) Pattern {
	np := Reify(npat)
	ip := Reify(ipat)
	pp := Reify(pat)
	return ip.Fmap(func(iv any) any {
		return func(nv any) any {
			i := toFloat(iv)
			n := toFloat(nv)
			if n == 0 {
				n = 1
			}
			a := FractionFromFloat(i).Div(FractionFromFloat(n)).Mod(FractionFromInt(1))
			b := a.Add(FractionFromInt(1).Div(FractionFromFloat(n)))
			return pp.Zoom(a, b)
		}
	}).AppLeft(np).SqueezeJoin()
}

// Linger selects fraction and repeats to fill cycle.
func (p Pattern) Linger(t float64) Pattern {
	if t == 0 {
		return Silence()
	}
	if t < 0 {
		return p.Zoom(FractionFromFloat(t).Add(FractionFromInt(1)), FractionFromInt(1)).SlowF(FractionFromFloat(t))
	}
	return p.Zoom(FractionFromInt(0), FractionFromFloat(t)).SlowF(FractionFromFloat(t))
}

// Segment samples pattern at n events per cycle.
func (p Pattern) Segment(rate any) Pattern {
	r := int(toFloat(Reify(rate).FirstCycleValue()))
	if r <= 0 {
		r = 1
	}
	return p.Struct(Pure(true).FastF(FractionFromInt(int64(r))))
}

// Seg alias for Segment
func (p Pattern) Seg(rate any) Pattern { return p.Segment(rate) }

// SwingBy delays second half of each slice.
func (p Pattern) SwingBy(swing, n any) Pattern {
	s := toFloat(Reify(swing).FirstCycleValue())
	nn := Reify(n)
	seqPat := Sequence(Pure(0), Pure(s/2))
	return p.Inside(nn, func(pat Pattern) Pattern { return pat.Late(seqPat) })
}

// Swing shorthand (1/3 swing)
func (p Pattern) Swing(n any) Pattern { return p.SwingBy(1.0/3, n) }

// Filter keeps haps where test on hap true.
func (p Pattern) Filter(test func(Hap) bool) Pattern {
	return p.FilterHaps(test)
}

// FilterWhen keeps haps where test on whole begin true.
func (p Pattern) FilterWhen(test func(Fraction) bool) Pattern {
	return p.FilterHaps(func(h Hap) bool {
		if h.Whole == nil {
			return false
		}
		return test(h.Whole.Begin)
	})
}

// Within applies fn to haps within a..b cyclePos, keeps others.
func (p Pattern) Within(a, b float64, fn func(Pattern) Pattern) Pattern {
	return Stack(
		fn(p.FilterWhen(func(t Fraction) bool { return t.CyclePos().Float64() >= a && t.CyclePos().Float64() <= b })),
		p.FilterWhen(func(t Fraction) bool { return t.CyclePos().Float64() < a || t.CyclePos().Float64() > b }),
	)
}

// Invert swaps true/false in binary pattern
func (p Pattern) Invert() Pattern { return p.Fmap(func(v any) any { if bv, ok := v.(bool); ok { return !bv }; return v }) }
func (p Pattern) Inv() Pattern    { return p.Invert() }

// Brak breaks every other cycle
func (p Pattern) Brak() Pattern {
	return p.When(SlowCat(Pure(false), Pure(true)), func(x Pattern) Pattern { return FastCat(x, Silence()).Late(FractionFromFloat(0.25)) })
}

// FirstCycleValue already in pattern_time.go
