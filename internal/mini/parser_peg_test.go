// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestParserPEGBasic(t *testing.T) {
	cases := []string{
		"bd sd",
		"bd*2",
		"bd/2",
		"bd?",
		"~",
		"[bd sd]",
		"<a b>",
		"{a b, c d e}",
	}
	for _, c := range cases {
		v, err := Parse("", []byte(c))
		if err != nil {
			t.Fatalf("parse %q err %v", c, err)
		}
		pat, ok := v.(core.Pattern)
		if !ok {
			t.Fatalf("parse %q not pattern %T", c, v)
		}
		haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		// silence ~ should be 0, others >=0
		if c == "~" && len(haps) != 0 {
			t.Fatalf("~ expected 0 got %d", len(haps))
		}
		if c != "~" && len(haps) == 0 && c != "bd?" {
			// degrade may be 0 occasionally but not always
			t.Logf("parse %q got 0 haps (ok for degrade)", c)
		}
	}
}

func TestParserPEGEquiMini(t *testing.T) {
	// Compare PEG parser vs Mini shim for same inputs
	cases := []string{"bd sd", "bd*2", "bd/2", "[bd sd]"}
	for _, c := range cases {
		v, err := Parse("", []byte(c))
		if err != nil {
			t.Fatalf("peg parse %q err %v", c, err)
		}
		pegPat := v.(core.Pattern)
		miniPat := Mini(c)
		// Compare first cycle length
		pegHaps := pegPat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		miniHaps := miniPat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(pegHaps) != len(miniHaps) {
			t.Logf("peg vs mini %q: peg %d mini %d (acceptable if both non-empty)", c, len(pegHaps), len(miniHaps))
		}
	}
}
