// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	shio "codeberg.org/uzu/saint-hubbins/internal/io"
)

func TestMIDISinkEmitsNoteOnThenOff(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 60}}

	if err := sink.Play(h, time.Now(), 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if len(mock.Snapshot()) == 0 {
		t.Fatal("no note-on was sent")
	}
	// The note-off is scheduled; give it room to land. Reading through
	// Snapshot rather than the Messages field directly is what makes this
	// safe under -race: the note-off fires on a timer goroutine, and a plain
	// field read has no happens-before relationship with that write even
	// after the mutex-guarded append inside SendNoteOff has completed.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if len(mock.Snapshot()) >= 2 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	msgs := mock.Snapshot()
	if len(msgs) < 2 {
		t.Fatalf("expected a note-off to follow, got %v", msgs)
	}
	if got, want := msgs[0], "on ch=0 note=60 vel=100"; got != want {
		t.Errorf("note-on = %q, want %q", got, want)
	}
	if got, want := msgs[1], "off ch=0 note=60"; got != want {
		t.Errorf("note-off = %q, want %q", got, want)
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
	// Wait well past where a (wrongly) scheduled note-off would land.
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

	if err := sink.Play(h, time.Now(), 0.5, 0); err != nil {
		t.Fatalf("Play: %v", err)
	}
	// Shortly after note-on, the note-off must not have landed yet.
	time.Sleep(30 * time.Millisecond)
	if msgs := mock.Snapshot(); len(msgs) != 1 {
		t.Fatalf("expected only the note-on at 30ms, got %v", msgs)
	}
	// Well past the 100ms default, the note-off must have landed.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if len(mock.Snapshot()) >= 2 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if msgs := mock.Snapshot(); len(msgs) < 2 {
		t.Fatalf("expected the default 100ms note-off to land, got %v", msgs)
	}
}
