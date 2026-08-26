// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs (4191 LOC) — core query engine.

package core

import (
	"fmt"
	"sort"
	"sync"
)

// Global steps flag (mirrors JS __steps) and string parser hook.
var stepsEnabled = true

// stringParser and stringParserMu guard the hook set by SetStringParser.
// Multiple goroutines evaluate patterns concurrently (the web console's
// /api/evaluate calls jsapi.Evaluate per request, each of which registers
// this hook), so both the write in SetStringParser and the read in Reify
// must go through the same mutex.
var (
	stringParserMu sync.RWMutex
	stringParser   func(string) Pattern
)

func CalculateSteps(v bool) { stepsEnabled = v }

func SetStringParser(p func(string) Pattern) {
	stringParserMu.Lock()
	stringParser = p
	stringParserMu.Unlock()
}

// getStringParser returns the currently registered hook, or nil. Every read
// of stringParser (Reify below, and evaluate.go) must go through this rather
// than touching the package var directly.
func getStringParser() func(string) Pattern {
	stringParserMu.RLock()
	defer stringParserMu.RUnlock()
	return stringParser
}

// Pattern is the core type: a function from State to []Hap, with optional steps metadata.
type Pattern struct {
	Query     func(State) []Hap
	Steps     *Fraction
	PureValue any
	pureLoc   any
}

// NewPattern creates a Pattern.
func NewPattern(query func(State) []Hap, steps ...*Fraction) Pattern {
	var s *Fraction
	if len(steps) > 0 {
		s = steps[0]
	}
	return Pattern{Query: query, Steps: s}
}

func (p Pattern) SetSteps(steps Fraction) Pattern {
	stepsCopy := steps
	p.Steps = &stepsCopy
	return p
}

func (p Pattern) WithSteps(fn func(Fraction) Fraction) Pattern {
	if !stepsEnabled {
		return p
	}
	if p.Steps == nil {
		return NewPattern(p.Query, nil)
	}
	newSteps := fn(*p.Steps)
	return NewPattern(p.Query, &newSteps)
}

func (p Pattern) HasSteps() bool { return p.Steps != nil }

// Pure creates a pattern that yields the given value for every cycle, split at cycle boundaries.
func Pure(value any) Pattern {
	query := func(state State) []Hap {
		spans := state.Span.SpanCycles()
		haps := make([]Hap, 0, len(spans))
		for _, subspan := range spans {
			whole := NewTimeSpan(subspan.Begin.Sam(), subspan.Begin.NextSam())
			haps = append(haps, NewHap(&whole, subspan, value, map[string]any{}))
		}
		return haps
	}
	result := NewPattern(query, func() *Fraction { f := FractionFromInt(1); return &f }())
	result.PureValue = value
	return result
}

func Silence() Pattern {
	return NewPattern(func(State) []Hap { return []Hap{} }, nil)
}

func IsPattern(v any) bool {
	_, ok := v.(Pattern)
	if ok {
		return true
	}
	_, ok = v.(*Pattern)
	return ok
}

func Reify(thing any) Pattern {
	if p, ok := thing.(Pattern); ok {
		return p
	}
	if p, ok := thing.(*Pattern); ok {
		if p != nil {
			return *p
		}
	}
	if s, ok := thing.(string); ok {
		if p := getStringParser(); p != nil {
			return p(s)
		}
	}
	return Pure(thing)
}

// ---- functor / applicative / monad ----

func (p Pattern) WithValue(fn func(any) any) Pattern {
	query := func(state State) []Hap {
		haps := p.Query(state)
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = h.WithValue(fn)
		}
		return out
	}
	result := NewPattern(query, p.Steps)
	result.PureValue = p.PureValue
	return result
}

func (p Pattern) Fmap(fn func(any) any) Pattern { return p.WithValue(fn) }

func (p Pattern) WithState(fn func(State) State) Pattern {
	return NewPattern(func(state State) []Hap { return p.Query(fn(state)) }, p.Steps)
}

func (p Pattern) AppWhole(wholeFunc func(*TimeSpan, *TimeSpan) *TimeSpan, patVal Pattern) Pattern {
	patFunc := p
	query := func(state State) []Hap {
		hapFuncs := patFunc.Query(state)
		hapVals := patVal.Query(state)
		var out []Hap
		for _, hf := range hapFuncs {
			for _, hv := range hapVals {
				s := hf.Part.Intersection(hv.Part)
				if s == nil {
					continue
				}
				whole := wholeFunc(hf.Whole, hv.Whole)
				// hf.Value is expected to be func(any) any
				fn, ok := hf.Value.(func(any) any)
				if !ok {
					continue
				}
				newVal := fn(hv.Value)
				newCtx := hv.CombineContext(hf.Context)
				out = append(out, NewHap(whole, *s, newVal, newCtx))
			}
		}
		return out
	}
	return NewPattern(query, nil)
}

func (p Pattern) AppBoth(patVal Pattern) Pattern {
	wholeFunc := func(a, b *TimeSpan) *TimeSpan {
		if a == nil || b == nil {
			return nil
		}
		inter := a.Intersection(*b)
		return inter
	}
	result := p.AppWhole(wholeFunc, patVal)
	if stepsEnabled {
		result.Steps = Lcm(p.Steps, patVal.Steps)
	}
	return result
}

func (p Pattern) AppLeft(patVal Pattern) Pattern {
	patFunc := p
	query := func(state State) []Hap {
		var haps []Hap
		for _, hf := range patFunc.Query(state) {
			wholeOrPart := hf.WholeOrPart()
			innerState := state.SetSpan(wholeOrPart)
			hapVals := patVal.Query(innerState)
			for _, hv := range hapVals {
				newPart := hf.Part.Intersection(hv.Part)
				if newPart == nil {
					continue
				}
				fn, ok := hf.Value.(func(any) any)
				if !ok {
					continue
				}
				newVal := fn(hv.Value)
				newCtx := hv.CombineContext(hf.Context)
				haps = append(haps, NewHap(hf.Whole, *newPart, newVal, newCtx))
			}
		}
		return haps
	}
	result := NewPattern(query, p.Steps)
	return result
}

func (p Pattern) AppRight(patVal Pattern) Pattern {
	patFunc := p
	query := func(state State) []Hap {
		var haps []Hap
		for _, hv := range patVal.Query(state) {
			wholeOrPart := hv.WholeOrPart()
			innerState := state.SetSpan(wholeOrPart)
			hapFuncs := patFunc.Query(innerState)
			for _, hf := range hapFuncs {
				newPart := hf.Part.Intersection(hv.Part)
				if newPart == nil {
					continue
				}
				fn, ok := hf.Value.(func(any) any)
				if !ok {
					continue
				}
				newVal := fn(hv.Value)
				newCtx := hv.CombineContext(hf.Context)
				haps = append(haps, NewHap(hv.Whole, *newPart, newVal, newCtx))
			}
		}
		return haps
	}
	result := NewPattern(query, patVal.Steps)
	return result
}

func (p Pattern) BindWhole(chooseWhole func(*TimeSpan, *TimeSpan) *TimeSpan, fn func(any) Pattern) Pattern {
	patVal := p
	query := func(state State) []Hap {
		var out []Hap
		for _, a := range patVal.Query(state) {
			innerPat := fn(a.Value)
			innerHaps := innerPat.Query(state.SetSpan(a.Part))
			for _, b := range innerHaps {
				whole := chooseWhole(a.Whole, b.Whole)
				part := b.Part
				newCtx := a.CombineContext(b.Context)
				out = append(out, NewHap(whole, part, b.Value, newCtx))
			}
		}
		return out
	}
	return NewPattern(query, nil)
}

func (p Pattern) Bind(fn func(any) Pattern) Pattern {
	wholeFunc := func(a, b *TimeSpan) *TimeSpan {
		if a == nil || b == nil {
			return nil
		}
		inter := a.Intersection(*b)
		return inter
	}
	return p.BindWhole(wholeFunc, fn)
}

func (p Pattern) Join() Pattern { return p.Bind(func(v any) Pattern {
	if pat, ok := v.(Pattern); ok {
		return pat
	}
	if pat, ok := v.(*Pattern); ok && pat != nil {
		return *pat
	}
	return Pure(v)
}) }

func (p Pattern) OuterBind(fn func(any) Pattern) Pattern {
	pat := p.BindWhole(func(a *TimeSpan, _ *TimeSpan) *TimeSpan { return a }, fn)
	pat.Steps = p.Steps
	return pat
}
func (p Pattern) OuterJoin() Pattern { return p.OuterBind(func(v any) Pattern {
	if pat, ok := v.(Pattern); ok {
		return pat
	}
	return Pure(v)
}) }
func (p Pattern) InnerBind(fn func(any) Pattern) Pattern {
	return p.BindWhole(func(_ *TimeSpan, b *TimeSpan) *TimeSpan { return b }, fn)
}
func (p Pattern) InnerJoin() Pattern { return p.InnerBind(func(v any) Pattern {
	if pat, ok := v.(Pattern); ok {
		return pat
	}
	return Pure(v)
}) }

// SqueezeJoin — like join but squeezes inner pattern into outer hap's duration.
func (p Pattern) SqueezeJoin() Pattern {
	pat := p
	query := func(state State) []Hap {
		var out []Hap
		for _, outer :=  range pat.Query(state) {
			innerPat, ok := outer.Value.(Pattern)
			if !ok {
				if pp, ok2 := outer.Value.(*Pattern); ok2 && pp != nil {
					innerPat = *pp
				} else {
					continue
				}
			}
			if outer.Whole == nil {
				continue
			}
			innerHaps := innerPat.Query(state.SetSpan(NewTimeSpan(FractionFromInt(0), FractionFromInt(1))))
			for _, inner := range innerHaps {
				// Map inner hap's span from [0,1) into outer.Part
					// Simple squeeze: scale inner part/whole proportionally into outer.Part
				// Use withCycle logic
				// For now, approximate by compressing cycle 0-1 into outer's whole
				if inner.Whole == nil {
					continue
				}
				// Convert inner times (0-1) to outer's time
				scale := func(f Fraction) Fraction {
					// f is in [0,1) or similar; map to outer.Whole.Begin + f * duration
					dur := outer.Whole.Duration()
					return outer.Whole.Begin.Add(f.Mul(dur))
				}
				newWhole := inner.Whole.WithTime(scale)
				newPart := inner.Part.WithTime(scale)
				// Intersect with outer.Part
				inter := newPart.Intersection(outer.Part)
				if inter == nil {
					continue
				}
				wholeInter := newWhole.Intersection(*outer.Whole)
				if wholeInter == nil {
					continue
				}
				newCtx := inner.CombineContext(outer.Context)
				out = append(out, NewHap(wholeInter, *inter, inner.Value, newCtx))
			}
		}
		return out
	}
	return NewPattern(query, p.Steps)
}

func (p Pattern) SqueezeBind(fn func(any) Pattern) Pattern {
	return p.Fmap(func(v any) any { return fn(v) }).SqueezeJoin()
}

// Query helpers
func (p Pattern) QueryArc(begin, end Fraction, controls ...map[string]any) []Hap {
	ctrl := map[string]any{}
	if len(controls) > 0 && controls[0] != nil {
		ctrl = controls[0]
	}
	state := NewState(NewTimeSpan(begin, end), ctrl)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("query panic: %v\n", r)
		}
	}()
	return p.Query(state)
}

func (p Pattern) SplitQueries() Pattern {
	pat := p
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			out = append(out, pat.Query(state.SetSpan(sub))...)
		}
		return out
	}, p.Steps)
}

func (p Pattern) WithQuerySpan(fn func(TimeSpan) TimeSpan) Pattern {
	return NewPattern(func(state State) []Hap { return p.Query(state.WithSpan(fn)) }, p.Steps)
}

func (p Pattern) WithQueryTime(fn func(Fraction) Fraction) Pattern {
	return NewPattern(func(state State) []Hap {
		return p.Query(state.WithSpan(func(s TimeSpan) TimeSpan { return s.WithTime(fn) }))
	}, p.Steps)
}

func (p Pattern) WithHapSpan(fn func(TimeSpan) TimeSpan) Pattern {
	return NewPattern(func(state State) []Hap {
		haps := p.Query(state)
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = h.WithSpan(fn)
		}
		return out
	}, p.Steps)
}

func (p Pattern) WithHapTime(fn func(Fraction) Fraction) Pattern {
	return p.WithHapSpan(func(s TimeSpan) TimeSpan { return s.WithTime(fn) })
}

func (p Pattern) WithHaps(fn func([]Hap, State) []Hap) Pattern {
	result := NewPattern(func(state State) []Hap { return fn(p.Query(state), state) }, p.Steps)
	return result
}

func (p Pattern) WithHap(fn func(Hap) Hap) Pattern {
	return p.WithHaps(func(haps []Hap, _ State) []Hap {
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = fn(h)
		}
		return out
	})
}

func (p Pattern) SetContext(ctx map[string]any) Pattern {
	return p.WithHap(func(h Hap) Hap { return h.SetContext(ctx) })
}

func (p Pattern) WithContext(fn func(map[string]any) map[string]any) Pattern {
	result := p.WithHap(func(h Hap) Hap { return h.SetContext(fn(h.Context)) })
	result.PureValue = p.PureValue
	return result
}

func (p Pattern) StripContext() Pattern { return p.WithHap(func(h Hap) Hap { return h.SetContext(map[string]any{}) }) }

func (p Pattern) WithLoc(start, end int) Pattern {
	loc := map[string]any{"start": start, "end": end}
	result := p.WithContext(func(ctx map[string]any) map[string]any {
		newCtx := map[string]any{}
		for k, v := range ctx {
			newCtx[k] = v
		}
		var locs []any
		if existing, ok := ctx["locations"]; ok {
			if sl, ok := existing.([]any); ok {
				locs = append(locs, sl...)
			}
		}
		locs = append(locs, loc)
		newCtx["locations"] = locs
		return newCtx
	})
	result.PureValue = p.PureValue
	return result
}

func (p Pattern) FilterHaps(test func(Hap) bool) Pattern {
	return NewPattern(func(state State) []Hap {
		haps := p.Query(state)
		var out []Hap
		for _, h := range haps {
			if test(h) {
				out = append(out, h)
			}
		}
		return out
	}, p.Steps)
}

func (p Pattern) FilterValues(test func(any) bool) Pattern {
	return NewPattern(func(state State) []Hap {
		haps := p.Query(state)
		var out []Hap
		for _, h := range haps {
			if test(h.Value) {
				out = append(out, h)
			}
		}
		return out
	}, p.Steps)
}

func (p Pattern) RemoveUndefineds() Pattern {
	return p.FilterValues(func(v any) bool { return v != nil })
}

func (p Pattern) OnsetsOnly() Pattern { return p.FilterHaps(func(h Hap) bool { return h.HasOnset() }) }
func (p Pattern) DiscreteOnly() Pattern { return p.FilterHaps(func(h Hap) bool { return h.Whole != nil }) }

func (p Pattern) FirstCycle(withContext ...bool) []Hap {
	withCtx := false
	if len(withContext) > 0 {
		withCtx = withContext[0]
	}
	pat := p
	if !withCtx {
		pat = pat.StripContext()
	}
	return pat.Query(NewState(NewTimeSpan(FractionFromInt(0), FractionFromInt(1)), nil))
}

func (p Pattern) SortHapsByPart() Pattern {
	return p.WithHaps(func(haps []Hap, _ State) []Hap {
		sorted := make([]Hap, len(haps))
		copy(sorted, haps)
		sort.Slice(sorted, func(i, j int) bool {
			a, b := sorted[i], sorted[j]
			if !a.Part.Begin.Equals(b.Part.Begin) {
				return a.Part.Begin.Lt(b.Part.Begin)
			}
			if !a.Part.End.Equals(b.Part.End) {
				return a.Part.End.Lt(b.Part.End)
			}
			if a.Whole != nil && b.Whole != nil {
				if !a.Whole.Begin.Equals(b.Whole.Begin) {
					return a.Whole.Begin.Lt(b.Whole.Begin)
				}
				return a.Whole.End.Lt(b.Whole.End)
			}
			return false
		})
		return sorted
	})
}

// Fast applies pattern of factors to pattern (structure from left).
func (p Pattern) Fast(factor Pattern) Pattern {
	// Simplified: sample factor's first value for now; full patternified version uses AppLeft
	// For single-value factor patterns (common case Pure), this is exact.
	// TODO: full AppLeft/Join patternified fast
	f := FractionFromInt(1)
	if haps := factor.FirstCycle(); len(haps) > 0 {
		f = toFraction(haps[0].Value)
	}
	return p.FastF(f)
}

// FastF speeds up pattern by factor frac (like JS fast).
func (p Pattern) FastF(frac Fraction) Pattern {
	return NewPattern(func(state State) []Hap {
		haps := p.Query(state.WithSpan(func(s TimeSpan) TimeSpan {
			return s.WithTime(func(t Fraction) Fraction { return t.Mul(frac) })
		}))
		out := make([]Hap, len(haps))
		for i, h := range haps {
			out[i] = h.WithSpan(func(s TimeSpan) TimeSpan {
				return s.WithTime(func(t Fraction) Fraction { return t.Div(frac) })
			})
		}
		return out
	}, p.Steps)
}

// Helper for slice map
func mapHapsSlice(haps []Hap, fn func(Hap) Hap) []Hap {
	out := make([]Hap, len(haps))
	for i, h := range haps {
		out[i] = fn(h)
	}
	return out
}

func (p Pattern) Slow(frac Fraction) Pattern { return p.FastF(FractionFromInt(1).Div(frac)) }

func toFraction(v any) Fraction {
	switch x := v.(type) {
	case Fraction:
		return x
	case *Fraction:
		return *x
	case int:
		return FractionFromInt(int64(x))
	case int64:
		return FractionFromInt(x)
	case float64:
		return FractionFromFloat(x)
	case string:
		f, err := ParseFraction(x)
		if err == nil {
			return f
		}
		return FractionFromInt(1)
	default:
		return FractionFromInt(1)
	}
}

// Extend Hap slice helper
func (p Pattern) queryWithHapMap(state State, fn func(Hap) Hap) []Hap {
	haps := p.Query(state)
	out := make([]Hap, len(haps))
	for i, h := range haps {
		out[i] = fn(h)
	}
	return out
}

// Add MapHaps method to []Hap via helper function
func mapHaps(haps []Hap, fn func(Hap) Hap) []Hap {
	out := make([]Hap, len(haps))
	for i, h := range haps {
		out[i] = fn(h)
	}
	return out
}

// Stack, Cat, FastCat, SlowCat
func Stack(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, pat := range pats {
			out = append(out, pat.Query(state)...)
		}
		return out
	}, nil)
}

func Cat(pats ...Pattern) Pattern { return SlowCat(pats...) }

func FastCat(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	n := len(pats)
	nFrac := FractionFromInt(int64(n))
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, cyc := range state.Span.SpanCycles() {
			base := cyc.Begin.Sam()
			for i, pat := range pats {
				segBegin := base.Add(FractionFromInt(int64(i)).Div(nFrac))
				segEnd := segBegin.Add(FractionFromInt(1).Div(nFrac))
				seg := NewTimeSpan(segBegin, segEnd)
				inter := seg.Intersection(cyc)
				if inter == nil {
					continue
				}
				// Map inter to pat's time: (t - segBegin) * n
				mappedSpan := inter.WithTime(func(t Fraction) Fraction {
					return t.Sub(segBegin).Mul(nFrac)
				})
				haps := pat.Query(state.SetSpan(mappedSpan))
				for _, h := range haps {
					newHap := h.WithSpan(func(s TimeSpan) TimeSpan {
						return s.WithTime(func(t Fraction) Fraction {
							return t.Div(nFrac).Add(segBegin)
						})
					})
					if newHap.Part.Intersection(cyc) != nil {
						out = append(out, newHap)
					}
				}
			}
		}
		return out
	}, nil)
}

// TimeCat helper
func TimeCat(args ...any) Pattern {
	// If args length even and first is duration-like, use weighted
	if len(args) >=2 && len(args)%2==0 {
		hasDur := false
		if _, ok := args[0].(Fraction); ok { hasDur=true }
		if _, ok := args[0].(*Fraction); ok { hasDur=true }
		if _, ok := args[0].(int); ok { hasDur=true }
		if _, ok := args[0].(float64); ok { hasDur=true }
		if hasDur {
			return TimeCatWeighted(args...)
		}
	}
	return NewPattern(func(state State) []Hap {
		return Stack(toPatterns(args)...).Query(state)
	}, nil)
}

func toPatterns(args []any) []Pattern {
	var out []Pattern
	for _, a := range args {
		if p, ok := a.(Pattern); ok {
			out = append(out, p)
		}
	}
	return out
}

func toTimeCatArgs(pats []Pattern, n Fraction) []any {
	var out []any
	for _, p := range pats {
		out = append(out, FractionFromInt(1), p)
	}
	return out
}

func SlowCat(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	// SlowCat cycles through pats, each for one cycle
	return NewPattern(func(state State) []Hap {
		// Which pat is active for the queried span?
		// For a query spanning multiple cycles, need to handle each subspan
		var out []Hap
		cycles := state.Span.SpanCycles()
		for _, cyc := range cycles {
			cycleIdx := cyc.Begin.Floor().Float64() // cycle number
			patIdx := Mod(int(cycleIdx), len(pats))
			pat := pats[patIdx]
			// Map cyc (which is within one cycle) to pat's cycle time
			// Pat's time is relative to cycle
			// The whole span is mapped relative to the start of this cycle.
			// Using each instant's own Sam() collapses the span, because the
			// end of cycle n has Sam() == n+1 and maps to zero.
			base := cyc.Begin.Sam()
			mappedSpan := cyc.WithTime(func(t Fraction) Fraction {
				return t.Sub(base)
			})
			haps := pat.Query(state.SetSpan(mappedSpan))
			// Map haps back
			for _, h := range haps {
				// Shift whole/part by cycle base
				shifted := h.WithSpan(func(s TimeSpan) TimeSpan {
					return s.WithTime(func(t Fraction) Fraction { return t.Add(base) })
				})
				out = append(out, shifted)
			}
		}
		return out
	}, nil)
}

func Sequence(pats ...Pattern) Pattern { return FastCat(pats...) }

// Polymeter repeats patterns to fit LCM of steps (mirrors JS polymeter).
func Polymeter(pats ...Pattern) Pattern {
	if len(pats) == 0 {
		return Silence()
	}
	var withSteps []Pattern
	for _, p := range pats {
		if p.Steps != nil && !p.Steps.Equals(FractionFromInt(0)) {
			withSteps = append(withSteps, p)
		}
	}
	if len(withSteps) == 0 {
		return Silence()
	}
	var steps []*Fraction
	for _, p := range withSteps {
		steps = append(steps, p.Steps)
	}
	lcm := Lcm(steps...)
	if lcm == nil || lcm.Equals(FractionFromInt(0)) {
		return Stack(withSteps...)
	}
	var expanded []Pattern
	for _, p := range withSteps {
		factor := lcm.Div(*p.Steps)
		expanded = append(expanded, p.SlowF(factor))
	}
	result := Stack(expanded...)
	result.Steps = lcm
	return result
}

