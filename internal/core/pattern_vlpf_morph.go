// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Remaining small Pattern methods: vlpf, xfade, morph, degrade helpers, etc.

package core

import "math"

// Vlpf is a filter alias (superdough v.lpf) — sets lpf freq
// JS: register('vlpf', (freq, pat) => pat.fmap((v) => ({...v, cutoff: freq * (v.velocity ?? 1) })))
func (p Pattern) Vlpf(freq any) Pattern {
	fp := Reify(freq)
	return p.Fmap(func(pv any) any {
		return func(fv any) any {
			f := toFloat(fv)
			if m, ok := pv.(map[string]any); ok {
				vel := 1.0
				if vv, ok := m["velocity"]; ok {
					vel = toFloat(vv)
				} else if vv, ok := m["vel"]; ok {
					vel = toFloat(vv)
				}
				cutoff := f * vel
				m2 := map[string]any{}
				for k, vv := range m {
					m2[k] = vv
				}
				m2["cutoff"] = cutoff
				m2["lpf"] = cutoff
				m2["ctf"] = cutoff
				return m2
			}
			return map[string]any{"value": pv, "cutoff": f, "lpf": f, "ctf": f}
		}
	}).AppBoth(fp)
}

func fadeGain(p float64) float64 {
	if p < 0.5 {
		return 1
	}
	return 1 - (p-0.5)/0.5
}

// XFade cross-fades between two patterns by position (0..1)
// JS: let gaina = pos.fmap((v) => ({ gain: fadeGain(v) })); let gainb = pos.fmap((v) => ({ gain: fadeGain(1 - v) })); return stack(a.mul(gaina), b.mul(gainb));
func XFade2(a Pattern, pos any, b Pattern) Pattern {
	posPat := Reify(pos)
	gainAPat := posPat.Fmap(func(v any) any {
		return map[string]any{"gain": fadeGain(toFloat(v))}
	})
	gainBPat := posPat.Fmap(func(v any) any {
		return map[string]any{"gain": fadeGain(1 - toFloat(v))}
	})
	aScaled := a.Fmap(func(av any) any {
		return func(gv any) any {
			gain := 1.0
			if gm, ok := gv.(map[string]any); ok {
				if gg, ok := gm["gain"]; ok {
					gain = toFloat(gg)
				}
			} else {
				gain = toFloat(gv)
			}
			if am, ok := av.(map[string]any); ok {
				m2 := map[string]any{}
				for k, vv := range am {
					m2[k] = vv
				}
				if ag, ok := am["gain"]; ok {
					m2["gain"] = toFloat(ag) * gain
				} else {
					m2["gain"] = gain
				}
				return m2
			}
			return map[string]any{"value": av, "gain": gain}
		}
	}).AppBoth(gainAPat)
	bScaled := b.Fmap(func(bv any) any {
		return func(gv any) any {
			gain := 1.0
			if gm, ok := gv.(map[string]any); ok {
				if gg, ok := gm["gain"]; ok {
					gain = toFloat(gg)
				}
			} else {
				gain = toFloat(gv)
			}
			if bm, ok := bv.(map[string]any); ok {
				m2 := map[string]any{}
				for k, vv := range bm {
					m2[k] = vv
				}
				if bg, ok := bm["gain"]; ok {
					m2["gain"] = toFloat(bg) * gain
				} else {
					m2["gain"] = gain
				}
				return m2
			}
			return map[string]any{"value": bv, "gain": gain}
		}
	}).AppBoth(gainBPat)
	return Stack(aScaled, bScaled)
}
func (p Pattern) XFade(pos any, other Pattern) Pattern { return XFade2(p, pos, other) }

// Morph takes two binary lists and morphs by factor by (0..1)
// JS _morph: by in [0,1], dur=1/from.length, positions = [pos/len where value truthy], arcs = zipWith((posa,va),(posb,vb) => { b= by*(posb-posa)+posa; e=b+dur; TimeSpan(b,e) }, positions(from), positions(to))
func MorphList(from, to []int, by float64) Pattern {
	if by < 0 {
		by = 0
	}
	if by > 1 {
		by = 1
	}
	if len(from) == 0 || len(to) == 0 {
		return Silence()
	}
	// Build positions where value truthy (non-zero)
	type posVal struct {
		pos Fraction
		val int
	}
	buildPositions := func(list []int) []posVal {
		out := []posVal{}
		n := len(list)
		for i, v := range list {
			if v != 0 {
				out = append(out, posVal{pos: FractionFromInt(int64(i)).Div(FractionFromInt(int64(n))), val: v})
			}
		}
		return out
	}
	posFrom := buildPositions(from)
	posTo := buildPositions(to)
	if len(posFrom) == 0 || len(posTo) == 0 {
		return Silence()
	}
	// If lengths differ, zip to min length (JS zipWith does min)
	m := len(posFrom)
	if len(posTo) < m {
		m = len(posTo)
	}
	dur := FractionFromInt(1).Div(FractionFromInt(int64(len(from))))
	byFrac := FractionFromFloat(by)
	arcs := make([]TimeSpan, 0, m)
	for i := 0; i < m; i++ {
		posa := posFrom[i].pos
		posb := posTo[i].pos
		// b = by*(posb-posa)+posa
		delta := posb.Sub(posa)
		b := byFrac.Mul(delta).Add(posa)
		e := b.Add(dur)
		arcs = append(arcs, NewTimeSpan(b, e))
	}
	query := func(state State) []Hap {
		cycle := state.Span.Begin.Sam()
		cycleArc := state.Span.CycleArc()
		out := []Hap{}
		for _, whole := range arcs {
			part := whole.Intersection(cycleArc)
			if part == nil {
				continue
			}
			// shift into cycle
			wholeShifted := whole.WithTime(func(f Fraction) Fraction { return f.Add(cycle) })
			partShifted := part.WithTime(func(f Fraction) Fraction { return f.Add(cycle) })
			out = append(out, NewHap(&wholeShifted, partShifted, true, map[string]any{}))
		}
		return out
	}
	return NewPattern(query, nil).SplitQueries()
}

// UndegradeBy helper (Degrade already in pattern_random.go)
func (p Pattern) UndegradeBy(prob float64) Pattern {
	return p.Fmap(func(v any) any { return v }).WithHaps(func(haps []Hap, s State) []Hap {
		return haps
	})
}

// Sometimes etc already exist as SometimesBy

// Helper for SequenceFromInts (used by morph)
func SequenceFromInts2(ints []int) Pattern { return SequenceFromInts(ints) }

// Additional distortion shims already in pattern_extra4.go; ensure vlpf registered
var _ = math.Log // keep math import used

// Ensure math import kept
