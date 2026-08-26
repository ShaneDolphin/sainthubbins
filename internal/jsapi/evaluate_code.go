// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// The one fallback rule shared by every place user text becomes a Pattern.

package jsapi

import (
	"strings"
	"unicode"

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
// "Mini yields real haps" needs two qualifiers, both discovered by testing
// inputs that could break the rule rather than ones that would confirm it:
//
//  1. mini.Mini never fails outright: any token it cannot make sense of
//     falls through to its own last resort, treating the *entire* unparsed
//     string as one literal word (core.Pure(tok)) — exactly right for a bare
//     sample name like "bd" or "supersaw", and exactly wrong for a chunk of
//     broken JS that merely lacked whitespace (`s("bd").nope()`) or stopped
//     splitting on whitespace because an unbalanced "(" left the tokenizer's
//     bracket-depth counter above zero (`s("bd" +`). Both parse
//     "successfully" into one hap whose value is the original source text
//     verbatim — the silent-failure shape this whole plan exists to remove,
//     just relocated from Evaluate into the fallback meant to catch
//     Evaluate's failures. hapsLookGenuine tells the two apart by charset
//     and, for the "bd:1" colon form, by whether the sample name collides
//     with a registered control name — see its doc comment.
//
//  2. A pattern can be genuinely, correctly silent on the one cycle this
//     function inspects — "~ ~" is empty forever, and "<~ bd>" (a mainstream
//     alternation idiom: rest this bar, hit the next) is empty on cycle 0
//     specifically. hapsLookGenuine needs a hap's *value* to classify and
//     has nothing to look at when there are zero haps, so a mini result
//     with no haps falls back to looksLikeMiniSource, which judges the
//     source text instead. See its doc comment for what that trades away.
func EvaluateCode(code string) (core.Pattern, error) {
	if pat, err := Evaluate(code); err == nil {
		return pat, nil
	} else {
		mini.RegisterStringParser()
		if m := mini.Mini(code); m.Query != nil {
			haps := m.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
			if len(haps) > 0 {
				if hapsLookGenuine(haps) {
					return m, nil
				}
			} else if trimmed := strings.TrimSpace(code); looksLikeMiniSource(trimmed) && !callsRegisteredGlobal(trimmed) {
				return m, nil
			}
		}
		return core.Silence(), err
	}
}

// hapsLookGenuine reports whether every hap's value looks like something
// mini-notation's real grammar could have produced, as opposed to its
// literal-word last resort echoing unparsed input verbatim. See
// EvaluateCode's doc comment for why the distinction matters.
func hapsLookGenuine(haps []core.Hap) bool {
	for _, h := range haps {
		if !valueLooksGenuine(h.Value) {
			return false
		}
	}
	return true
}

// valueLooksGenuine checks one hap value. A plain string must be built
// entirely from mini's step_char set (isSafeMiniAtom); anything else — a
// number, nil, a bool — is trusted as-is, since a raw unvalidated string is
// the only shape mini's literal-word fallback ever produces.
//
// A map[string]any is what mini's "bd:1" colon syntax produces:
// {"s": <name>, "n": <number-or-string>} (internal/mini/mini.go's parseToken,
// the three "map[string]any{...}" returns in its colon branch). Checking its
// string members with the same atom rule is necessary but not sufficient:
// "gain:0.5" is a syntactically valid JS labeled statement, fails Evaluate
// for an unrelated reason (a bare number is not a Pattern), and mini's
// colon branch happily turns it into a plausible-looking control bag
// {"s": "gain", "n": 0.5} — "gain" and "0.5" both pass the atom check on
// their own. What gives it away is that "gain" is not a sample name from
// any drum library; it is one of this API's own registered control names
// (the `controls` table in registry.go), and no real mini-notation pattern
// names a sample after a control it could just call directly. Any "s"
// value that collides with a registered control name is rejected on that
// basis alone, regardless of what "n" turns out to be.
func valueLooksGenuine(v any) bool {
	switch val := v.(type) {
	case string:
		return isSafeMiniAtom(val)
	case map[string]any:
		if s, ok := val["s"].(string); ok {
			if _, isControlName := controls[s]; isControlName {
				return false
			}
		}
		for _, sub := range val {
			if !valueLooksGenuine(sub) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// isMiniStepChar reports whether r is one of mini-notation's step_char set —
// internal/mini/krill.peg:108's step_char rule: `unicode_letter / [0-9~] /
// "-" / "#" / "." / "^" / "_"` — the character class mini's own leaf grammar
// uses to build an atom/word. This mirrors krill.peg, the grammar's source
// of truth, not internal/mini/mini.go's GetLeafLocations, which is only an
// ASCII approximation of it used for editor highlighting. Deliberately
// Unicode-aware: "bä sd" is two Unicode-letter words, not garbage, and an
// ASCII-only version of this check rejected it.
func isMiniStepChar(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	return strings.ContainsRune("~-#.^_", r)
}

// isSafeMiniAtom reports whether s is built entirely from mini-notation's
// step_char set — see isMiniStepChar.
func isSafeMiniAtom(s string) bool {
	for _, r := range s {
		if !isMiniStepChar(r) {
			return false
		}
	}
	return true
}

// looksLikeMiniSource reports whether code, independent of what (if
// anything) it parsed into, could plausibly be pure mini-notation source:
// step_chars (isMiniStepChar) plus mini's own structural operators and
// whitespace. It exists only to accept a mini result that produced zero
// haps on the one cycle EvaluateCode inspects — hapsLookGenuine needs a
// hap's value to classify and has nothing to work with when there isn't
// one, so a genuinely (or just currently) silent pattern falls back to
// judging the source text instead.
//
// This is a best-effort heuristic, not a proof: mini-notation and JS share
// most of their ASCII punctuation, so a JS snippet that (a) fails to
// evaluate, (b) uses only these characters, and (c) parses via mini into
// zero haps on this cycle would be wrongly accepted as silence rather than
// reported as a JS error.
//
// The accepted set is mini's grammar, not a guess about which characters
// "look like" JavaScript — that guess is what previously excluded ":",
// "(", ")", "%" and "|" and broke every silent pattern using them. Only a
// quote, ";", "+", "=", "&" and "$" are genuinely JS-only, and the narrow
// residual risk is judged better than the alternative: a stricter rule
// makes "~ ~", "<~ bd>" and "<~ bd:3>" — all correct, silent-on-this-cycle
// mini-notation — report as errors.
//
// This only runs when mini produced ZERO haps. Anything with haps is judged
// by hapsLookGenuine, whose stricter charset is deliberately untouched.
func looksLikeMiniSource(code string) bool {
	for _, r := range code {
		if isMiniStepChar(r) {
			continue
		}
		switch r {
		// Mini's structural operators, taken from krill.peg rather than from
		// intuition about which punctuation "looks like" JavaScript:
		//   [] <> {} , * / ! @ ?   grouping, alternation, repeat, degrade
		//   :                      op_tail — the sample index in bd:3
		//   ( )                    op_bjorklund — euclid, bd(3,8)
		//   %                      polymeter_steps — {a b, c}%4
		//   |                      choose_tail — random choice
		// The last four were excluded on the theory that JS needs them and
		// mini does not. Mini needs all four, and every silent pattern using
		// one reported a JavaScript syntax error: <~ bd:3>, <~ bd(3,8)>,
		// ~|~, {~ ~, ~}%4.
		case '[', ']', '<', '>', '{', '}', ',', '*', '/', '!', '@', '?',
			':', '(', ')', '%', '|',
			' ', '\t', '\n', '\r', '\u00a0':
			continue
		}
		return false
	}
	return true
}

// callsRegisteredGlobal reports whether code contains an identifier
// immediately followed by "(" that names one of this API's own bound
// globals \u2014 a control (registry.go's `controls`), a variadic combinator
// (`variadic`), or one of the two standalone functions `mini`/`silence`.
//
// looksLikeMiniSource widening to accept "(" and ")" (op_bjorklund, bd(3,8))
// opened exactly this hole: mini's euclid branch splits the text between
// the parens on "," WITHOUT tracking bracket depth (internal/mini/mini.go's
// parseToken, the `strings.Split(inside, ",")` in its "(" handling \u2014 every
// other comma-splitter in that file uses the depth-aware splitAtDepth0).
// So `stack(s(bd*4), s(hh*8).gain(0.4))` \u2014 the console's own default
// snippet with the quotes forgotten \u2014 reads as euclid("stack", pulses=0,
// steps=0) because Atoi fails on both non-numeric halves, which is a valid
// Pattern that happens to produce zero haps forever. hapsLookGenuine never
// runs (there's nothing to check), and looksLikeMiniSource accepts the
// source because every character in it \u2014 letters, "( ) * ," and space \u2014 is
// legitimately part of mini's grammar. `stack(bd, sd)`, `cat(...)`,
// `note(...)` and `sequence(...)` all fail the identical way.
//
// This is the same "our own name is not a sample name" test
// valueLooksGenuine already applies to the colon form ("gain:0.5"), applied
// here to the call form instead: no real mini-notation pattern is a euclid
// rhythm gated on a sample named "stack" or "note" \u2014 those are this API's
// own names for something else entirely. `bd(3,8)` and `foo(1,2)` are
// unaffected, since neither "bd" nor "foo" is a bound global.
func callsRegisteredGlobal(code string) bool {
	runes := []rune(code)
	for i := 0; i < len(runes); {
		if !isJSIdentStart(runes[i]) {
			i++
			continue
		}
		start := i
		for i < len(runes) && isJSIdentPart(runes[i]) {
			i++
		}
		if i < len(runes) && runes[i] == '(' && isBoundGlobalName(string(runes[start:i])) {
			return true
		}
	}
	return false
}

// isBoundGlobalName reports whether name is bound directly into the JS
// global scope by register() (registry.go): a control, a variadic
// combinator, or one of the two standalone functions mini/silence.
func isBoundGlobalName(name string) bool {
	if _, ok := controls[name]; ok {
		return true
	}
	if _, ok := variadic[name]; ok {
		return true
	}
	return name == "mini" || name == "silence"
}

// isJSIdentStart and isJSIdentPart approximate JS identifier syntax closely
// enough for callsRegisteredGlobal's purpose \u2014 finding a run of characters
// that could be one of a small, known list of ASCII names \u2014 without
// needing the real grammar. Unicode letters are included because
// unicode.IsLetter already is elsewhere in this file (isMiniStepChar); no
// bound global name uses one, so it never changes the answer, only the
// scan.
func isJSIdentStart(r rune) bool {
	return unicode.IsLetter(r) || r == '_' || r == '$'
}

func isJSIdentPart(r rune) bool {
	return isJSIdentStart(r) || unicode.IsDigit(r)
}
