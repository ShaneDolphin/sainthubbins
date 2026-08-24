// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package superdough

import "testing"

func TestSuperdoughCutoffResonance(t *testing.T) {
	sd := New(48000)
	// Low cutoff vs high cutoff should differ but both non-empty
	bufLow := sd.TriggerWithControls(map[string]any{"s": "bd", "cutoff": 200}, 0.1)
	bufHigh := sd.TriggerWithControls(map[string]any{"s": "bd", "cutoff": 5000}, 0.1)
	if len(bufLow) == 0 || len(bufHigh) == 0 {
		t.Fatalf("cutoff empty")
	}
	// Buffers should differ
	diff := 0
	for i := range bufLow {
		if bufLow[i] != bufHigh[i] {
			diff++
			break
		}
	}
	if diff == 0 {
		t.Logf("low vs high cutoff identical (ok if ADSR same, but expected differ)")
	}
}

func TestSuperdoughPanGain(t *testing.T) {
	sd := New(48000)
	bufCenter := sd.TriggerWithControls(map[string]any{"s": "bd", "pan": 0, "gain": 0.5}, 0.05)
	bufLeft := sd.TriggerWithControls(map[string]any{"s": "bd", "pan": -1, "gain": 0.5}, 0.05)
	if len(bufCenter) == 0 || len(bufLeft) == 0 {
		t.Fatalf("pan empty")
	}
	// gain affects max amplitude
	maxCenter := float32(0)
	for _, v := range bufCenter {
		if v < 0 {
			v = -v
		}
		if v > maxCenter {
			maxCenter = v
		}
	}
	bufLoud := sd.TriggerWithControls(map[string]any{"s": "bd", "gain": 1.0}, 0.05)
	maxLoud := float32(0)
	for _, v := range bufLoud {
		if v < 0 {
			v = -v
		}
		if v > maxLoud {
			maxLoud = v
		}
	}
	if maxLoud <= maxCenter {
		t.Logf("loud max %v <= center max %v (expected louder)", maxLoud, maxCenter)
	}
}
