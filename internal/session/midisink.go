// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Playing scheduled events through a MIDI interface.

package session

import (
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	shio "codeberg.org/uzu/saint-hubbins/internal/io"
)

// MIDISink plays haps through any MIDIInterface. The note-off is scheduled on
// a timer so the scheduler callback never blocks for the length of a note.
type MIDISink struct {
	Out shio.MIDIInterface
}

func (m *MIDISink) Play(h core.Hap, at time.Time, cps, duration float64) error {
	note, vel, ch, ok := shio.HapToNote(h)
	if !ok {
		return nil
	}
	hold := time.Duration(duration * float64(time.Second))
	if hold <= 0 {
		hold = 100 * time.Millisecond
	}
	// Cyclist.OnTrigger fires every hap in a lookahead window (roughly
	// 100-200ms of events) in one tight loop, well before those events are
	// actually due. Sending the note-on synchronously here would make every
	// note in a batch sound at the same instant, collapsing the rhythm. So
	// the note-on itself is scheduled against at, the event's absolute
	// target time, exactly as OSCSink.Play threads at into SendAt. Do not
	// "simplify" this back to a direct SendNoteOn call.
	//
	// time.Until(at) can be negative or zero when the event is already due;
	// time.AfterFunc treats that as "fire immediately", which is exactly the
	// behaviour wanted here.
	time.AfterFunc(time.Until(at), func() {
		_ = m.Out.SendNoteOn(ch, note, vel)
		time.AfterFunc(hold, func() { _ = m.Out.SendNoteOff(ch, note) })
	})
	return nil
}
