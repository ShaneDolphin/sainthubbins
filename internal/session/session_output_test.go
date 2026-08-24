// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"errors"
	"sync"
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// recorder is a Sink that remembers what it was asked to play.
type recorder struct {
	mu   sync.Mutex
	haps []core.Hap
}

func (r *recorder) Play(h core.Hap, at time.Time, cps, duration float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.haps = append(r.haps, h)
	return nil
}

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.haps)
}

func TestSessionSendsHapsToItsSink(t *testing.T) {
	s := NewSession()
	rec := &recorder{}
	s.SetSink(rec)

	if _, err := s.Evaluate("bd*4"); err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	// Drive the trigger path directly rather than waiting on wall-clock ticks,
	// so the test is fast and deterministic.
	h := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0]
	s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))

	if rec.count() != 1 {
		t.Fatalf("sink received %d haps, want 1", rec.count())
	}
}

func TestSessionWithoutSinkDoesNotPanic(t *testing.T) {
	s := NewSession()
	if _, err := s.Evaluate("bd"); err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	h := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0]
	s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))
}

// erroringSink always fails, so tests can drive Session.OnError without a
// real socket.
type erroringSink struct{ err error }

func (e *erroringSink) Play(h core.Hap, at time.Time, cps, duration float64) error {
	return e.err
}

func TestSessionReportsSinkErrorsViaOnError(t *testing.T) {
	s := NewSession()
	wantErr := errors.New("boom")
	s.SetSink(&erroringSink{err: wantErr})

	var mu sync.Mutex
	var calls int
	var got error
	s.OnError = func(err error) {
		mu.Lock()
		defer mu.Unlock()
		calls++
		got = err
	}

	if _, err := s.Evaluate("bd*4"); err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	h := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0]

	// Trigger the failing sink several times: a single event type failing
	// every tick must report once, not flood the output.
	for i := 0; i < 3; i++ {
		s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("OnError called %d times, want exactly 1 (flood guard)", calls)
	}
	if got != wantErr {
		t.Fatalf("OnError err = %v, want %v", got, wantErr)
	}
	if s.ErrCount() != 3 {
		t.Fatalf("ErrCount() = %d, want 3 (every failure counted even while suppressed)", s.ErrCount())
	}
}
