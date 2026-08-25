// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Standard MIDI File encoding — no dependencies, just the container format.

package io

import "encoding/binary"

// TimedEvent is a raw MIDI message at an absolute tick.
type TimedEvent struct {
	Tick uint32
	Data []byte
}

// writeVLQ encodes a variable-length quantity: seven bits per byte, with the
// high bit set on every byte except the last.
func writeVLQ(v uint32) []byte {
	if v == 0 {
		return []byte{0}
	}
	var stack []byte
	for v > 0 {
		stack = append(stack, byte(v&0x7F))
		v >>= 7
	}
	out := make([]byte, 0, len(stack))
	for i := len(stack) - 1; i >= 0; i-- {
		b := stack[i]
		if i != 0 {
			b |= 0x80
		}
		out = append(out, b)
	}
	return out
}

// EncodeSMF writes a format-0 single-track file. Events must be sorted by tick;
// deltas are derived from the gaps between them.
func EncodeSMF(ticksPerQuarter int, events []TimedEvent) []byte {
	var track []byte
	var last uint32
	for _, e := range events {
		delta := e.Tick - last
		last = e.Tick
		track = append(track, writeVLQ(delta)...)
		track = append(track, e.Data...)
	}
	// End of track.
	track = append(track, 0x00, 0xFF, 0x2F, 0x00)

	out := make([]byte, 0, len(track)+22)
	out = append(out, "MThd"...)
	out = binary.BigEndian.AppendUint32(out, 6)
	out = binary.BigEndian.AppendUint16(out, 0)                       // format 0
	out = binary.BigEndian.AppendUint16(out, 1)                       // one track
	out = binary.BigEndian.AppendUint16(out, uint16(ticksPerQuarter)) // division
	out = append(out, "MTrk"...)
	out = binary.BigEndian.AppendUint32(out, uint32(len(track)))
	out = append(out, track...)
	return out
}
