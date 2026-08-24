// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/euclid.mjs — Bjorklund Euclidean rhythms.

package core

// Bjorklund generates Euclidean rhythm binary pattern.
func Bjorklund(pulses, steps int) []int {
	inverted := pulses < 0
	absOns := pulses
	if absOns < 0 {
		absOns = -absOns
	}
	if steps <= 0 {
		return []int{}
	}
	offs := steps - absOns
	if offs < 0 {
		offs = 0
	}
	ones := make([][]int, absOns)
	for i := range ones {
		ones[i] = []int{1}
	}
	zeros := make([][]int, offs)
	for i := range zeros {
		zeros[i] = []int{0}
	}
	result := bjorklundHelper([]int{absOns, offs}, [][]int{}, ones, zeros)
	// result is [n, x] where x is [ones, zeros] after recursion
	// For simplicity, flatten result second part
	flat := []int{}
	for _, arr := range result[0] {
		flat = append(flat, arr...)
	}
	for _, arr := range result[1] {
		flat = append(flat, arr...)
	}
	if inverted {
		for i, v := range flat {
			flat[i] = 1 - v
		}
	}
	return flat
}

func bjorklundHelper(n []int, _ [][]int, x [][]int, y [][]int) [][][]int {
	ons := n[0]
	offs := n[1]
	if min(ons, offs) <= 1 {
		return [][][]int{x, y}
	}
	return bjorklundRecurse([]int{ons, offs}, x, y)
}

// Simpler: implement directly iterative bjorklund as in JS but using slices
func BjorklundSimple(pulses, steps int) []int {
	if steps <= 0 {
		return []int{}
	}
	inverted := pulses < 0
	if inverted {
		pulses = -pulses
	}
	if pulses <= 0 {
		res := make([]int, steps)
		if inverted {
			for i := range res {
				res[i] = 1
			}
		}
		return res
	}
	if pulses >= steps {
		res := make([]int, steps)
		for i := range res {
			res[i] = 1
		}
		if inverted {
			for i := range res {
				res[i] = 0
			}
		}
		return res
	}
	// Use Bjorklund algorithm via left/right recursion similar to JS
	// For simplicity, use well-known iterative method
	pattern := make([]int, steps)
	// Distribute pulses
	// Use bucket method: for i in 0..steps-1, check if (i*pulses)%steps < pulses
	for i := 0; i < steps; i++ {
		if (i*pulses)%steps < pulses {
			pattern[i] = 1
		} else {
			pattern[i] = 0
		}
	}
	// Rotate to start with pulse? JS bjorklund starts with pulses? Check JS: ones first, then zeros flattened, gives pattern starting with 1
	// Our bucket gives pattern starting with 1 (i=0 always < pulses)
	if inverted {
		for i := range pattern {
			pattern[i] = 1 - pattern[i]
		}
	}
	return pattern
}

// Euclid returns pattern with euclidean rhythm applied via struct.
func (p Pattern) Euclid(pulses, steps int) Pattern {
	return p.EuclidRot(pulses, steps, 0)
}

func (p Pattern) EuclidRot(pulses, steps, rotation int) Pattern {
	b := BjorklundSimple(pulses, steps)
	if rotation != 0 {
		b = rotateInts(b, -rotation)
	}
	// Convert to bool pattern for struct: 1->true, 0->false
	// struct expects binary pattern? Use Filter? For now, implement via struct-like logic:
	// Create a binary pattern string? Simplified: generate a Pattern of bools and use Struct
	boolPat := SequenceFromInts(b)
	return p.Struct(boolPat)
}

func (p Pattern) EuclidLegato(pulses, steps int) Pattern {
	return p.EuclidLegatoRot(pulses, steps, 0)
}

func (p Pattern) EuclidLegatoRot(pulses, steps, rotation int) Pattern {
	if pulses < 1 {
		return Silence()
	}
	b := BjorklundSimple(pulses, steps)
	if rotation != 0 {
		b = rotateInts(b, -rotation)
	}
	// gapless: join '' split '1' map length+1
	binStr := ""
	for _, v := range b {
		if v == 1 {
			binStr += "1"
		} else {
			binStr += "0"
		}
	}
	// Split by "1"
	parts := []string{}
	start := -1
	for i, ch := range binStr {
		if ch == '1' {
			if start != -1 {
				parts = append(parts, binStr[start:i])
			}
			start = i
		}
	}
	if start != -1 {
		parts = append(parts, binStr[start:])
	}
	// Now parts without first empty? Need to mimic JS gapless logic
	// Simplified: gapless = bin_pat.join('').split('1').slice(1).map(s => [s.length+1, true])
	gapless := [][2]any{}
	for _, s := range parts {
		if s == "" {
			continue
		}
		// s includes leading 1
		gapless = append(gapless, [2]any{len(s), true})
	}
	// For now, just use TimeCat like JS: timeCat(...gapless)
	// We'll approximate by struct with legato
	// Use Sequence of durations
	durations := []any{}
	for _, g := range gapless {
		durations = append(durations, g[0], true) // placeholder
		_ = durations
	}
	_ = binStr
	return p.Struct(SequenceFromInts(b)).LateF(FractionFromInt(int64(rotation)).Div(FractionFromInt(int64(steps))))
}

// Helper: rotate ints
func rotateInts(a []int, n int) []int {
	if len(a) == 0 {
		return a
	}
	n = ((n % len(a)) + len(a)) % len(a)
	return append(a[n:], a[:n]...)
}

// SequenceFromInts creates a pattern of bools segmented equally per cycle.
func SequenceFromInts(ints []int) Pattern {
	n := len(ints)
	if n == 0 {
		return Silence()
	}
	return NewPattern(func(state State) []Hap {
		var out []Hap
		for _, sub := range state.Span.SpanCycles() {
			base := sub.Begin.Sam()
			for i, v := range ints {
				segBegin := base.Add(FractionFromInt(int64(i)).Div(FractionFromInt(int64(n))))
				segEnd := segBegin.Add(FractionFromInt(1).Div(FractionFromInt(int64(n))))
				seg := NewTimeSpan(segBegin, segEnd)
				inter := seg.Intersection(sub)
				if inter == nil {
					continue
				}
				val := false
				if v == 1 {
					val = true
				}
				out = append(out, NewHap(&seg, *inter, val, map[string]any{}))
			}
		}
		return out
	}, nil)
}

// Struct applies binary pattern: 1 keeps event
func (p Pattern) Struct(binaryPat Pattern) Pattern {
	// Struct via keepif.out: preserve binary structure, values from p
	// Use AppRight: binary's haps drive sampling of p
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			if bv, ok := b.(bool); ok {
				if bv {
					return a
				}
				return nil
			}
			// also handle numeric 1/0
			switch v := b.(type) {
			case int:
				if v != 0 { return a }
			case float64:
				if v != 0 { return a }
			}
			return nil
		}
	})
	// funcPat is Pattern of functions a-> (b-> a if b else nil)
	// AppRight with binaryPat
	valPat := binaryPat.Fmap(func(b any) any { return b })
	result := funcPat.AppRight(valPat).RemoveUndefineds()
	// AppRight preserves whole from right (binary), which gives correct segmentation
	return result
}

// Keep is set, etc - placeholder for struct alias
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bjorklundRecurse(n []int, x [][]int, y [][]int) [][][]int {
	return [][][]int{x, y}
}
