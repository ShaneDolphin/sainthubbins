// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

import "math"

// Perlin returns simple pseudo-perlin signal (sin-based)
func Perlin() Pattern {
	return Signal(func(t Fraction) float64 {
		x := t.Float64() * 2 * math.Pi
		return (math.Sin(x) + math.Sin(x*2.3)*0.5 + math.Sin(x*4.1)*0.25) / 1.75
	})
}

func PerlinWith(seed float64) Pattern {
	return Signal(func(t Fraction) float64 {
		x := t.Float64()*2*math.Pi + seed
		return math.Sin(x)*0.5 + math.Sin(x*2.3)*0.3
	})
}
