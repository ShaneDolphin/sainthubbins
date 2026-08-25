// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	shio "codeberg.org/uzu/saint-hubbins/internal/io"
)

// waitForCount bounded-polls mock until it has at least n messages, or fails
// the test after a deadline. Polling (rather than a flat sleep) keeps this
// fast on a quiet machine and robust on a loaded one.
func waitForCount(t *testing.T, mock *shio.MockMIDI, n int) []string {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if msgs := mock.Snapshot(); len(msgs) >= n {
			return msgs
		}
		time.Sleep(5 * time.Millisecond)
	}
	return mock.Snapshot()
}

func TestMIDISinkEmitsNoteOnThenOff(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 60}}

	// at is in the near future, as it would be for a hap inside the
	// scheduler's lookahead window: the note-on must not land until then.
	at := time.Now().Add(60 * time.Millisecond)
	if err := sink.Play(h, at, 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}

	// Shortly after Play returns but well before at, nothing should have
	// been sent yet — this is what proves the note-on is scheduled against
	// at rather than fired synchronously.
	time.Sleep(15 * time.Millisecond)
	if msgs := mock.Snapshot(); len(msgs) != 0 {
		t.Fatalf("note-on fired before its scheduled time: %v", msgs)
	}

	msgs := waitForCount(t, mock, 2)
	if len(msgs) < 2 {
		t.Fatalf("expected a note-on then a note-off, got %v", msgs)
	}
	if got, want := msgs[0], "on ch=0 note=60 vel=100"; got != want {
		t.Errorf("note-on = %q, want %q", got, want)
	}
	if got, want := msgs[1], "off ch=0 note=60"; got != want {
		t.Errorf("note-off = %q, want %q", got, want)
	}
}

// TestMIDISinkSchedulesNoteOnAgainstAt confirms the note-on itself lands at
// (approximately) at, not at call time — the core of the timing fix: without
// it, every hap in a Cyclist lookahead batch would sound at the same instant.
func TestMIDISinkSchedulesNoteOnAgainstAt(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 67}}

	delay := 80 * time.Millisecond
	start := time.Now()
	at := start.Add(delay)
	if err := sink.Play(h, at, 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}

	msgs := waitForCount(t, mock, 1)
	elapsed := time.Since(start)
	if len(msgs) < 1 {
		t.Fatalf("expected a note-on, got %v", msgs)
	}
	if elapsed < delay {
		t.Errorf("note-on landed after %v, before its scheduled delay of %v", elapsed, delay)
	}
}

// TestMIDISinkAlreadyDueFiresImmediately confirms an at value in the past
// (or now) is treated as "fire immediately" rather than blocking Play or
// being dropped — time.AfterFunc's documented behaviour for a non-positive
// duration.
func TestMIDISinkAlreadyDueFiresImmediately(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 60}}

	at := time.Now().Add(-50 * time.Millisecond)
	start := time.Now()
	if err := sink.Play(h, at, 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("Play blocked for %v on an already-due event; it must return promptly", elapsed)
	}

	msgs := waitForCount(t, mock, 1)
	if len(msgs) < 1 || msgs[0] != "on ch=0 note=60 vel=100" {
		t.Fatalf("expected an immediate note-on, got %v", msgs)
	}
}

// TestMIDISinkSkipsPitchlessHap covers the requirement that a hap with no
// resolvable pitch is skipped entirely: no note-on is sent, and (just as
// importantly) no note-off timer is scheduled that could later fire a stray
// message against a mock the caller believes is quiescent.
func TestMIDISinkSkipsPitchlessHap(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"gain": 0.8}}

	if err := sink.Play(h, time.Now(), 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}
	// Wait well past where a (wrongly) scheduled note-on/note-off would land.
	time.Sleep(150 * time.Millisecond)
	if msgs := mock.Snapshot(); len(msgs) != 0 {
		t.Errorf("expected no messages for a pitchless hap, got %v", msgs)
	}
}

// TestMIDISinkZeroDurationDefaultsHold confirms that a non-positive duration
// falls back to the brief's 100ms default hold rather than firing the
// note-off immediately (or never).
func TestMIDISinkZeroDurationDefaultsHold(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 64}}

	at := time.Now()
	if err := sink.Play(h, at, 0.5, 0); err != nil {
		t.Fatalf("Play: %v", err)
	}
	// Wait for the note-on to land, then confirm the note-off hasn't yet.
	waitForCount(t, mock, 1)
	time.Sleep(30 * time.Millisecond)
	if msgs := mock.Snapshot(); len(msgs) != 1 {
		t.Fatalf("expected only the note-on at 30ms, got %v", msgs)
	}
	// Well past the 100ms default, the note-off must have landed.
	msgs := waitForCount(t, mock, 2)
	if len(msgs) < 2 {
		t.Fatalf("expected the default 100ms note-off to land, got %v", msgs)
	}
}
