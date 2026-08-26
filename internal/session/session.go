// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live session — evaluation, pattern, scheduler and output.

package session

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/jsapi"
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

	// OnError is called when the sink reports a per-event failure — a dial
	// error, a write error, an encode error. Without this, sink.Play errors
	// were discarded (`_ = sink.Play(...)`), so a bad host or an unencodable
	// value failed silently and the tool exited 0. Only the first error is
	// reported; every event that fails after it just increments errCount, so
	// a pattern that fails on every tick logs one line instead of flooding
	// the output — a caller who wants a total can read errCount.
	OnError  func(error)
	errOnce  sync.Once
	errCount atomic.Int64
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
		if err := sink.Play(h, at, cps, duration); err != nil {
			s.errCount.Add(1)
			s.errOnce.Do(func() {
				if s.OnError != nil {
					s.OnError(err)
				}
			})
		}
	}
	return s
}

// ErrCount reports how many sink errors have occurred so far, including the
// first one reported through OnError and every one suppressed after it.
func (r *Session) ErrCount() int64 { return r.errCount.Load() }

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

// Evaluate resolves code to a Pattern via jsapi.EvaluateCode (JS first,
// mini-notation fallback) and installs it as the session's live pattern.
// A JS error that mini-notation cannot rescue is returned to the caller —
// runPlay (cmd/saint-hubbins) relies on this to abort before ever touching
// a socket, rather than silently streaming a literal-string hap or nothing
// at all to SuperDirt.
func (r *Session) Evaluate(code string) (core.Pattern, error) {
	pat, err := jsapi.EvaluateCode(code)
	if err != nil {
		return core.Silence(), err
	}
	r.mu.Lock()
	r.Pattern = pat
	r.mu.Unlock()
	r.Cyclist.SetPattern(pat)
	return pat, nil
}

func (r *Session) Hush() {
	r.mu.Lock()
	r.Pattern = core.Silence()
	r.mu.Unlock()
	r.Cyclist.SetPattern(core.Silence())
}
