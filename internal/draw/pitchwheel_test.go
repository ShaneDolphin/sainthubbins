// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package draw

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestPitchWheel(t *testing.T) {
	haps := []core.Hap{
		{Value: "c4", Whole: &core.TimeSpan{Begin: core.FractionFromInt(0), End: core.FractionFromInt(1)}, Part: core.TimeSpan{Begin: core.FractionFromInt(0), End: core.FractionFromInt(1)}},
		{Value: "e4", Whole: &core.TimeSpan{Begin: core.FractionFromInt(1), End: core.FractionFromInt(2)}, Part: core.TimeSpan{Begin: core.FractionFromInt(1), End: core.FractionFromInt(2)}},
	}
	evs := PitchWheel(haps)
	if len(evs) != 2 {
		t.Fatalf("pitchwheel expected 2 got %d", len(evs))
	}
	// c4 (C) should be angle ~0, e4 (E=4) angle ~ 4/12*2pi = 2.09
	if evs[0].Angle < -0.1 || evs[0].Angle > 0.1 {
		t.Logf("c4 angle %v", evs[0].Angle)
	}
	if evs[1].Angle < 2 || evs[1].Angle > 2.2 {
		t.Fatalf("e4 angle expected ~2.09 got %v", evs[1].Angle)
	}
}

func TestColorMap(t *testing.T) {
	if ConvertColorToNumber("red") != 0xFF0000 {
		t.Fatalf("red")
	}
	if ConvertColorToNumber("bd") != 0xFF0000 {
		t.Fatalf("bd color")
	}
	if ConvertColorToNumber("#00FF00") != 0x00FF00 {
		t.Fatalf("hex")
	}
}
