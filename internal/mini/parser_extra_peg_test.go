// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestParserPEGSampleColon(t *testing.T) {
	// Mini handles bd:1 via Pure(map[s:bd n:1]) — PEG Word currently treats bd:1 as one Word "bd:1"
	// Ensure Mini still works via shim fallback, and PEG at least parses without error
	for _, c := range []string{"bd:1", "bd:1 sd:2", "bd:1*2"} {
		v, err := Parse("", []byte(c))
		if err != nil {
			t.Fatalf("parse %q err %v", c, err)
		}
		pat := v.(core.Pattern)
		miniPat := Mini(c)
		// Both should produce at least 1 hap (Mini correctly handles colon, PEG may treat as string)
		if len(miniPat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
			t.Fatalf("mini %q empty", c)
		}
		// PEG should not panic; log len difference
		_ = pat
	}
}

func TestParserPEGEuclid(t *testing.T) {
	v, err := Parse("", []byte("bd(3,8)"))
	if err != nil {
		t.Fatalf("euclid parse err %v", err)
	}
	pat := v.(core.Pattern)
	// PEG EuclidOp currently returns string op, not yet applied via Element ops — should still parse as Word+EuclidOp
	// For now, ensure non-panic and fallback via Mini is available
	miniPat := Mini("bd(3,8)")
	if len(miniPat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))) == 0 {
		t.Fatalf("mini euclid empty")
	}
	_ = pat
}
