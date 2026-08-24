// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package motion

// Motion is device motion/orientation → pattern stub
type Motion struct {
	Enabled bool
	X, Y, Z float64
}

func New() *Motion { return &Motion{} }
