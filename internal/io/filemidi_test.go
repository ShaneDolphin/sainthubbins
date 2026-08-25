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

func TestFileMIDICursorNeverSet(t *testing.T) {
	// If At() is never called, cursor stays at 0
	f := NewFileMIDI(480)
	if err := f.SendNoteOn(0, 60, 100); err != nil {
		t.Fatalf("SendNoteOn: %v", err)
	}
	if err := f.SendNoteOff(0, 60); err != nil {
		t.Fatalf("SendNoteOff: %v", err)
	}

	path := filepath.Join(t.TempDir(), "never_at.mid")
	if err := f.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, _ := os.ReadFile(path)
	// Should still produce a valid file
	if !bytes.HasPrefix(b, []byte("MThd")) {
		t.Error("not a MIDI file even when At() never called")
	}
}

func TestFileMIDICursorBackwards(t *testing.T) {
	// Cursor can move backwards; EncodeSMF handles sorting
	f := NewFileMIDI(480)
	f.At(1000)
	if err := f.SendNoteOn(0, 60, 100); err != nil {
		t.Fatalf("SendNoteOn at 1000: %v", err)
	}

	f.At(500) // Move backwards
	if err := f.SendNoteOn(1, 61, 100); err != nil {
		t.Fatalf("SendNoteOn at 500: %v", err)
	}

	f.At(2000) // Move forward
	if err := f.SendNoteOff(0, 60); err != nil {
		t.Fatalf("SendNoteOff at 2000: %v", err)
	}

	path := filepath.Join(t.TempDir(), "backwards.mid")
	if err := f.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, _ := os.ReadFile(path)
	// All events should be present despite backwards cursor movement
	if !bytes.Contains(b, []byte{0x90, 60, 100}) {
		t.Error("note-on at tick 1000 missing")
	}
	if !bytes.Contains(b, []byte{0x91, 61, 100}) {
		t.Error("note-on at tick 500 missing")
	}
	if !bytes.Contains(b, []byte{0x80, 60, 0}) {
		t.Error("note-off at tick 2000 missing")
	}
}

func TestFileMIDIClampingBehavior(t *testing.T) {
	// Test channel masking: channel 20 should mask to 4 (20 & 0x0F = 4)
	f := NewFileMIDI(480)
	f.At(0)
	if err := f.SendNoteOn(20, 60, 100); err != nil {
		t.Fatalf("SendNoteOn with channel 20: %v", err)
	}

	// Test note clamping: note 200 should mask to 72 (200 & 0x7F = 72)
	if err := f.SendNoteOff(0, 200); err != nil {
		t.Fatalf("SendNoteOff with note 200: %v", err)
	}

	// Test velocity clamping: velocity 200 should mask to 72 (200 & 0x7F = 72)
	f.At(100)
	if err := f.SendCC(0, 7, 200); err != nil {
		t.Fatalf("SendCC with val 200: %v", err)
	}

	path := filepath.Join(t.TempDir(), "clamped.mid")
	if err := f.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, _ := os.ReadFile(path)

	// Channel 20 & 0x0F = 4, so status byte is 0x90 | 4 = 0x94
	if !bytes.Contains(b, []byte{0x94, 60, 100}) {
		t.Error("note-on with channel 20 should be masked to channel 4 (0x94)")
	}

	// Note 200 & 0x7F = 72
	if !bytes.Contains(b, []byte{0x80, 72, 0}) {
		t.Error("note-off with note 200 should be masked to 72")
	}

	// CC value 200 & 0x7F = 72
	if !bytes.Contains(b, []byte{0xB0, 7, 72}) {
		t.Error("CC with val 200 should be masked to 72")
	}
}
