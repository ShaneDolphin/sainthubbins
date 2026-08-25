// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Turning pattern events into MIDI notes.

package io

import (
	"strconv"
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// drumNotes maps the drum names the engine knows to General MIDI percussion.
var drumNotes = map[string]int{
	"bd": 36, "sd": 38, "hh": 42, "oh": 46, "ch": 42, "cp": 39,
}

// HapToNote resolves a hap to a MIDI note. Precedence matches the audio
// renderer: n, then note, then a drum name. An event carrying no pitch returns
// ok == false so callers skip it rather than emitting a spurious middle C.
func HapToNote(h core.Hap) (note, velocity, channel int, ok bool) {
	velocity, channel = 100, 0
	m, isBag := h.Value.(map[string]any)
	if !isBag {
		if s, isStr := h.Value.(string); isStr {
			if n, found := drumNotes[s]; found {
				return n, velocity, 9, true
			}
			if n, found := noteNameToMIDI(s); found {
				return n, velocity, channel, true
			}
		}
		return 0, 0, 0, false
	}

	if g, found := m["gain"]; found {
		velocity = int(toF(g) * 127)
		if velocity < 1 {
			velocity = 1
		}
		if velocity > 127 {
			velocity = 127
		}
	}
	if c, found := m["channel"]; found {
		channel = int(toF(c))
	}
	if v, found := m["n"]; found {
		return int(toF(v)), velocity, channel, true
	}
	if v, found := m["note"]; found {
		if s, isStr := v.(string); isStr {
			if n, parsed := noteNameToMIDI(s); parsed {
				return n, velocity, channel, true
			}
			return 0, 0, 0, false
		}
		return int(toF(v)), velocity, channel, true
	}
	if v, found := m["s"]; found {
		if s, isStr := v.(string); isStr {
			if n, found := drumNotes[s]; found {
				return n, velocity, 9, true
			}
		}
	}
	return 0, 0, 0, false
}

func toF(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	}
	return 0
}

// noteNameToMIDI parses "c4", "f#3", "eb2". Middle C is c4 = 60.
func noteNameToMIDI(s string) (int, bool) {
	if len(s) < 2 {
		return 0, false
	}
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := strings.ToLower(s)
	semi, found := base[n[0]]
	if !found {
		return 0, false
	}
	i := 1
	if len(n) > 2 && (n[1] == '#' || n[1] == 'b') {
		if n[1] == '#' {
			semi++
		} else {
			semi--
		}
		i = 2
	}
	oct, err := strconv.Atoi(n[i:])
	if err != nil {
		return 0, false
	}
	return (oct+1)*12 + semi, true
}

// RenderMIDI queries cycles of pat and records paired note-on/note-off events.
// One cycle is one bar of four quarters, matching the offline audio renderer.
func RenderMIDI(pat core.Pattern, cycles, ticksPerQuarter int) *FileMIDI {
	f := NewFileMIDI(ticksPerQuarter)
	ticksPerCycle := float64(ticksPerQuarter * 4)
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(int64(cycles)))
	for _, h := range haps {
		if h.Whole == nil || !h.HasOnset() {
			continue
		}
		note, vel, ch, ok := HapToNote(h)
		if !ok {
			continue
		}
		f.At(uint32(h.Whole.Begin.Float64() * ticksPerCycle))
		_ = f.SendNoteOn(ch, note, vel)
		f.At(uint32(h.Whole.End.Float64() * ticksPerCycle))
		_ = f.SendNoteOff(ch, note)
	}
	return f
}
