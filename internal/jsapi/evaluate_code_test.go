// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import (
	"strings"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestEvaluateCodeTriesJSFirst(t *testing.T) {
	pat, err := EvaluateCode(`s("bd sd")`)
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2", len(haps))
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Errorf("haps[0].Value = %#v, want a control bag carrying s:bd", haps[0].Value)
	}
}

// TestEvaluateCodeFallsBackToMini is the case the plan must not break:
// "bd sd" is two bare identifiers, not valid JS, but is exactly the
// mini-notation every existing template and eval gate depends on.
func TestEvaluateCodeFallsBackToMini(t *testing.T) {
	pat, err := EvaluateCode("bd sd")
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2 — mini-notation fallback must still work", len(haps))
	}
	if haps[0].Value != "bd" || haps[1].Value != "sd" {
		t.Errorf("values = %v, %v, want bd, sd", haps[0].Value, haps[1].Value)
	}
}

// TestEvaluateCodeReportsJSErrorWhenMiniAlsoFails is the critical half: a
// broken JS expression that mini-notation cannot rescue must surface the
// JS error, not silence and not a literal-string hap.
func TestEvaluateCodeReportsJSErrorWhenMiniAlsoFails(t *testing.T) {
	_, err := EvaluateCode(`s("bd" +`)
	if err == nil {
		t.Fatal("want an error for unparseable JS that mini-notation cannot rescue, got nil")
	}
}

func TestEvaluateCodeReportsErrorForBadMethod(t *testing.T) {
	_, err := EvaluateCode(`s("bd").nope()`)
	if err == nil {
		t.Fatal("want an error for a nonexistent method, got nil")
	}
	if !strings.Contains(err.Error(), "nope") {
		t.Errorf("error %q should mention the offending method", err.Error())
	}
}

// TestMiniNotationCorpusSurvivesTheFallback locks hapsLookGenuine's charset
// against the mini grammar it claims to mirror.
//
// hapsLookGenuine decides "did mini really parse this, or just echo the
// source back at me" from the characters in a hap's string value, and it
// tracks internal/mini's step_char set by doc comment only — nothing fails
// if the grammar grows a character the classifier does not know about. The
// symptom would be nasty and narrow: one mini form starts reporting a
// JavaScript syntax error instead of playing, while every other form and
// every gate stays green.
//
// So drive real mini-notation through the real fallback. Each of these is
// invalid JavaScript, so each one reaches mini only by failing JS first.
func TestMiniNotationCorpusSurvivesTheFallback(t *testing.T) {
	corpus := []string{
		"bd sd", "bd*2 sd", "bd:3 sd:1", "<bd sd> hh", "[bd sd] hh",
		"bd(3,8)", "bd . sd hh", "bd? sd", "bd!2 sd", "bd@3 sd",
		"~ bd", "bd,hh*4", "{bd sd, hh hh hh}", "bd [sd [hh hh]]",
		"c#4 eb3", "bd_sd", "bd.sd", "bd^2",
	}
	for _, src := range corpus {
		t.Run(src, func(t *testing.T) {
			pat, err := EvaluateCode(src)
			if err != nil {
				t.Fatalf("mini-notation %q reported a JS error: %v", src, err)
			}
			haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
			if len(haps) == 0 {
				t.Fatalf("mini-notation %q produced no haps", src)
			}
			// The echo signature this guard exists to catch: one hap whose
			// value is the source text handed straight back.
			//
			// Only meaningful for sources that cannot legitimately BE one
			// atom. "bd.sd" and "bd^2" are single atoms — . and ^ are both
			// step_chars — so value == source is the correct parse there,
			// not an echo. Checking those too is how this assertion first
			// failed against perfectly good code.
			if len(haps) == 1 && strings.ContainsAny(src, " []<>{},") {
				if s, ok := haps[0].Value.(string); ok && s == src {
					t.Errorf("%q came back as a literal-source hap, not a parse", src)
				}
			}
		})
	}
}
