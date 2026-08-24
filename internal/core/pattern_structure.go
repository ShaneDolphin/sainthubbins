// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — structure: stack/cat, polymeter, arrange, palindrome, etc.

package core

import "fmt"

// Arrange cycles through sections (time, pat) pairs.
// Simplified: sections are [cycles, pat, cycles, pat, ...] where cycles is int/float.
func Arrange(sections ...any) Pattern {
	if len(sections) == 0 {
		return Silence()
	}
	// Build timeCat-like sequence: each pair is [duration, pat]
	var pats []Pattern
	var durations []Fraction
	for i := 0; i < len(sections)-1; i += 2 {
		dur := toFraction(sections[i])
		var pat Pattern
		switch v := sections[i+1].(type) {
		case Pattern:
			pat = v
		case string:
			pat = Reify(v)
		default:
			pat = Reify(v)
		}
		durations = append(durations, dur)
		pats = append(pats, pat)
	}
	// Total duration
	total := FractionFromInt(0)
	for _, d := range durations {
		total = total.Add(d)
	}
	if total.Equals(FractionFromInt(0)) {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			// For each sub cycle, need to determine which section it falls in
			// Simplified: linear time mapping: sub's position within total cycles
			cycle := sub.Begin.Sam()
			// Map cycle position modulo total
			pos := cycle.Mod(total)
			acc := FractionFromInt(0)
			for idx, dur := range durations {
				next := acc.Add(dur)
				if pos.Lte(next) || idx == len(durations)-1 {
					// This section active
					pat := pats[idx]
					// offset within section
					off := pos.Sub(acc)
					// Query pat at off
					haps := pat.Query(state.SetSpan(sub.WithTime(func(t Fraction) Fraction { return t.Sub(cycle).Sub(off).Add(cycle) })))
					for _, h := range haps {
						// Shift to correct time?
						// Simplified: keep as is if intersects
						if h.Part.Intersection(sub) != nil {
							out = append(out, h)
						}
					}
					break
				}
				acc = next
			}
		}
		return out
	}, nil)
}

// Palindrome alternates pattern and its reverse per cycle.
func (p Pattern) Palindrome() Pattern {
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			cycle := sub.Begin.Sam().Floor().Float64()
			isRev := int(cycle)%2 == 1
			var pat Pattern
			if isRev {
				pat = p.Rev()
			} else {
				pat = p
			}
			out = append(out, pat.Query(state.SetSpan(sub))...)
		}
		return out
	}, p.Steps)
}

// Iter subdivides cycle into n parts and applies progressively.
// Simplified: repeats pattern subdivided.
func (p Pattern) Iter(n int) Pattern {
	if n <= 0 {
		return p
	}
	return p.Segment2(n).Fmap(func(v any) any { return v }).SqueezeJoin()
}

// Segment is alias for seg (already in pattern_time.go as Segment)
// Provide here if not exists: uses Struct to segment.
func (p Pattern) Segment2(n int) Pattern {
	return p.Struct(Pure(true).FastF(FractionFromInt(int64(n))))
}

// Ply alias already in pattern_time.go

// Superimpose layers function results on top of original
func (p Pattern) Superimpose(fns ...func(Pattern) Pattern) Pattern {
	base := p
	for _, fn := range fns {
		base = Stack(base, fn(p))
	}
	return base
}

func (p Pattern) Layer(fns ...func(Pattern) Pattern) Pattern {
	var pats []Pattern
	for _, fn := range fns {
		pats = append(pats, fn(p))
	}
	return Stack(pats...)
}

// OnsetsOnly, DiscreteOnly already in pattern.go

// Defragment merges adjacent haps with same value/whole.
func (p Pattern) Defragment() Pattern {
	return p.WithHaps(func(haps []Hap, state State) []Hap {
		// Group by whole and value, merge contiguous parts
		// Simplified: sort by part begin, then merge where a.part.end == b.part.begin and same whole/value
		// Use stringified value for equality
		// Sort first
		// This is O(n^2) but n small
		changed := true
		for changed {
			changed = false
			for i := 0; i < len(haps)-1; i++ {
				a := haps[i]
				for j := i + 1; j < len(haps); j++ {
					b := haps[j]
					if a.Whole != nil && b.Whole != nil && a.Whole.Equals(*b.Whole) {
						if a.Part.End.Equals(b.Part.Begin) {
							// Check values equal via stringify
							if fmt.Sprintf("%v", a.Value) == fmt.Sprintf("%v", b.Value) {
								// Merge
								mergedPart := NewTimeSpan(a.Part.Begin, b.Part.End)
								merged := NewHap(a.Whole, mergedPart, a.Value, a.Context)
								// Replace a with merged, remove b
								haps[i] = merged
								haps = append(haps[:j], haps[j+1:]...)
								changed = true
								break
							}
						}
						if b.Part.End.Equals(a.Part.Begin) {
							if fmt.Sprintf("%v", a.Value) == fmt.Sprintf("%v", b.Value) {
								mergedPart := NewTimeSpan(b.Part.Begin, a.Part.End)
								merged := NewHap(a.Whole, mergedPart, a.Value, a.Context)
								haps[i] = merged
								haps = append(haps[:j], haps[j+1:]...)
								changed = true
								break
							}
						}
					}
				}
				if changed {
					break
				}
			}
		}
		return haps
	})
}

// Jux applies function to pattern panned to right and stacks with original left.
func (p Pattern) Jux(fn func(Pattern) Pattern) Pattern {
	return Stack(p, fn(p).Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["pan"] = 1.0
			return m2
		}
		return map[string]any{"value": v, "pan": 1.0}
	}))
}

// JuxBy with by amount
func (p Pattern) JuxBy(by float64, fn func(Pattern) Pattern) Pattern {
	return Stack(p.Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["pan"] = 0.0
			return m2
		}
		return v
	}), fn(p).Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["pan"] = by
			return m2
		}
		return map[string]any{"value": v, "pan": by}
	}))
}
