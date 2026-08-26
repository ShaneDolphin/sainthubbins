// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// The one fallback rule shared by every place user text becomes a Pattern.

package jsapi

import (
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// EvaluateCode resolves user-typed text to a Pattern for every call site
// that accepts it — the CLI's eval/render/midi/play commands and the web
// console's /api/evaluate and /api/pianoroll handlers. It is the single
// place the JS-first, mini-notation-fallback rule lives; callers must not
// hand-copy it — this package has already had the identical "logic outside
// the shared guard" bug three times.
//
// The rule: try JS first via Evaluate. If that fails AND the mini-notation
// parser yields real haps, use the mini-notation pattern instead — this is
// what keeps bare mini-notation like `bd sd` working, since two bare
// identifiers are not valid JS. Otherwise — JS failed and mini produced
// nothing usable — return the JS error rather than falling back to a
// literal-string hap or a silently empty pattern. That "otherwise" is the
// entire point of this plan: a user who mistypes `s("bd" +` must see a
// syntax error, not silence.
//
// "Mini yields real haps" needs one more qualifier than it sounds like it
// does. mini.Mini never fails outright: any token it cannot make sense of
// falls through to its own last resort, treating the *entire* unparsed
// string as one literal word (core.Pure(tok)) — exactly right for a bare
// sample name like "bd" or "supersaw", and exactly wrong for a chunk of
// broken JS that merely lacked whitespace (`s("bd").nope()`) or stopped
// splitting on whitespace because an unbalanced "(" left the tokenizer's
// bracket-depth counter above zero (`s("bd" +`). Both of those parse "successfully"
// into one hap whose value is the original source text verbatim — the
// silent-failure shape this whole plan exists to remove, just relocated
// from Evaluate into the fallback meant to catch Evaluate's failures.
// hapsLookGenuine tells the two apart by charset: every value mini's real
// grammar ever assigns to an atom is built from letters, digits, and the
// punctuation its own leaf grammar accepts (see GetLeafLocations in
// internal/mini/mini.go) — quotes, parens, and a bare "+" never appear in a
// value unless mini gave up and echoed raw input back.
func EvaluateCode(code string) (core.Pattern, error) {
	if pat, err := Evaluate(code); err == nil {
		return pat, nil
	} else {
		mini.RegisterStringParser()
		if m := mini.Mini(code); m.Query != nil {
			haps := m.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
			if len(haps) > 0 && hapsLookGenuine(haps) {
				return m, nil
			}
			// "~" and "" are mini-notation's own spellings of "nothing" —
			// core.Evaluate treated them as valid, silent input rather than
			// an error, and this preserves that.
			if len(haps) == 0 {
				if trimmed := strings.TrimSpace(code); trimmed == "" || trimmed == "~" {
					return m, nil
				}
			}
		}
		return core.Silence(), err
	}
}

// hapsLookGenuine reports whether every string-valued hap looks like a
// value mini-notation's real grammar could have produced, as opposed to
// its literal-word last resort echoing unparsed input verbatim. See
// EvaluateCode's doc comment for why the distinction matters.
func hapsLookGenuine(haps []core.Hap) bool {
	for _, h := range haps {
		if s, ok := h.Value.(string); ok && !isSafeMiniAtom(s) {
			return false
		}
	}
	return true
}

// isSafeMiniAtom reports whether s only contains characters mini-notation's
// leaf grammar accepts inside a word: letters, digits, and #-^_.~ (the same
// set internal/mini/mini.go's GetLeafLocations documents as step_char).
func isSafeMiniAtom(s string) bool {
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case strings.ContainsRune("#-^_.~", r):
		default:
			return false
		}
	}
	return true
}
