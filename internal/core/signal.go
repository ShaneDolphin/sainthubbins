// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/signal.mjs (1099 LOC) — continuous patterns.

package core

import (
	"fmt"
	"math"
)

// Signal creates a continuous pattern from a function Fraction -> float64.
// Whole is nil, value sampled at Part midpoint, like JS.
func Signal(fn func(Fraction) float64) Pattern {
	query := func(state State) []Hap {
		// For continuous signals, we sample at midpoint of query span (or subspans split at cycles)
		// JS splits at cycles then samples midpoint of each subspan
		subspans := state.Span.SpanCycles()
		var haps []Hap
		for _, sub := range subspans {
			mid := sub.Midpoint()
			val := fn(mid)
			// continuous: whole=nil, part=sub
			haps = append(haps, NewHap(nil, sub, val, map[string]any{}))
		}
		return haps
	}
	return NewPattern(query, nil)
}

func Saw() Pattern { return Signal(func(t Fraction) float64 { return t.CyclePos().Float64() }) }
func Saw2() Pattern { return Saw() }
func Isaw() Pattern { return Signal(func(t Fraction) float64 { return 1 - t.CyclePos().Float64() }) }
func Isaw2() Pattern { return Isaw() }
func Square() Pattern { return Signal(func(t Fraction) float64 {
	if t.CyclePos().Float64() < 0.5 { return 1 } else { return 0 }
}) }
func Square2() Pattern { return Square() }
func Tri() Pattern { return Signal(func(t Fraction) float64 {
	pos := t.CyclePos().Float64()
	if pos < 0.5 { return pos * 2 } else { return 2 - pos*2 }
}) }
func Tri2() Pattern { return Tri() }
func Sine() Pattern { return Signal(func(t Fraction) float64 { return math.Sin(2*math.Pi*t.CyclePos().Float64()) }) }
func Sine2() Pattern { return Sine() }
func RandPat() Pattern { return Signal(func(t Fraction) float64 {
	// Simple deterministic pseudo-random based on cycle count
	n := t.Sam().Float64()
	// Use sin-based hash
	return math.Mod(math.Sin(n*127.1)*43758.5453, 1)
}) }

// Steady returns a pattern with constant value (alias for Pure but distinct from signal sampling)
func Steady(value float64) Pattern { return Pure(value) }

func AddPat(a, b Pattern) Pattern {
	return a.Fmap(func(av any) any {
		return func(bv any) any {
			af := toFloat(av)
			bf := toFloat(bv)
			return af + bf
		}
	}).AppBoth(b.Fmap(func(bv any) any { return bv })).Fmap(func(v any) any {
		// Actually AppBoth of func pattern and value pattern: need to handle differently
		// Simplified: use Bind for arithmetic
		return v
	})
	// placeholder - real impl would be a._opMix etc. For now, use simple hap-wise_add via query
}

func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case int32:
		return float64(x)
	case uint:
		return float64(x)
	case uint64:
		return float64(x)
	case Fraction:
		return x.Float64()
	case *Fraction:
		if x != nil {
			return x.Float64()
		}
		return 0
	case string:
		var f float64
		_, err := fmt.Sscanf(x, "%f", &f)
		if err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
}
