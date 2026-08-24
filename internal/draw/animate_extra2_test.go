// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package draw

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestAnimateRescaleMove(t *testing.T) {
	p := core.Pure("a")
	rescaled := Rescale(0.5, p)
	if len(rescaled.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("Rescale expected haps")
	}
	moved := MoveXY(10, 20, p)
	haps := moved.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("MoveXY expected haps")
	}
	if m, ok := haps[0].Value.(map[string]any); !ok || m["x"] != 10 || m["y"] != 20 {
		t.Fatalf("MoveXY expected x=10 y=20 got %v", haps[0].Value)
	}
	zoomed := ZoomIn(2.0, p)
	haps2 := zoomed.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("ZoomIn expected haps")
	}
	if m, ok := haps2[0].Value.(map[string]any); !ok || m["zoom"] != 2.0 {
		t.Fatalf("ZoomIn expected zoom=2.0 got %v", haps2[0].Value)
	}
}

func TestFramerSequence(t *testing.T) {
	f := NewFramer()
	if f == nil {
		t.Fatalf("NewFramer nil")
	}
	// Tick should not panic
	evs := f.Tick(0.5)
	if len(evs) != 0 {
		t.Fatalf("Tick expected nil or empty")
	}
	d := NewDrawer("canvas")
	d.Draw(nil)
	GetDrawContext("id")
	CleanupDraw(false)
	CleanupDrawContext("session")
}
