// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package gamepad

// Gamepad is Gamepad API → pattern stub
type Gamepad struct {
	Connected bool
	Axes      []float64
	Buttons   []bool
}

func New() *Gamepad { return &Gamepad{} }
