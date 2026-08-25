// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileMIDISatisfiesInterface(t *testing.T) {
	var _ MIDIInterface = NewFileMIDI(480)
}

func TestFileMIDIRecordsAndWrites(t *testing.T) {
	f := NewFileMIDI(480)
	f.At(0)
	if err := f.SendNoteOn(0, 60, 100); err != nil {
		t.Fatalf("SendNoteOn: %v", err)
	}
	f.At(480)
	if err := f.SendNoteOff(0, 60); err != nil {
		t.Fatalf("SendNoteOff: %v", err)
	}

	path := filepath.Join(t.TempDir(), "out.mid")
	if err := f.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.HasPrefix(b, []byte("MThd")) {
		t.Errorf("not a MIDI file: % X", b[:4])
	}
	if !bytes.Contains(b, []byte{0x90, 60, 100}) {
		t.Error("note-on not present in the file")
	}
}

func TestFileMIDIEmptyStillWritesValidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.mid")
	if err := NewFileMIDI(480).Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, _ := os.ReadFile(path)
	if !bytes.Contains(b, []byte{0xFF, 0x2F, 0x00}) {
		t.Error("an empty file still needs an end-of-track event")
	}
}
