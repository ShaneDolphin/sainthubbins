// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live session — evaluation, pattern and scheduler.

package session

import (
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// Session ties evaluation, pattern and scheduler together — the live console state.
type Session struct {
	Pattern core.Pattern
	Cyclist *core.Cyclist
}

// NewSession creates a new live session (these go to eleven).
func NewSession() *Session {
	mini.RegisterStringParser()
	return &Session{
		Cyclist: core.NewCyclist(0.1, nil, nil),
		Pattern: core.Silence(),
	}
}

func (r *Session) Evaluate(code string) (core.Pattern, error) {
	pat, _, err := core.Evaluate(code, nil)
	if err != nil {
		pat = mini.Mini(code)
		if pat.Query == nil {
			pat = core.Pure(code)
		}
		err = nil
	}
	r.Pattern = pat
	r.Cyclist.SetPattern(pat)
	return pat, err
}

func (r *Session) Hush() {
	r.Pattern = core.Silence()
	r.Cyclist.SetPattern(r.Pattern)
}
