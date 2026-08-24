// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live session — evaluation, pattern, scheduler and output.

package session

import (
	"context"
	"sync"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
	"codeberg.org/uzu/saint-hubbins/internal/osc"
)

// Sink receives scheduled events. Keeping this an interface means the session
// knows nothing about OSC, and the trigger path can be tested without a socket.
type Sink interface {
	Play(h core.Hap, at time.Time, cps, duration float64) error
}

// OSCSink plays events through SuperDirt.
type OSCSink struct{ Client *osc.Client }

func (s *OSCSink) Play(h core.Hap, at time.Time, cps, duration float64) error {
	return s.Client.SendAt(at, osc.DirtAddress, osc.DirtArgs(h, cps, duration)...)
}

// Session ties evaluation, pattern, scheduler and output together.
type Session struct {
	mu      sync.RWMutex
	Pattern core.Pattern
	Cyclist *core.Cyclist
	sink    Sink
}

// NewSession creates a new live session (these go to eleven).
func NewSession() *Session {
	mini.RegisterStringParser()
	s := &Session{
		Cyclist: core.NewCyclist(0.1, nil, nil),
		Pattern: core.Silence(),
	}
	// Cyclist computes targetTime as absolute seconds since the Unix epoch.
	s.Cyclist.OnTrigger = func(h core.Hap, deadline, duration, cps, targetTime float64) {
		s.mu.RLock()
		sink := s.sink
		s.mu.RUnlock()
		if sink == nil {
			return
		}
		at := time.Unix(0, int64(targetTime*1e9))
		_ = sink.Play(h, at, cps, duration)
	}
	return s
}

// SetSink installs the output. Passing nil silences the session.
func (r *Session) SetSink(s Sink) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sink = s
}

// Start runs the scheduler until ctx is cancelled.
func (r *Session) Start(ctx context.Context) error { return r.Cyclist.Start(ctx) }

// Stop halts the scheduler.
func (r *Session) Stop() { r.Cyclist.Stop() }

func (r *Session) Evaluate(code string) (core.Pattern, error) {
	pat, _, err := core.Evaluate(code, nil)
	if err != nil {
		pat = mini.Mini(code)
		if pat.Query == nil {
			pat = core.Pure(code)
		}
		err = nil
	}
	r.mu.Lock()
	r.Pattern = pat
	r.mu.Unlock()
	r.Cyclist.SetPattern(pat)
	return pat, err
}

func (r *Session) Hush() {
	r.mu.Lock()
	r.Pattern = core.Silence()
	r.mu.Unlock()
	r.Cyclist.SetPattern(core.Silence())
}
