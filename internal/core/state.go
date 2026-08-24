// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/state.mjs
package core

// State holds the time span being queried plus controls (pattern context).
// Mirrors JS State class (span + controls map).
type State struct {
	Span     TimeSpan
	Controls map[string]any
}

func NewState(span TimeSpan, controls map[string]any) State {
	if controls == nil {
		controls = map[string]any{}
	}
	// copy controls to preserve immutability
	cp := make(map[string]any, len(controls))
	for k, v := range controls {
		cp[k] = v
	}
	return State{Span: span, Controls: cp}
}

func (s State) SetSpan(span TimeSpan) State {
	return NewState(span, s.Controls)
}

func (s State) WithSpan(fn func(TimeSpan) TimeSpan) State {
	return s.SetSpan(fn(s.Span))
}

func (s State) SetControls(controls map[string]any) State {
	merged := make(map[string]any, len(s.Controls)+len(controls))
	for k, v := range s.Controls {
		merged[k] = v
	}
	for k, v := range controls {
		merged[k] = v
	}
	return NewState(s.Span, merged)
}
