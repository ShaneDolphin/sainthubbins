// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/pattern.mjs — COMPOSERS numeric ops and helpers.

package core

import "math"

// toFloat helper already in signal.go, but define here if not exists
// Note: toFloat defined in signal.go, reuse.

// Add adds pattern/value to this pattern (numeric). When either side is a
// control bag, the addition lands on the bag's numeric field instead of
// flattening it — see addValues in pattern_arith.go.
func (p Pattern) Add(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return addValues(a, b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Sub(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) - toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Mul(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) * toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Div(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			bf := toFloat(b)
			if bf == 0 {
				return 0.0
			}
			return toFloat(a) / bf
		}
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

// Range scales unipolar 0..1 pattern to min..max
func (p Pattern) Range(min, max float64) Pattern {
	diff := max - min
	return p.Mul(diff).Add(min)
}

func (p Pattern) Rangex(min, max float64) Pattern {
	// exponential: log scale
	logMin := math.Log(min)
	logMax := math.Log(max)
	// _range with log values then exp
	// For now, approximate: p.Mul(logMax-logMin).Add(logMin) then exp
	return p.Mul(logMax - logMin).Add(logMin).Fmap(func(v any) any { return math.Exp(toFloat(v)) })
}

func (p Pattern) Range2(min, max float64) Pattern {
	// bipolar -1..1 to min..max
	return p.FromBipolar().Range(min, max)
}

func (p Pattern) ToBipolar() Pattern {
	return p.Fmap(func(v any) any { return toFloat(v)*2 - 1 })
}

func (p Pattern) FromBipolar() Pattern {
	return p.Fmap(func(v any) any { return (toFloat(v) + 1) / 2 })
}

func (p Pattern) Mod(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return math.Mod(toFloat(a), toFloat(b)) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Pow(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return math.Pow(toFloat(a), toFloat(b)) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Band(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return int(toFloat(a)) & int(toFloat(b)) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Bor(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return int(toFloat(a)) | int(toFloat(b)) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Bxor(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return int(toFloat(a)) ^ int(toFloat(b)) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Blshift(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return int(toFloat(a)) << uint(int(toFloat(b))) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Brshift(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return int(toFloat(a)) >> uint(int(toFloat(b))) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Lt(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) < toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Gt(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) > toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Lte(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) <= toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Gte(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return toFloat(a) >= toFloat(b) }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Eq(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return a == b }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Eqt(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return a == b }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Ne(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return a != b }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Net(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any { return a != b }
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) And(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			ab, aok := a.(bool)
			bb, bok := b.(bool)
			if aok && bok {
				return ab && bb
			}
			af := toFloat(a)
			bf := toFloat(b)
			return af != 0 && bf != 0
		}
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Or(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			ab, aok := a.(bool)
			bb, bok := b.(bool)
			if aok && bok {
				return ab || bb
			}
			af := toFloat(a)
			bf := toFloat(b)
			return af != 0 || bf != 0
		}
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Func(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			if fn, ok := b.(func(any) any); ok {
				return fn(a)
			}
			if fn, ok := a.(func(any) any); ok {
				return fn(b)
			}
			return b
		}
	})
	return funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Round() Pattern {
	return p.Fmap(func(v any) any { return math.Round(toFloat(v)) })
}

func (p Pattern) FloorPat() Pattern {
	return p.Fmap(func(v any) any { return math.Floor(toFloat(v)) })
}

func (p Pattern) CeilPat() Pattern {
	return p.Fmap(func(v any) any { return math.Ceil(toFloat(v)) })
}

// Gap creates silence with steps (JS: gap(steps) -> Pattern with steps but no haps)
func (p Pattern) Gap(steps int) Pattern { return Gap(steps) }
func (p Pattern) SilenceP() Pattern    { return Silence() }
func (p Pattern) NothingP() Pattern    { return Nothing() }

// Set combines patterns: values from other override this, but keep this structure (like JS set)
// Simplified: merge bags via AppBoth with keep logic
func (p Pattern) Set(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			// b overrides a
			if am, ok := a.(map[string]any); ok {
				if bm, ok2 := b.(map[string]any); ok2 {
					merged := map[string]any{}
					for k, v := range am {
						merged[k] = v
					}
					for k, v := range bm {
						merged[k] = v
					}
					return merged
				}
			}
			// if b is map, return b, else merge?
			if bm, ok := b.(map[string]any); ok {
				return bm
			}
			return b
		}
	})
	return funcPat.AppLeft(otherPat.Fmap(func(b any) any { return b }))
}

// Keep keeps values from this where not in other
func (p Pattern) Keep(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			if am, ok := a.(map[string]any); ok {
				if bm, ok2 := b.(map[string]any); ok2 {
					merged := map[string]any{}
					for k, v := range am {
						merged[k] = v
					}
					for k, v := range bm {
						if _, exists := merged[k]; !exists {
							merged[k] = v
						}
					}
					return merged
				}
			}
			return a
		}
	})
	return funcPat.AppLeft(otherPat.Fmap(func(b any) any { return b }))
}

func (p Pattern) Log2() Pattern {
	return p.Fmap(func(v any) any { return math.Log2(toFloat(v)) })
}

func (p Pattern) AsNumber() Pattern { return p.Fmap(func(v any) any { return toFloat(v) }) }

func (p Pattern) PressBy(r float64) Pattern {
	return p.Fmap(func(x any) any { return Pure(x).Compress(FractionFromFloat(r), FractionFromInt(1)) }).SqueezeJoin()
}
func (p Pattern) Press() Pattern { return p.PressBy(0.5) }
func (p Pattern) Hush() Pattern  { return Silence() }

func (p Pattern) KeepIf(other any) Pattern {
	otherPat := Reify(other)
	funcPat := p.Fmap(func(a any) any {
		return func(b any) any {
			if keep, ok := b.(bool); ok {
				if keep {
					return a
				}
				return nil
			}
			if toFloat(b) != 0 {
				return a
			}
			return nil
		}
	})
	res := funcPat.AppBoth(otherPat.Fmap(func(b any) any { return b }))
	return res.FilterValues(func(v any) bool { return v != nil })
}
