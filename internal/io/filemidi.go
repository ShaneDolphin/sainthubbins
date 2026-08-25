// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// A MIDIInterface that writes a Standard MIDI File instead of talking to a device.

package io

import "os"

// FileMIDI records MIDI messages against a tick cursor and writes them out as
// a Standard MIDI File. MIDIInterface carries no timing, so callers move the
// cursor with At before sending each event.
type FileMIDI struct {
	TicksPerQuarter int
	cursor          uint32
	events          []TimedEvent
}

// Compile-time assertion that FileMIDI satisfies MIDIInterface.
var _ MIDIInterface = (*FileMIDI)(nil)

func NewFileMIDI(ticksPerQuarter int) *FileMIDI {
	if ticksPerQuarter <= 0 {
		ticksPerQuarter = 480
	}
	return &FileMIDI{TicksPerQuarter: ticksPerQuarter}
}

// At moves the cursor. Subsequent messages are stamped at this tick.
func (f *FileMIDI) At(tick uint32) { f.cursor = tick }

// NoteOnCount reports how many note-on messages have been recorded so far.
// Callers use this to tell a pattern that resolved to real notes apart from
// one that resolved to nothing — a mini-notation pattern of bare numerics,
// for example, produces a header-only file with zero note-ons and no error.
func (f *FileMIDI) NoteOnCount() int {
	n := 0
	for _, e := range f.events {
		if len(e.Data) > 0 && e.Data[0]&0xF0 == 0x90 {
			n++
		}
	}
	return n
}

func (f *FileMIDI) record(data ...byte) error {
	f.events = append(f.events, TimedEvent{Tick: f.cursor, Data: data})
	return nil
}

func (f *FileMIDI) SendNoteOn(channel, note, velocity int) error {
	return f.record(byte(0x90|channel&0x0F), byte(note&0x7F), byte(velocity&0x7F))
}

func (f *FileMIDI) SendNoteOff(channel, note int) error {
	return f.record(byte(0x80|channel&0x0F), byte(note&0x7F), 0)
}

func (f *FileMIDI) SendCC(channel, cc, val int) error {
	return f.record(byte(0xB0|channel&0x0F), byte(cc&0x7F), byte(val&0x7F))
}

func (f *FileMIDI) Close() error { return nil }

// Write encodes everything recorded so far to path.
func (f *FileMIDI) Write(path string) error {
	return os.WriteFile(path, EncodeSMF(f.TicksPerQuarter, f.events), 0o644)
}
