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

// TestZeroHapMiniNotationIsNotAnError covers mini-notation that is
// genuinely (or just currently) silent on the one cycle EvaluateCode
// inspects, so hapsLookGenuine has no hap value to classify at all.
// "~ ~" is empty forever; "<~ bd>" and "bd?1" are empty on cycle 0
// specifically (a mainstream alternation idiom, and a 100%-probability
// degrade). All three are valid JS syntax errors and must fall back to
// mini rather than surfacing the JS error.
func TestZeroHapMiniNotationIsNotAnError(t *testing.T) {
	for _, src := range []string{"~ ~", "[~ ~]", "<~ bd>", "bd?1"} {
		t.Run(src, func(t *testing.T) {
			if _, err := EvaluateCode(src); err != nil {
				t.Errorf("EvaluateCode(%q) = %v, want a silent pattern, not a JS error", src, err)
			}
		})
	}
}

// TestUnicodeMiniNotationAtoms covers isMiniStepChar being Unicode-aware,
// not ASCII-only: "bä sd" is two Unicode-letter words, not garbage.
func TestUnicodeMiniNotationAtoms(t *testing.T) {
	pat, err := EvaluateCode("bä sd")
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2", len(haps))
	}
	if haps[0].Value != "bä" || haps[1].Value != "sd" {
		t.Errorf("values = %v, %v, want bä, sd", haps[0].Value, haps[1].Value)
	}
}

// TestColonSyntaxCollidingWithAControlNameIsAnError is the review's
// concrete counterexample: "gain:0.5" is valid JS (a labeled statement),
// fails Evaluate for an unrelated reason (a bare number is not a Pattern),
// and mini's colon branch would otherwise happily turn it into a
// plausible-looking control bag {"s": "gain", "n": 0.5}. Both "gain" and
// "0.5" pass the atom charset on their own — this must be caught by the
// control-name collision in valueLooksGenuine specifically.
func TestColonSyntaxCollidingWithAControlNameIsAnError(t *testing.T) {
	if _, err := EvaluateCode("gain:0.5"); err == nil {
		t.Fatal(`EvaluateCode("gain:0.5") = nil error, want a JS error — "gain" is a control name, not a sample`)
	}
}

// TestColonSyntaxWithARealSampleNameStillWorks guards against
// TestColonSyntaxCollidingWithAControlNameIsAnError's fix overreaching:
// legitimate "bd:1" sample:n syntax must still work.
func TestColonSyntaxWithARealSampleNameStillWorks(t *testing.T) {
	pat, err := EvaluateCode("bd:1")
	if err != nil {
		t.Fatalf("EvaluateCode: %v", err)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("got %d haps, want 1", len(haps))
	}
	m, ok := haps[0].Value.(map[string]any)
	if !ok || m["s"] != "bd" {
		t.Errorf("value = %#v, want a control bag carrying s:bd", haps[0].Value)
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
