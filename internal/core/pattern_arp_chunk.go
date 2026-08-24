// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Additional Pattern methods: collect/arp, chunk family, iterBack, repeatCycles, shrink/grow stubs.

package core

import "fmt"

// Collect groups congruent haps (same whole+part) into Hap with []Hap value.
func (p Pattern) Collect() Pattern {
	return p.WithHaps(func(haps []Hap, state State) []Hap {
		if len(haps) == 0 {
			return haps
		}
		groups := [][]Hap{}
		for _, h := range haps {
			found := -1
			for i, g := range groups {
				if len(g) > 0 && h.Part.Equals(g[0].Part) {
					wholeEq := (h.Whole == nil && g[0].Whole == nil) || (h.Whole != nil && g[0].Whole != nil && h.Whole.Equals(*g[0].Whole))
					if wholeEq {
						found = i
						break
					}
				}
			}
			if found == -1 {
				groups = append(groups, []Hap{h})
			} else {
				groups[found] = append(groups[found], h)
			}
		}
		out := make([]Hap, 0, len(groups))
		for _, g := range groups {
			first := g[0]
			out = append(out, NewHap(first.Whole, first.Part, g, map[string]any{}))
		}
		return out
	})
}

// ArpWith selects indices via func on collected haps.
func (p Pattern) ArpWith(fn func([]Hap) Pattern) Pattern {
	return p.Collect().Fmap(func(v any) any {
		haps, _ := v.([]Hap)
		return Reify(fn(haps))
	}).InnerJoin().WithHap(func(h Hap) Hap {
		// h.value is Hap? Actually innerJoin yields hap with value = Hap's value? JS does h.value.value
		// Our Fmap returns Pattern, innerJoin flattens, then withHap extracts value.value
		if inner, ok := h.Value.(Hap); ok {
			return NewHap(h.Whole, h.Part, inner.Value, h.CombinedContext(inner))
		}
		// fallback: if value is []Hap etc
		return h
	})
}

// Arp selects indices pattern.
func (p Pattern) Arp(indices any) Pattern {
	ip := Reify(indices)
	return p.ArpWith(func(haps []Hap) Pattern {
		if len(haps) == 0 {
			return Silence()
		}
		return ip.Fmap(func(v any) any {
			i := int(toFloat(v))
			idx := Mod(i, len(haps))
			if idx < 0 {
				idx += len(haps)
			}
			return haps[idx]
		})
	})
}

// IterBack plays subdivisions in reverse (mirrors JS _iter with back=true)
// JS: iterBack n => chunk into n, reverse each chunk's order via Rev, repeatCycles
func (p Pattern) IterBack(n int) Pattern {
	if n <= 0 {
		return p
	}
	// Implement via ChunkBack with Rev: chunk into n, reverse, then repeatCycles
	return p.ChunkBack(n, func(pat Pattern) Pattern { return pat.Rev() })
}

// ChunkBack etc (Chunk/RepeatCycles already in pattern_random.go)
func (p Pattern) SlowChunk(n int, fn func(Pattern) Pattern) Pattern { return p.Chunk(n, fn) }
func (p Pattern) ChunkBack(n int, fn func(Pattern) Pattern) Pattern {
	if n <= 0 {
		return p
	}
	return p.When(pulsedBinary(n, true), func(pat Pattern) Pattern { return fn(pat) }).RepeatCycles(n)
}
func (p Pattern) FastChunk(n int, fn func(Pattern) Pattern) Pattern {
	if n <= 0 {
		return p
	}
	// fast chunk without repeatCycles
	binary := pulsedBinary(n, false)
	return p.When(binary, fn)
}

// pulsedBinary helper already in pattern_random.go; re-expose wrapper for chunk file
func pulsedBinaryChunk(n int, back bool) Pattern { return pulsedBinary(n, back) }

// Shrink/Grow stubs (stepwise). Full impl requires shrinklist which is complex; provide simplified.
func (p Pattern) Shrink(amount any) Pattern {
	if p.Steps == nil {
		return Silence()
	}
	// Simplified: compress by amount fraction
	amt := toFloat(Reify(amount).FirstCycleValue())
	if amt == 0 {
		return Silence()
	}
	f := FractionFromFloat(amt)
	// shrinklist not fully implemented; approximate via Fast
	return p.Fast(Pure(f))
}
func (p Pattern) Grow(amount any) Pattern {
	if p.Steps == nil {
		return Silence()
	}
	amt := toFloat(Reify(amount).FirstCycleValue())
	f := FractionFromFloat(amt)
	return p.Fast(Pure(FractionFromInt(1).Div(f)))
}

// Tour inserts patterns progressively (simplified as SlowCat of variants)
func (p Pattern) Tour(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return p
	}
	all := append([]Pattern{p}, pats...)
	return SlowCat(all...)
}

// Zip combines patterns stepwise (simplified as Stack)
func Zip(pats ...Pattern) Pattern { return Stack(pats...) }

// Hush already in pattern_ops.go (Silence)

// XFade cross-fades
func XFade(a Pattern, pos any, b Pattern) Pattern {
	posPat := Reify(pos)
	return a.Fmap(func(av any) any {
		return func(bv any) any {
			// pos determines gain: simplified ignore pos, just stack
			_ = bv
			return av
		}
	}).AppBoth(posPat).OuterJoin()
}

// Helper for WithHap with context merging
func (h Hap) CombinedContext(other Hap) map[string]any {
	out := map[string]any{}
	for k, v := range h.Context {
		out[k] = v
	}
	for k, v := range other.Context {
		out[k] = v
	}
	// locations
	if locs, ok := out["locations"]; ok {
		_ = locs
	}
	return out
}

// WithHap helper already exists? Ensure Hap has WithHap via Pattern.WithHap
// Additional helpers for stringify debugging
func fmtHap(h Hap) string { return fmt.Sprintf("%v:%v", h.Value, h.Part) }
