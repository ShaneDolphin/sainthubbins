// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
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
