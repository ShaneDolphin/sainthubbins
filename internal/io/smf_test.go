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
	if !bytes.Contains(out, []byte("MTrk")) {
		t.Fatal("missing MTrk chunk")
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
