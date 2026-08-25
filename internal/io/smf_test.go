// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"bytes"
	"testing"
)

func TestWriteVLQ(t *testing.T) {
	cases := map[uint32][]byte{
		0:       {0x00},
		127:     {0x7F},
		128:     {0x81, 0x00},
		8192:    {0xC0, 0x00},
		1048576: {0xC0, 0x80, 0x00},
	}
	for in, want := range cases {
		if got := writeVLQ(in); !bytes.Equal(got, want) {
			t.Errorf("writeVLQ(%d) = % X, want % X", in, got, want)
		}
	}
}

func TestEncodeSMFStructure(t *testing.T) {
	out := EncodeSMF(480, []TimedEvent{
		{Tick: 0, Data: []byte{0x90, 60, 100}}, // note on
		{Tick: 480, Data: []byte{0x80, 60, 0}}, // note off a quarter later
	})
	if !bytes.HasPrefix(out, []byte("MThd")) {
		t.Fatalf("missing MThd header: % X", out[:4])
	}
	// Verify exact header field values.
	// Format at offset 8-9 must be 0x0000 (format 0)
	if out[8] != 0x00 || out[9] != 0x00 {
		t.Errorf("format field should be 0, got [% X]", out[8:10])
	}
	// Track count at offset 10-11 must be 0x0001 (one track)
	if out[10] != 0x00 || out[11] != 0x01 {
		t.Errorf("ntrks field should be 1, got [% X]", out[10:12])
	}
	// Division at offset 12-13 must be 0x01E0 (480 in big-endian)
	if out[12] != 0x01 || out[13] != 0xE0 {
		t.Errorf("division field should be 480, got [% X]", out[12:14])
	}
	if !bytes.Contains(out, []byte("MTrk")) {
		t.Fatal("missing MTrk chunk")
	}
	// Verify MTrk chunk length matches actual track data.
	mtkIdx := bytes.Index(out, []byte("MTrk"))
	if mtkIdx < 0 {
		t.Fatal("MTrk not found")
	}
	// Length field is 4 bytes after "MTrk"
	expectedLen := uint32(len(out) - mtkIdx - 8) // 8 = 4 ("MTrk") + 4 (length field)
	actualLen := uint32(out[mtkIdx+4])<<24 | uint32(out[mtkIdx+5])<<16 |
		uint32(out[mtkIdx+6])<<8 | uint32(out[mtkIdx+7])
	if actualLen != expectedLen {
		t.Errorf("MTrk length field is %d, expected %d", actualLen, expectedLen)
	}
	// End-of-track meta event is mandatory.
	if !bytes.Contains(out, []byte{0xFF, 0x2F, 0x00}) {
		t.Error("missing end-of-track meta event")
	}
}

func TestEncodeSMFUsesDeltaTimes(t *testing.T) {
	// Two events at the same tick must produce a zero delta for the second.
	out := EncodeSMF(480, []TimedEvent{
		{Tick: 100, Data: []byte{0x90, 60, 100}},
		{Tick: 100, Data: []byte{0x90, 64, 100}},
	})
	// The second event's delta is 0x00 immediately before its status byte.
	if !bytes.Contains(out, []byte{0x00, 0x90, 64, 100}) {
		t.Errorf("second simultaneous event should have a zero delta: % X", out)
	}
}

func TestEncodeSMFSortsOutOfOrderEvents(t *testing.T) {
	// Out-of-order input should produce the same output as sorted input.
	unsorted := []TimedEvent{
		{Tick: 100, Data: []byte{0x90, 60, 100}},
		{Tick: 50, Data: []byte{0x90, 64, 100}},
		{Tick: 200, Data: []byte{0x80, 60, 0}},
	}
	sorted := []TimedEvent{
		{Tick: 50, Data: []byte{0x90, 64, 100}},
		{Tick: 100, Data: []byte{0x90, 60, 100}},
		{Tick: 200, Data: []byte{0x80, 60, 0}},
	}

	outUnsorted := EncodeSMF(480, unsorted)
	outSorted := EncodeSMF(480, sorted)

	if !bytes.Equal(outUnsorted, outSorted) {
		t.Errorf("out-of-order input should produce same output as sorted\nUnsorted: % X\nSorted:   % X", outUnsorted, outSorted)
	}
}

// TestEncodeSMFPreservesSimultaneousEventOrder guards sort stability, which
// is load-bearing: RenderMIDI on repeated notes ("c4 c4") emits
// on@0, off@960, on@960, off@1920, and it is stability that keeps the
// note-off before the note-on at the shared tick 960. An unstable sort would
// let SendNoteOn(960) and SendNoteOff(960) swap, retriggering-then-killing
// the repeated note in a DAW — with the rest of the suite still green.
//
// A single tied pair cannot catch that reliably: below Go's insertion-sort
// threshold sort.Slice happens to match sort.SliceStable, and — verified
// empirically against this exact package — even above that threshold a lone
// tied pair sitting among otherwise-distinct, already-sorted ticks still
// comes out in input order under the plain (unstable) sort; pdqsort's
// pattern-detecting fast paths paper over it. What actually forces a
// divergence is SEVERAL tied groups: four ticks, four events tied at each,
// interleaved round-robin. Confirmed by repeatedly running this exact
// fixture through sort.Slice vs sort.SliceStable: they disagree on every
// run. (See the fix report for the swap-and-rerun proof against the real
// EncodeSMF.)
func TestEncodeSMFPreservesSimultaneousEventOrder(t *testing.T) {
	ticks := []uint32{100, 200, 300, 400}
	var events []TimedEvent
	var groups [4][]TimedEvent // groups[i] = the four tied events at ticks[i], in input order
	id := 0
	for round := 0; round < 4; round++ {
		for i, tick := range ticks {
			// note = 60+id uniquely identifies this event so its position in
			// the output is unambiguous.
			e := TimedEvent{Tick: tick, Data: []byte{0x90, byte(60 + id), 100}}
			events = append(events, e)
			groups[i] = append(groups[i], e)
			id++
		}
	}
	if len(events) < 13 {
		t.Fatalf("fixture has only %d events, want at least 13 (above the insertion-sort threshold)", len(events))
	}

	out := EncodeSMF(480, events)

	// Every tied group must appear, back-to-back (zero delta between
	// consecutive members), in its original input order.
	for i, group := range groups {
		var expected []byte
		for j, e := range group {
			if j > 0 {
				expected = append(expected, 0x00) // zero delta: same tick as predecessor
			}
			expected = append(expected, e.Data...)
		}
		if !bytes.Contains(out, expected) {
			t.Errorf("tied group at tick %d should preserve input order: expected % X in % X", ticks[i], expected, out)
		}
	}
}

func TestEncodeSMFDoesNotMutateCaller(t *testing.T) {
	// The caller's slice should not be modified.
	original := []TimedEvent{
		{Tick: 200, Data: []byte{0x80, 60, 0}},
		{Tick: 100, Data: []byte{0x90, 60, 100}},
	}
	original_copy := make([]TimedEvent, len(original))
	copy(original_copy, original)

	_ = EncodeSMF(480, original)

	if !bytes.Equal(
		bytes.Join([][]byte{original[0].Data, original[1].Data}, nil),
		bytes.Join([][]byte{original_copy[0].Data, original_copy[1].Data}, nil),
	) || original[0].Tick != original_copy[0].Tick || original[1].Tick != original_copy[1].Tick {
		t.Errorf("EncodeSMF modified the caller's slice: original %+v, after %+v", original_copy, original)
	}
}

func TestEncodeSMFEmptyEventList(t *testing.T) {
	// Empty event list should produce a valid file with only the end-of-track event.
	out := EncodeSMF(480, []TimedEvent{})
	if !bytes.HasPrefix(out, []byte("MThd")) {
		t.Fatal("missing MThd header")
	}
	if !bytes.Contains(out, []byte("MTrk")) {
		t.Fatal("missing MTrk chunk")
	}
	if !bytes.Contains(out, []byte{0xFF, 0x2F, 0x00}) {
		t.Error("missing end-of-track meta event")
	}
	// Track should contain only the end-of-track event (4 bytes) plus its delta (1 byte).
	mtkIdx := bytes.Index(out, []byte("MTrk"))
	if mtkIdx < 0 {
		t.Fatal("MTrk not found")
	}
	trackLen := uint32(out[mtkIdx+4])<<24 | uint32(out[mtkIdx+5])<<16 |
		uint32(out[mtkIdx+6])<<8 | uint32(out[mtkIdx+7])
	// Minimal track: 0x00 (zero delta) + 0xFF 0x2F 0x00 (end-of-track) = 4 bytes
	if trackLen != 4 {
		t.Errorf("empty event list should have track length 4, got %d", trackLen)
	}
}
