// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package gamepad

import "testing"

func TestGamepadExtra(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatalf("nil")
	}
	g.Connected = true
	g.Axes = []float64{0.5, -0.5}
	g.Buttons = []bool{true, false}
	if !g.Connected || len(g.Axes) != 2 || len(g.Buttons) != 2 {
		t.Fatalf("gamepad fields %v", g)
	}
	g2 := New()
	if g2 == nil {
		t.Fatalf("nil2")
	}
}
