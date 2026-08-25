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

func TestEncodeSMFPreservesSimultaneousEventOrder(t *testing.T) {
	// Events at the same tick must preserve their input order (stability).
	// We verify this by checking that a note-off before a note-on at the same
	// tick produces them in that order.
	events := []TimedEvent{
		{Tick: 100, Data: []byte{0x80, 60, 0}},   // note off
		{Tick: 100, Data: []byte{0x90, 64, 100}}, // note on
	}
	out := EncodeSMF(480, events)

	// First event has delta from 0 to 100 (encoded as 0x64), then note off data.
	// Second event has delta from 100 to 100 (encoded as 0x00), then note on data.
	// We look for the sequence after the first delta: note-off data, then zero delta, then note-on data.
	expected := []byte{0x80, 60, 0, 0x00, 0x90, 64, 100}
	if !bytes.Contains(out, expected) {
		t.Errorf("simultaneous events should preserve input order: expected % X in % X", expected, out)
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
