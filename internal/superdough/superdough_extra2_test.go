// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package superdough

import "testing"

func TestSuperdoughCutoffPan(t *testing.T) {
	sd := New(48000)
	bufLow := sd.TriggerWithControls(map[string]any{"s": "bd", "cutoff": 100, "gain": 1.0}, 0.1)
	sd2 := New(48000)
	bufHigh := sd2.TriggerWithControls(map[string]any{"s": "bd", "cutoff": 20000, "gain": 1.0}, 0.1)
	if len(bufLow) != len(bufHigh) {
		t.Fatalf("len mismatch %d vs %d", len(bufLow), len(bufHigh))
	}
	bufPanL := sd.TriggerWithControls(map[string]any{"s": "bd", "pan": -1, "gain": 1.0}, 0.05)
	bufPanR := sd.TriggerWithControls(map[string]any{"s": "bd", "pan": 1, "gain": 1.0}, 0.05)
	if len(bufPanL) == 0 || len(bufPanR) == 0 {
		t.Fatalf("pan empty")
	}
	_ = bufLow
	_ = bufHigh
}

func TestSuperdoughADSR(t *testing.T) {
	sd := New(48000)
	buf := sd.TriggerWithControls(map[string]any{"s": "bd", "attack": 0.01, "decay": 0.05, "sustain": 0.5, "release": 0.1, "gain": 1.0}, 0.2)
	if len(buf) == 0 {
		t.Fatalf("ADSR empty")
	}
	max := float32(0)
	for _, v := range buf {
		if v < 0 {
			v = -v
		}
		if v > max {
			max = v
		}
	}
	if max == 0 {
		t.Fatalf("max 0")
	}
}
