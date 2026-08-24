// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "math"

// --- ceil / floor aliases (JS register 'ceil'/'floor' -> Go CeilPat/FloorPat) ---
func (p Pattern) Ceil() Pattern { return p.CeilPat() }
func (p Pattern) Floor() Pattern { return p.FloorPat() }

// --- Beat (JS __beat -> _compress + innerJoin) ---
func (p Pattern) Beat(t, div any) Pattern {
	tFrac := toFrac(Reify(t).FirstCycleValue())
	dFrac := toFrac(Reify(div).FirstCycleValue())
	if dFrac.Equals(FractionFromInt(0)) {
		return p
	}
	b := tFrac.Mod(dFrac).Div(dFrac)
	e := tFrac.Add(FractionFromInt(1)).Div(dFrac)
	return p.Fmap(func(v any) any { return Pure(v).Compress(b, e) }).InnerJoin()
}
func Beat(t, div any, pat Pattern) Pattern { return pat.Beat(t, div) }

// helper toFrac
func toFrac(v any) Fraction {
	switch x := v.(type) {
	case Fraction:
		return x
	case *Fraction:
		if x != nil {
			return *x
		}
	case int:
		return FractionFromInt(int64(x))
	case int64:
		return FractionFromInt(x)
	case float64:
		return FractionFromFloat(x)
	case string:
		if f, err := ParseFraction(x); err == nil {
			return f
		}
	}
	return FractionFromFloat(toFloat(v))
}
// --- ZoomArc (JS zoomArc -> pat.zoom(a.begin,a.end)) ---
func (p Pattern) ZoomArc(a TimeSpan) Pattern { return p.Zoom(a.Begin, a.End) }
func (p Pattern) Zoomarc(a TimeSpan) Pattern { return p.ZoomArc(a) }

// --- LoopAtCps (JS _loopAt with cps param, simplified to LoopAt) ---
func (p Pattern) LoopAtCps(factor, cps any) Pattern { return p.LoopAt(factor) }
func (p Pattern) Loopatcps(factor, cps any) Pattern { return p.LoopAtCps(factor, cps) }

// --- StutWith / EchoWith aliases ---
func (p Pattern) StutWith(times, time any, fn func(Pattern, int) Pattern) Pattern {
	return p.EchoWith(times, time, fn)
}
func (p Pattern) Stutwith(times, time any, fn func(Pattern, int) Pattern) Pattern {
	return p.StutWith(times, time, fn)
}
func (p Pattern) Echowith(times, time any, fn func(Pattern, int) Pattern) Pattern {
	return p.EchoWith(times, time, fn)
}

// --- Flux / JuxFlip aliases (JS flux = juxFlip) ---
func (p Pattern) Flux(fn func(Pattern) Pattern) Pattern { return p.JuxFlip(fn) }
func (p Pattern) JuxFlip(fn func(Pattern) Pattern) Pattern { return p.Jux(fn) }
func (p Pattern) Juxflip(fn func(Pattern) Pattern) Pattern { return p.JuxFlip(fn) }
func (p Pattern) FluxBy(by float64, fn func(Pattern) Pattern) Pattern {
	return p.JuxFlipBy(by, fn)
}
func (p Pattern) JuxFlipBy(by float64, fn func(Pattern) Pattern) Pattern {
	return p.JuxBy(by, fn)
}
func (p Pattern) Fluxby(by float64, fn func(Pattern) Pattern) Pattern { return p.FluxBy(by, fn) }
func (p Pattern) Juxflipby(by float64, fn func(Pattern) Pattern) Pattern {
	return p.JuxFlipBy(by, fn)
}
// --- ChunkInto / ChunkBackInto (JS via into + iter) ---
func (p Pattern) ChunkInto(n int, fn func(Pattern) Pattern) Pattern { return p.Chunk(n, fn) }
func (p Pattern) Chunkinto(n int, fn func(Pattern) Pattern) Pattern { return p.ChunkInto(n, fn) }
func (p Pattern) ChunkBackInto(n int, fn func(Pattern) Pattern) Pattern {
	return p.ChunkBack(n, fn)
}
func (p Pattern) Chunkbackinto(n int, fn func(Pattern) Pattern) Pattern {
	return p.ChunkBackInto(n, fn)
}

// --- ApplyN (JS applyN: loop func n times) ---
func (p Pattern) ApplyN(n int, fn func(Pattern) Pattern) Pattern {
	result := p
	for i := 0; i < n; i++ {
		result = fn(result)
	}
	return result
}
func ApplyN(n int, fn func(Pattern) Pattern, pat Pattern) Pattern {
	return pat.ApplyN(n, fn)
}

// --- Reset / Restart aliases (JS keepif.reset etc.) ---
func (p Pattern) Reset(args ...any) Pattern {
	// JS reset = keepif.reset — alignment-based; approximate as passthrough if not pattern-of-patterns
	// to avoid empty for Pure strings, return p
	if len(args) == 0 {
		return p
	}
	// If p is pattern-of-patterns, use ResetJoin; otherwise return p
	hasPat := false
	for _, h := range p.Query(NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil)) {
		if _, ok := h.Value.(Pattern); ok {
			hasPat = true
			break
		}
		if _, ok := h.Value.(*Pattern); ok {
			hasPat = true
			break
		}
	}
	if hasPat {
		return p.ResetJoin()
	}
	return p
}
func (p Pattern) ResetAll(args ...any) Pattern { return p.Reset(args...) }
func (p Pattern) Restart(args ...any) Pattern {
	if len(args) == 0 {
		return p
	}
	hasPat := false
	for _, h := range p.Query(NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil)) {
		if _, ok := h.Value.(Pattern); ok {
			hasPat = true
			break
		}
		if _, ok := h.Value.(*Pattern); ok {
			hasPat = true
			break
		}
	}
	if hasPat {
		return p.RestartJoin()
	}
	return p
}
func (p Pattern) RestartAll(args ...any) Pattern { return p.Restart(args...) }

// Ensure math import used
var _ = math.Ceil
