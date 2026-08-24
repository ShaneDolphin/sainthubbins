// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package csound

// Engine is csound @csound/browser stub
type Engine struct {
	Loaded bool
}

func New() *Engine { return &Engine{} }

func (e *Engine) CompileOrc(orc string) error { return nil }
func (e *Engine) EvalScore(score string) error { return nil }
