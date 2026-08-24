// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package draw

import (
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"math"
)

// PitchWheelEvent is drawable for pitchwheel (polar pitch vs time)
type PitchWheelEvent struct {
	Angle    float64       `json:"angle"`
	Radius   float64       `json:"radius"`
	Value    any           `json:"value"`
	Context  map[string]any `json:"context,omitempty"`
}

// PitchWheel maps haps to pitchwheel polar coordinates
// JS: pitchwheel maps note midi to angle, time to radius similar to spiral but with pitch
func PitchWheel(haps []core.Hap) []PitchWheelEvent {
	evs := make([]PitchWheelEvent, 0, len(haps))
	for _, h := range haps {
		if h.Whole == nil {
			continue
		}
		// Angle from pitch class (C=0 etc.)
		pitch := 0.0
		switch v := h.Value.(type) {
		case string:
			pitch = float64(noteToMidi(v) % 12)
		case map[string]any:
			if n, ok := v["note"]; ok {
				if s, ok := n.(string); ok {
					pitch = float64(noteToMidi(s) % 12)
				} else if f, ok := n.(float64); ok {
					pitch = math.Mod(f, 12)
				}
			} else if v2, ok := v["n"]; ok {
				switch x := v2.(type) {
				case string:
					pitch = float64(noteToMidi(x) % 12)
				case float64:
					pitch = math.Mod(x, 12)
				case int:
					pitch = math.Mod(float64(x), 12)
				}
			}
		case float64:
			pitch = math.Mod(v, 12)
		case int:
			pitch = math.Mod(float64(v), 12)
		}
		angle := pitch / 12 * 2 * math.Pi
		radius := h.Whole.Begin.Sam().Float64() + 1 + h.Whole.Duration().Float64()
		evs = append(evs, PitchWheelEvent{
			Angle:   angle,
			Radius:  radius,
			Value:   h.Value,
			Context: h.Context,
		})
	}
	return evs
}

// ColorMap mirrors JS color.mjs colorMap (basic)
var ColorMap = map[string]string{
	"bd": "red", "sd": "blue", "hh": "yellow", "cp": "green",
	"c": "red", "d": "orange", "e": "yellow", "f": "green", "g": "blue", "a": "purple", "b": "pink",
}

// ConvertColorToNumber mirrors JS convertColorToNumber
func ConvertColorToNumber(color string) int {
	if v, ok := ColorMap[color]; ok {
		color = v
	}
	// Simple hex parse #rrggbb
	if len(color) == 7 && color[0] == '#' {
		var n int
		for _, c := range color[1:] {
			n <<= 4
			if c >= '0' && c <= '9' {
				n |= int(c - '0')
			} else if c >= 'a' && c <= 'f' {
				n |= int(c - 'a' + 10)
			} else if c >= 'A' && c <= 'F' {
				n |= int(c - 'A' + 10)
			}
		}
		return n
	}
	m := map[string]int{"red": 0xFF0000, "green": 0x00FF00, "blue": 0x0000FF, "yellow": 0xFFFF00, "orange": 0xFF8800, "purple": 0x8800FF, "pink": 0xFF88FF, "white": 0xFFFFFF, "black": 0x000000}
	if v, ok := m[color]; ok {
		return v
	}
	return 0x888888
}

func noteToMidi(note string) int {
	if len(note) == 0 {
		return 60
	}
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := note
	// lower
	b := []byte(n)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	n = string(b)
	letter := n[0]
	semi, ok := base[letter]
	if !ok {
		return 60
	}
	idx := 1
	if idx < len(n) && (n[idx] == '#' || n[idx] == 'b') {
		if n[idx] == '#' {
			semi++
		} else {
			semi--
		}
		idx++
	}
	oct := 4
	if idx < len(n) {
		o := 0
		mult := 1
		// simple parseInt suffix
		start := idx
		neg := false
		if start < len(n) && n[start] == '-' {
			neg = true
			start++
		}
		for _, c := range n[start:] {
			if c >= '0' && c <= '9' {
				o = o*10 + int(c-'0')
			}
		}
		if neg {
			o = -o
		}
		if o != 0 || (len(n) > idx && n[idx] >= '0' && n[idx] <= '9') {
			oct = o
		}
		// handle mult for negative?
		_ = mult
		_ = neg
	}
	return (oct+1)*12 + semi
}
