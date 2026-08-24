// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package hydra

// Hydra is stub for hydra-synth bridge
type Hydra struct {
	Enabled bool
}

func New() *Hydra { return &Hydra{Enabled: false} }

func (h *Hydra) Eval(code string) error { return nil }
