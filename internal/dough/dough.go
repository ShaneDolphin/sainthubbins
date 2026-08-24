// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package dough

// Synth is dough-synth wrapper stub
type Synth struct {
	Name string
}

func New(name string) *Synth { return &Synth{Name: name} }

func (s *Synth) Trigger(freq, dur float64) []float32 {
	// stub sine
	return make([]float32, int(dur*48000))
}
