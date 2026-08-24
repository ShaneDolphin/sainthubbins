// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package superdough

import "testing"

func TestSuperdoughNoteFreq(t *testing.T) {
	sd := New(48000)
	buf := sd.TriggerWithControls(map[string]any{"note": "c4"}, 0.1)
	if len(buf) == 0 {
		t.Fatalf("c4 empty")
	}
	buf2 := sd.TriggerWithControls(map[string]any{"note": "c5"}, 0.1)
	if len(buf2) == 0 {
		t.Fatalf("c5 empty")
	}
	// Different notes should produce different buffers (at least not identical)
	same := true
	for i := range buf {
		if buf[i] != buf2[i] {
			same = false
			break
		}
	}
	if same {
		t.Fatalf("c4 and c5 identical")
	}
}

func TestSuperdoughSampleFreq(t *testing.T) {
	sd := New(48000)
	bufBd := sd.Trigger("bd", 0.05)
	bufSd := sd.Trigger("sd", 0.05)
	if len(bufBd) == 0 || len(bufSd) == 0 {
		t.Fatalf("bd/sd empty")
	}
}
