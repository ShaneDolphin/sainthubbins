// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — randomness, degrade, choose, chunk/iter.

package core

import (
	"hash/fnv"
	"math"
	"math/rand"
)

// deterministic rand for a given Fraction time (cycle pos + cycle)
func pseudoRand(t Fraction) float64 {
	// Use fnv hash of numerator/denominator for determinism
	h := fnv.New64a()
	// Write fraction as string
	s := t.Show()
	h.Write([]byte(s))
	v := h.Sum64()
	// Convert to 0..1
	return float64(v%1000000) / 1000000.0
}

// Rand returns a continuous random signal 0..1 (sampled per hap mid)
func Rand() Pattern {
	return Signal(func(t Fraction) float64 {
		// hash of cycle + pos
		return pseudoRand(t)
	})
}

// Choose picks randomly from list pattern.
func (p Pattern) Choose(list []any) Pattern {
	if len(list) == 0 {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		haps := p.Query(state)
		var out []Hap
		for _, h := range haps {
			// Use hap's whole begin as seed
			seed := h.WholeOrPart().Begin
			r := pseudoRand(seed)
			idx := int(math.Floor(r * float64(len(list))))
			if idx >= len(list) {
				idx = len(list) - 1
			}
			if idx < 0 {
				idx = 0
			}
			val := list[idx]
			// If val is Pattern, need to query? Simplify: if Pattern, join
			if pat, ok := val.(Pattern); ok {
				inner := pat.Query(state.SetSpan(h.Part))
				for _, ih := range inner {
					whole := ih.Whole
					if whole == nil {
						whole = h.Whole
					}
					out = append(out, NewHap(whole, ih.Part, ih.Value, ih.Context))
				}
			} else {
				out = append(out, NewHap(h.Whole, h.Part, val, h.Context))
			}
		}
		return out
	}, p.Steps)
}

// DegradeBy randomly drops haps with prob.
func (p Pattern) DegradeBy(prob float64) Pattern {
	return p.FilterHaps(func(h Hap) bool {
		t := h.WholeOrPart().Begin
		return pseudoRand(t) >= prob
	})
}

func (p Pattern) Degrade() Pattern { return p.DegradeBy(0.5) }

// SometimesBy applies func with prob
func (p Pattern) SometimesBy(prob float64, fn func(Pattern) Pattern) Pattern {
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			base := sub.Begin.Sam()
			// per-cycle decision
			if pseudoRand(base) < prob {
				out = append(out, fn(p).Query(state.SetSpan(sub))...)
			} else {
				out = append(out, p.Query(state.SetSpan(sub))...)
			}
		}
		return out
	}, p.Steps)
}

func (p Pattern) Sometimes(fn func(Pattern) Pattern) Pattern { return p.SometimesBy(0.5, fn) }

// RandPat already in signal.go; alias
func RandPat2() Pattern { return Rand() }

// Shuffle reorders events randomly per cycle
func (p Pattern) Shuffle(n int) Pattern {
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			haps := p.Query(state.SetSpan(sub))
			if len(haps) <= 1 {
				out = append(out, haps...)
				continue
			}
			// deterministic shuffle per cycle
			seed := sub.Begin.Sam().Num ^ sub.Begin.Sam().Den
			r := rand.New(rand.NewSource(int64(seed + 0x9e3779b9)))
			perm := r.Perm(len(haps))
			shuffled := make([]Hap, len(haps))
			for i, pi := range perm {
				shuffled[i] = haps[pi]
			}
			// Need to reassign parts to shuffled order? Simplified: just permute values
			// Keep original parts but shuffle values
			for i := range shuffled {
				origPart := haps[i].Part
				origWhole := haps[i].Whole
				out = append(out, NewHap(origWhole, origPart, shuffled[i].Value, shuffled[i].Context))
			}
		}
		return out
	}, p.Steps)
}

// Chunk applies func to chunked parts
func (p Pattern) Chunk(n int, fn func(Pattern) Pattern) Pattern {
	if n <= 0 {
		return p
	}
	return p.When(pulsedBinary(n, false), fn).RepeatCycles(n)
}

// RepeatCycles repeats each cycle n times
func (p Pattern) RepeatCycles(n int) Pattern {
	if n <= 0 {
		return Silence()
	}
	nFrac := FractionFromInt(int64(n))
	return NewPattern(func(state State) []Hap {
		cycle := state.Span.Begin.Sam()
		sourceCycle := cycle.Div(nFrac).Sam()
		delta := cycle.Sub(sourceCycle)
		mappedState := state.WithSpan(func(s TimeSpan) TimeSpan { return s.WithTime(func(t Fraction) Fraction { return t.Sub(delta) }) })
		haps := p.Query(mappedState)
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = h.WithSpan(func(s TimeSpan) TimeSpan { return s.WithTime(func(t Fraction) Fraction { return t.Add(delta) }) })
		}
		return out
	}, p.Steps).SplitQueries()
}

func pulsedBinary(n int, back bool) Pattern {
	binary := make([]bool, n)
	binary[0] = true
	// create sequence pattern
	ints := make([]int, n)
	for i := range ints {
		if binary[i] {
			ints[i] = 1
		}
	}
	seq := SequenceFromInts(ints)
	if back {
		seq = seq.Rev()
	}
	return seq
}
