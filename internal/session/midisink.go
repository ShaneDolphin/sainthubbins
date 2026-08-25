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
	if err := m.Out.SendNoteOn(ch, note, vel); err != nil {
		return err
	}
	hold := time.Duration(duration * float64(time.Second))
	if hold <= 0 {
		hold = 100 * time.Millisecond
	}
	time.AfterFunc(hold, func() { _ = m.Out.SendNoteOff(ch, note) })
	return nil
}
