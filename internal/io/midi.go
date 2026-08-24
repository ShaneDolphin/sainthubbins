// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/midi/* — MIDI via Interface.

package io

import (
	"fmt"
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// MIDIInterface abstracts MIDI output (native/webmidi/WASM/mock).
type MIDIInterface interface {
	SendNoteOn(channel, note, velocity int) error
	SendNoteOff(channel, note int) error
	SendCC(channel, cc, val int) error
	Close() error
}

// MockMIDI records messages for tests.
type MockMIDI struct {
	Messages []string
}

func (m *MockMIDI) SendNoteOn(ch, note, vel int) error {
	m.Messages = append(m.Messages, fmt.Sprintf("on ch=%d note=%d vel=%d", ch, note, vel))
	return nil
}
func (m *MockMIDI) SendNoteOff(ch, note int) error {
	m.Messages = append(m.Messages, fmt.Sprintf("off ch=%d note=%d", ch, note))
	return nil
}
func (m *MockMIDI) SendCC(ch, cc, val int) error {
	m.Messages = append(m.Messages, fmt.Sprintf("cc ch=%d cc=%d val=%d", ch, cc, val))
	return nil
}
func (m *MockMIDI) Close() error { return nil }

// MIDIFromHaps converts haps to MIDI using Mock.
func MIDIFromHaps(haps []core.Hap, midi MIDIInterface) error {
	for _, h := range haps {
		if m, ok := h.Value.(map[string]any); ok {
			noteVal, hasNote := m["note"]
			chanVal := 0
			if c, ok := m["midichan"]; ok {
				switch v := c.(type) {
				case int: chanVal = v
				case float64: chanVal = int(v)
				}
			}
			if hasNote {
				var midiNote int
				switch v := noteVal.(type) {
				case string:
					n, err := core.NoteToMidi(v)
					if err == nil { midiNote = n }
				case int:
					midiNote = v
				case float64:
					midiNote = int(v)
				}
				vel := 100
				if v, ok := m["velocity"]; ok {
					switch x := v.(type) {
					case float64: vel = int(x*127)
					case int: vel = x
					}
				}
				_ = midi.SendNoteOn(chanVal, midiNote, vel)
				_ = midi.SendNoteOff(chanVal, midiNote)
			}
		}
	}
	return nil
}

// OSCInterface for SuperDirt
type OSCInterface interface {
	Send(address string, args ...any) error
	Close() error
}

type MockOSC struct { Messages []string }
func (m *MockOSC) Send(addr string, args ...any) error {
	m.Messages = append(m.Messages, fmt.Sprintf("%s %v", addr, args))
	return nil
}
func (m *MockOSC) Close() error { return nil }

// Serial/MQTT/Gamepad stubs
type SerialInterface interface { Write(data []byte) (int, error); Close() error }
type MockSerial struct { Buf []byte }
func (m *MockSerial) Write(data []byte) (int, error) { m.Buf = append(m.Buf, data...); return len(data), nil }
func (m *MockSerial) Close() error { return nil }
