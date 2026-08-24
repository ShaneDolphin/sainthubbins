// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/mini/mini.mjs (261 LOC) + krill.pegjs (303 LOC)
// This is a minimal pure-Go mini-notation implementation (subset) for Phase 2.

package mini

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// Mini parses a mini-notation string into a Pattern.
// Supports: "bd sd", "bd ~ sd", "bd*2", "bd/2", "bd(3,8)", "<a b>", "[a b]" etc. via simple tokenizer.
// For full PEG parity, use pigeon-generated parser (krill.peg) — this is Phase 2 shim.

func Mini(input string) core.Pattern {
	input = strings.TrimSpace(input)
	// Remove surrounding quotes if present
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}
	if input == "" {
		return core.Silence()
	}
	return parseSequence(input)
}

func parseSequence(input string) core.Pattern {
	// Stacking binds loosest: "bd*4, hh*8" is two layers, not a sequence.
	// Split on depth-0 commas before anything else so operators inside a
	// nested group are never treated as separators.
	if parts := splitAtDepth0(input, ','); len(parts) > 1 {
		pats := make([]core.Pattern, 0, len(parts))
		for _, part := range parts {
			if part = strings.TrimSpace(part); part != "" {
				pats = append(pats, Mini(part))
			}
		}
		if len(pats) == 1 {
			return pats[0]
		}
		if len(pats) > 1 {
			return core.Stack(pats...)
		}
	}
	// Random choice binds tighter than stacking: "a | b" picks one per cycle.
	if parts := splitAtDepth0(input, '|'); len(parts) > 1 {
		choices := make([]any, 0, len(parts))
		for _, part := range parts {
			if part = strings.TrimSpace(part); part != "" {
				choices = append(choices, Mini(part))
			}
		}
		if len(choices) == 1 {
			if p, ok := choices[0].(core.Pattern); ok {
				return p
			}
		}
		if len(choices) > 1 {
			return core.Pure(0).Choose(choices)
		}
	}
	// Very simplified: split by spaces respecting brackets/angles
	// For now, handle: "bd sd", "bd ~", "bd*2", "bd:1", etc.
	// Tokenize on spaces not inside []<>()
	tokens := splitMiniTokens(input)
	if len(tokens) == 0 {
		return core.Silence()
	}
	// Depth-0 "," and "|" are handled at the top of this function, before
	// tokenizing, so any isolated separator reaching this point is a degenerate
	// input such as "," on its own; buildSteps skips it and falls through to
	// Silence below.
	//
	// A single real token used to shortcut straight to parseToken(tokens[0]),
	// bypassing weight stripping — "bd@3" would then reach parseToken still
	// carrying its "@3" suffix, and since parseToken no longer understands
	// "@" the raw suffix leaked into the value. Every token, including a lone
	// one, must go through buildSteps (and so splitStepBase) first.
	pats, weights, hadRange := buildSteps(tokens)
	if len(pats) == 0 {
		return core.Silence()
	}
	if len(pats) == 1 {
		// Exactly one step: there are no siblings for its weight to be
		// relative to, so the weight is discarded — the value still went
		// through splitWeight to strip a trailing "@n", but the timing
		// already spans the full cycle without TimeCatWeighted's
		// involvement.
		return pats[0]
	}
	if hadRange {
		// This branch doesn't support weighted timing against a range
		// (there's no sibling duration for its weight to be relative to
		// across a value/value transition), so weights are discarded same
		// as they always were for this path — only a leaked suffix in the
		// value itself would be new breakage, and buildSteps already
		// stripped that.
		return core.FastCat(pats...)
	}
	weightedAny := false
	for _, w := range weights {
		if w != 1 {
			weightedAny = true
			break
		}
	}
	if !weightedAny {
		return core.FastCat(pats...)
	}
	weighted := make([]any, 0, len(pats)*2)
	for i, p := range pats {
		weighted = append(weighted, weights[i], p)
	}
	return core.TimeCatWeighted(weighted...)
}

// buildSteps resolves a token list (as produced by splitMiniTokens) into the
// sequence of steps it names. It is the one place that combines the three
// things every list-building call site needs — parseSequence's main loop,
// the "<...>" alternation branch and the "{...}" polymeter branch — so a fix
// to any of them only has to happen once:
//
//   - ".." range expansion: "a .. b" between two integer tokens expands into
//     one Pure(v) step per value in the range, inclusive, in the direction
//     from a to b. Only the first such range in the list is expanded — this
//     matches the pre-existing single-range behavior of this grammar.
//   - "!" replicate: splitStepBase's reps count turns one token into that
//     many sibling steps.
//   - "@" weight: splitStepBase's weight is returned alongside each step so
//     a caller that divides the cycle by relative duration (parseSequence)
//     can use it; a caller for which a step is always exactly one slot
//     ("<...>", "{...}") can simply ignore the weights slice.
//
// A range step always carries weight 1, and hadRange reports whether a
// range was found — parseSequence uses that to fall back to plain
// concatenation for the whole list, since this grammar has never supported
// weighted timing against a range.
func buildSteps(tokens []string) (pats []core.Pattern, weights []float64, hadRange bool) {
	appendToken := func(tok string) {
		base, reps, w := splitStepBase(tok)
		pat := parseToken(base)
		for k := 0; k < reps; k++ {
			pats = append(pats, pat)
			weights = append(weights, w)
		}
	}

	rangeAt := -1
	for i := 0; i < len(tokens); i++ {
		if tokens[i] != ".." || i == 0 || i+1 >= len(tokens) {
			continue
		}
		if _, err1 := strconv.Atoi(tokens[i-1]); err1 != nil {
			continue
		}
		if _, err2 := strconv.Atoi(tokens[i+1]); err2 != nil {
			continue
		}
		rangeAt = i
		break
	}

	if rangeAt < 0 {
		for _, tok := range tokens {
			if tok == ".." || tok == "|" || tok == "," {
				continue
			}
			appendToken(tok)
		}
		return pats, weights, false
	}

	for j := 0; j < rangeAt-1; j++ {
		if t := tokens[j]; t == ".." || t == "|" || t == "," {
			continue
		}
		appendToken(tokens[j])
	}
	l, _ := strconv.Atoi(tokens[rangeAt-1])
	r, _ := strconv.Atoi(tokens[rangeAt+1])
	if l <= r {
		for v := l; v <= r; v++ {
			pats = append(pats, core.Pure(v))
			weights = append(weights, 1)
		}
	} else {
		for v := l; v >= r; v-- {
			pats = append(pats, core.Pure(v))
			weights = append(weights, 1)
		}
	}
	for j := rangeAt + 2; j < len(tokens); j++ {
		if t := tokens[j]; t == ".." || t == "|" || t == "," {
			continue
		}
		appendToken(tokens[j])
	}
	return pats, weights, true
}

// splitWeight separates a trailing @n weight from a token. "bd@3" yields
// ("bd", 3). A token with no weight yields a weight of 1, so callers can treat
// every step uniformly.
func splitWeight(tok string) (string, float64) {
	i := indexAtDepth0(tok, "@")
	if i <= 0 {
		return tok, 1
	}
	w, err := strconv.ParseFloat(strings.TrimSpace(tok[i+1:]), 64)
	if err != nil || w <= 0 {
		return tok[:i], 1
	}
	return tok[:i], w
}

// splitReplicate separates a trailing !n from a token. "bd!3" yields ("bd", 3);
// a bare "bd!" yields ("bd", 2). A token with no ! yields a count of 1.
func splitReplicate(tok string) (string, int) {
	i := indexAtDepth0(tok, "!")
	if i <= 0 {
		return tok, 1
	}
	rest := strings.TrimSpace(tok[i+1:])
	if rest == "" {
		return tok[:i], 2
	}
	n, err := strconv.Atoi(rest)
	if err != nil || n < 1 {
		return tok[:i], 1
	}
	return tok[:i], n
}

// splitStepBase resolves the two suffixes a mini-notation step token can
// carry — replicate (!n) and weight (@n) — the same way parseSequence
// resolves them for its own tokens. Any call site that treats a token as one
// slot in a sequence-like list (a step in "a b c", an alternative in
// "<a b>") must run it through here before parseToken, or the raw suffix
// falls straight through parseToken's fallback and leaks into the value —
// parseToken itself no longer understands either suffix, since both need
// context (sibling durations for @, sibling slots for !) that only the
// caller building the list has.
//
// splitReplicate runs first so the replicate count is read from the token as
// written; splitWeight then strips @n from what's left. Callers with nothing
// for a weight to be relative to (an alternation slot in "<...>" is always
// exactly one full cycle) can simply ignore the returned weight.
func splitStepBase(tok string) (base string, reps int, weight float64) {
	base, reps = splitReplicate(tok)
	base, weight = splitWeight(base)
	return base, reps, weight
}

func splitMiniTokens(s string) []string {
	var tokens []string
	var cur strings.Builder
	depth := 0
	for _, r := range s {
		switch r {
		case '[', '<', '(', '{':
			depth++
			cur.WriteRune(r)
		case ']', '>', ')', '}':
			depth--
			cur.WriteRune(r)
		case ' ', '\t', '\n':
			if depth == 0 {
				if cur.Len() > 0 {
					tokens = append(tokens, cur.String())
					cur.Reset()
				}
			} else {
				cur.WriteRune(r)
			}
		// Handle '|' and ',' as separate tokens when depth 0? Keep ',' inside curly as separator but depth already handles
		// For range ".." we keep as part of tokenization via spaces; but "|" should be kept inside token if no spaces like "a|b"
		// We will not split on '|' here — handle at sequence/token level
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

// --- Bracket-depth-aware scanning helpers -----------------------------------
// Mini-notation operators (* / ( @ ! ? | ,) must never be matched inside a
// nested group. Scanning with strings.Index/Contains splits tokens mid-bracket
// and yields garbage atoms like "[bd". These helpers only match at depth 0.

// bracketDelta reports the nesting change contributed by r.
func bracketDelta(r byte) int {
	switch r {
	case '[', '<', '(', '{':
		return 1
	case ']', '>', ')', '}':
		return -1
	}
	return 0
}

// indexAtDepth0 returns the index of the first byte of s that appears in chars
// at bracket depth 0, or -1 when there is none.
func indexAtDepth0(s, chars string) int {
	depth := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if depth == 0 && strings.IndexByte(chars, c) >= 0 {
			return i
		}
		depth += bracketDelta(c)
		if depth < 0 {
			depth = 0
		}
	}
	return -1
}

// containsAtDepth0 reports whether any byte in chars occurs at bracket depth 0.
func containsAtDepth0(s, chars string) bool { return indexAtDepth0(s, chars) >= 0 }

// splitAtDepth0 splits s on sep, ignoring separators inside brackets.
func splitAtDepth0(s string, sep byte) []string {
	var parts []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if depth == 0 && c == sep {
			parts = append(parts, s[start:i])
			start = i + 1
			continue
		}
		depth += bracketDelta(c)
		if depth < 0 {
			depth = 0
		}
	}
	return append(parts, s[start:])
}

// closerFor returns the bracket that closes open, or 0 if open is not an
// opening bracket.
func closerFor(open byte) byte {
	switch open {
	case '[':
		return ']'
	case '<':
		return '>'
	case '(':
		return ')'
	case '{':
		return '}'
	}
	return 0
}

// unwrapGroup returns the inner text of a token that is wholly enclosed by a
// single matched bracket pair — "[a b]" -> "a b" — so the group is parsed as a
// unit before any operator scanning happens.
//
// The closing bracket must match the opening one. Without that check "[bd)"
// and "<bd]" would parse as though correctly closed, turning a typo into a
// plausible but different pattern instead of leaving it as a visible literal.
func unwrapGroup(tok string) (inner string, open byte, ok bool) {
	if len(tok) < 2 {
		return "", 0, false
	}
	open = tok[0]
	if bracketDelta(open) != 1 {
		return "", 0, false
	}
	depth := 0
	for i := 0; i < len(tok); i++ {
		depth += bracketDelta(tok[i])
		if depth == 0 {
			// Closed early: the group does not span the whole token
			// (e.g. "[a b]*2"), so leave it to the operator handlers.
			if i != len(tok)-1 {
				return "", 0, false
			}
			if tok[i] != closerFor(open) {
				return "", 0, false
			}
			return tok[1 : len(tok)-1], open, true
		}
	}
	return "", 0, false
}

func parseToken(tok string) core.Pattern {
	// Handle rest
	if tok == "~" || tok == "-" || tok == "_" {
		return core.Silence()
	}
	// A token that is entirely one bracket group is parsed as a unit first.
	// Operators inside it belong to the group, not to this token.
	if inner, open, ok := unwrapGroup(tok); ok {
		switch open {
		case '[':
			return parseSequence(inner)
		case '<':
			// Depth-0 commas stack independent alternations, matching the way
			// they stack sequences inside "[...]": "<a b, c d>" layers
			// SlowCat(a, b) with SlowCat(c, d). Splitting here also stops a
			// trailing comma riding along on a token and producing an
			// empty-valued hap.
			var pats []core.Pattern
			for _, part := range splitAtDepth0(inner, ',') {
				toks := splitMiniTokens(strings.TrimSpace(part))
				if len(toks) == 0 {
					continue
				}
				// Each token here is one alternative occupying exactly one
				// cycle, the same as a step in a sequence — so it needs the
				// same buildSteps treatment a sequence step list gets, not a
				// bare parseToken call. "!" expands into repeated
				// alternatives ("<bd!3 sd>" is bd, bd, bd, sd across four
				// cycles); ".." expands a range the same as it does at
				// sequence level ("<0 .. 3>" is 0, 1, 2, 3 across four
				// cycles); "@" has nothing to be relative to inside a
				// one-cycle slot, so its weight is simply discarded — but
				// both suffixes still have to be stripped from the value,
				// or parseToken's fallback leaks the raw "!3"/"@3" text into
				// the hap.
				sub, _, _ := buildSteps(toks)
				pats = append(pats, core.SlowCat(sub...))
			}
			switch len(pats) {
			case 0:
				return core.Silence()
			case 1:
				return pats[0]
			}
			return core.Stack(pats...)
		}
		// '{' (polymeter) and '(' fall through to their existing handlers.
	}
	// Handle range operator inside token like "0..4" without spaces
	if strings.Contains(tok, "..") && !strings.HasPrefix(tok, "..") && !strings.HasSuffix(tok, "..") {
		parts := strings.Split(tok, "..")
		if len(parts) == 2 {
			leftStr := strings.TrimSpace(parts[0])
			rightStr := strings.TrimSpace(parts[1])
			if leftStr != "" && rightStr != "" {
				if l, err1 := strconv.Atoi(leftStr); err1 == nil {
					if r, err2 := strconv.Atoi(rightStr); err2 == nil {
						if l <= r {
							pats := make([]core.Pattern, 0, r-l+1)
							for i := l; i <= r; i++ {
								pats = append(pats, core.Pure(i))
							}
							return core.FastCat(pats...)
						}
						pats := make([]core.Pattern, 0, l-r+1)
						for i := l; i >= r; i-- {
							pats = append(pats, core.Pure(i))
						}
						return core.FastCat(pats...)
					}
				}
			}
		}
	}
	// Handle choice operator "|" — random choice (Choose)
	if containsAtDepth0(tok, "|") {
		// Avoid infinite recursion if tok is exactly "|" (should be handled as tokenization, but just in case)
		if tok != "|" {
			parts := splitAtDepth0(tok, '|')
			choices := make([]any, 0, len(parts))
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}
				choices = append(choices, parseToken(part))
			}
			if len(choices) > 0 {
				return core.Pure(0).Choose(choices)
			}
		}
	}
	// Handle curly polymeter "{a b, c d e}" and "{a b, c d e}%3" / "{a b, c d e}*3"
	if len(tok) >= 2 && tok[0] == '{' {
		closeIdx := strings.LastIndex(tok, "}")
		if closeIdx > 0 {
			inner := tok[1:closeIdx]
			suffix := strings.TrimSpace(tok[closeIdx+1:])

			// Each comma-separated layer is a list of steps. Each token here
			// is one slot in that list, the same as a step in a sequence, so
			// it needs buildSteps's treatment: "!" expands into repeated
			// sibling steps, ".." expands a range the same as it does at
			// sequence level, and "@" is stripped from the value (there is
			// no sibling duration for its weight to be relative to inside a
			// rate-based list, so it is discarded like it is for "<...>").
			// Calling parseToken directly, as on a bare token, would leak a
			// raw "!2"/"@2"/".." suffix straight into the value instead.
			var layers [][]core.Pattern
			for _, part := range splitAtDepth0(inner, ',') {
				toks := splitMiniTokens(strings.TrimSpace(part))
				if len(toks) == 0 {
					continue
				}
				steps, _, _ := buildSteps(toks)
				layers = append(layers, steps)
			}
			if len(layers) == 0 {
				return core.Silence()
			}

			// Steps per cycle: %n if given, otherwise the first layer's length.
			// This is what makes it a polymeter — every layer runs at the same
			// rate, and a layer whose length does not divide that rate lands on
			// different elements each cycle.
			stepsPerCycle := len(layers[0])
			if strings.HasPrefix(suffix, "%") {
				// %n's digits end where the next operator (*n or /n) begins —
				// "%4*2" is steps-per-cycle 4 with a *2 still to apply, not a
				// malformed "%4*2" that silently discards the *2.
				rest := suffix[1:]
				digits := 0
				for digits < len(rest) && rest[digits] >= '0' && rest[digits] <= '9' {
					digits++
				}
				if digits > 0 {
					if v, err := strconv.Atoi(rest[:digits]); err == nil && v > 0 {
						stepsPerCycle = v
					}
				}
				suffix = strings.TrimSpace(rest[digits:])
			}

			pats := make([]core.Pattern, 0, len(layers))
			for _, steps := range layers {
				// SlowCat gives one element per cycle; speeding it up by the
				// step count gives that many elements per cycle, wrapping
				// around the layer's own length.
				pats = append(pats, core.SlowCat(steps...).
					FastF(core.NewFraction(int64(stepsPerCycle), 1)))
			}
			base := pats[0]
			if len(pats) > 1 {
				base = core.Stack(pats...)
			}

			// A leftover *n or /n still applies on top.
			if strings.HasPrefix(suffix, "*") {
				if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
					return base.FastF(core.FractionFromFloat(v))
				}
			} else if strings.HasPrefix(suffix, "/") {
				if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
					return base.SlowF(core.FractionFromFloat(v))
				}
			}
			return base
		}
	}
	// Handle suffix modifiers _N and ?. @ weight and ! replicate are both
	// handled by parseSequence instead: @ needs sibling steps' relative
	// durations, and ! needs to add sibling steps of its own, so neither can
	// be resolved from a single token in isolation.
	if containsAtDepth0(tok, "_?") {
		// Handle degrade ? and ?0.5
		if containsAtDepth0(tok, "?") {
			parts := splitAtDepth0(tok, '?')
			base := parts[0]
			pat := parseToken(base)
			// degradeBy 0.5 default
			prob := 0.5
			if len(parts) > 1 && parts[1] != "" {
				if v, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
					prob = v
				}
			}
			return pat.DegradeBy(prob)
		}
		// Handle weight alias _: "bd_2" means weight 2 (elongate).
		if containsAtDepth0(tok, "_") {
			parts := splitAtDepth0(tok, '_')
			base := parts[0]
			weightStr := parts[1]
			pat := parseToken(base)
			if w, err := strconv.ParseFloat(strings.TrimSpace(weightStr), 64); err == nil && w != 0 {
				// Return pat with steps weight via WithSteps
				return pat.WithSteps(func(f core.Fraction) core.Fraction { return f.Mul(core.FractionFromFloat(w)) })
			}
			return pat
		}
	}
	// Handle euclid: "bd(3,8)" or "bd(3,8,2)"
	if indexAtDepth0(tok, "(") > 0 && strings.HasSuffix(tok, ")") {
		idx := indexAtDepth0(tok, "(")
		base := tok[:idx]
		inside := tok[idx+1 : len(tok)-1]
		parts := strings.Split(inside, ",")
		if len(parts) >= 2 {
			pulses, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			steps, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			rotation := 0
			if len(parts) >= 3 {
				rotation, _ = strconv.Atoi(strings.TrimSpace(parts[2]))
			}
			basePat := parseToken(base)
			if rotation != 0 {
				return basePat.EuclidRot(pulses, steps, rotation)
			}
			return basePat.Euclid(pulses, steps)
		}
	}
	// Handle fast/slow: "bd*2" "bd/2" "bd*2/3" etc.
	if containsAtDepth0(tok, "*/") {
		// Split base and modifiers
		baseEnd := indexAtDepth0(tok, "*/")
		if baseEnd > 0 {
			base := tok[:baseEnd]
			mods := tok[baseEnd:]
			pat := parseToken(base)
			// Parse modifiers sequentially
			i := 0
			for i < len(mods) {
				op := mods[i]
				i++
				numStr := ""
				for i < len(mods) && mods[i] >= '0' && mods[i] <= '9' || (i < len(mods) && mods[i] == '.') {
					numStr += string(mods[i])
					i++
				}
				if numStr == "" {
					continue
				}
				val, _ := strconv.ParseFloat(numStr, 64)
				if op == '*' {
					pat = pat.FastF(core.FractionFromFloat(val))
				} else if op == '/' {
					pat = pat.SlowF(core.FractionFromFloat(val))
				}
			}
			return pat
		}
	}
	// Handle angle brackets: "<a b>" = slowcat
	if len(tok) >= 2 && tok[0] == '<' && tok[len(tok)-1] == '>' {
		inner := tok[1 : len(tok)-1]
		tokens := splitMiniTokens(inner)
		pats := make([]core.Pattern, len(tokens))
		for i, t := range tokens {
			pats[i] = parseToken(t)
		}
		return core.SlowCat(pats...)
	}
	// Handle square brackets: "[a b]" = fastcat
	if len(tok) >= 2 && tok[0] == '[' && tok[len(tok)-1] == ']' {
		inner := tok[1 : len(tok)-1]
		tokens := splitMiniTokens(inner)
		pats := make([]core.Pattern, len(tokens))
		for i, t := range tokens {
			pats[i] = parseToken(t)
		}
		return core.FastCat(pats...)
	}
	// Handle stack with comma: "bd,sd" ? Not standard mini, but handle
	if containsAtDepth0(tok, ",") && !strings.Contains(tok, ":") {
		parts := splitAtDepth0(tok, ',')
		pats := make([]core.Pattern, len(parts))
		for i, p := range parts {
			pats[i] = parseToken(strings.TrimSpace(p))
		}
		return core.Stack(pats...)
	}
	// Handle note with colon sample syntax: "bd:1" => s + n
	if strings.Contains(tok, ":") && !strings.Contains(tok, "..") {
		parts := strings.Split(tok, ":")
		if len(parts) == 2 {
			base := strings.TrimSpace(parts[0])
			nStr := strings.TrimSpace(parts[1])
			if base != "" && nStr != "" {
				// Try int n
				if n, err := strconv.Atoi(nStr); err == nil {
					// Direct map with s and n to preserve sample
					return core.Pure(map[string]any{"s": base, "n": n})
				}
				if f, err := strconv.ParseFloat(nStr, 64); err == nil {
					return core.Pure(map[string]any{"s": base, "n": f})
				}
				// Fallback: treat as string n (e.g., bd:foo)
				return core.Pure(map[string]any{"s": base, "n": nStr})
			}
		}
	}
	// For generic, try numeric or note
	if _, err := strconv.ParseFloat(tok, 64); err == nil {
		// numeric pure
		return core.Pure(tok)
	}
	// Default: pure string value (e.g., "bd", "c3")
	// If parsable as note, keep as string; controls will interpret
	return core.Pure(tok)
}

// PatternifyAST stub for compatibility with krill.peg generated parser
func PatternifyAST(input string) core.Pattern {
	return Mini(input)
}

// GetLeafLocations returns [from,to] offsets of leaf atoms inside quoted mini string.
// code should include quotes, e.g. "\"bd sd\"". start is global offset of opening quote (default 0).
// Mirrors JS getLeafLocations(code, start) which parses via krill and returns leaf positions.
func GetLeafLocations(code string, offsets ...int) [][2]int {
	start := 0
	if len(offsets) > 0 {
		start = offsets[0]
	}
	// Extract content between outer quotes (support " ' `)
	if len(code) < 2 {
		return nil
	}
	quote := code[0]
	if quote != '"' && quote != '\'' && quote != '`' {
		// try to find first quote
		return nil
	}
	if code[len(code)-1] != quote {
		// not properly quoted, fallback to scanning whole
		quote = 0
	}
	var content string
	var contentOffset int
	if quote != 0 {
		content = code[1 : len(code)-1]
		contentOffset = start + 1
	} else {
		content = code
		contentOffset = start
	}
	// Scan content for leaf atoms: sequences of letters/digits/#/_-./~?* but treat ~ as rest? For locations, we want words like bd, sd, c3, 2 etc.
	// Simplistic: find runs of [A-Za-z0-9#_\-\.]+ and numbers, ignoring punctuation * / [ ] < > ( ) , | . etc.
	// Use manual scan to capture positions.
	var locs [][2]int
	i := 0
	for i < len(content) {
		c := content[i]
		// Separators that delimit leaves: whitespace and operators * / @ ! ? : ( ) [ ] < > { } , |  and also whitespace.
		// Note: '.' '_' '#' '-' '^' '~' are part of leaf per step_char, so NOT separators here (they are consumed as leaf chars).
		// For "bd*2", '*' separates "bd" and "2"; for "bd@2", '@' separates.
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == ',' || c == '|' || c == '[' || c == ']' || c == '<' || c == '>' || c == '(' || c == ')' || c == '{' || c == '}' || c == ':' || c == '@' || c == '!' || c == '?' || c == '*' || c == '/' {
			i++
			continue
		}
		// Start of leaf: consume while allowed leaf chars: letters, digits, #, -, ^, _, ., ~ (step_char)
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '#' || c == '-' || c == '^' || c == '_' || c == '.' || c == '~' {
			startLocal := i
			for i < len(content) && ((content[i] >= 'a' && content[i] <= 'z') || (content[i] >= 'A' && content[i] <= 'Z') || (content[i] >= '0' && content[i] <= '9') || content[i] == '#' || content[i] == '-' || content[i] == '^' || content[i] == '_' || content[i] == '.' || content[i] == '~') {
				i++
			}
			// Record leaf if not just "~" or "_" or "."? But include anyway for completeness; filter out single-char separators that are not meaningful?
			leaf := content[startLocal:i]
			// Skip "~" alone? The JS getLeaves would include "~" as atom? But its location would be leaf as well? For "bd ~ sd", leaves would include "bd", "~", "sd"? However earlier TestMiniRest expects "bd ~ sd" to give >=2 haps, not necessarily leaf count. For locations, we probably want to include "~" but JS test for collections locations didn't include "~" maybe? Check JS test for "bd ~ sd" not given. We'll include "~" as leaf for now, but filter out if leaf == "~" or "_" or "."? The JS spec for step_char includes ~, but step rule filters out "." and "_" alone via !{...}. So "." and "_" alone are not leaves. But "~" is allowed? The step_char includes "~", and step rule does `!{ const s = chars.join(""); return (s === ".") || (s === "_") }` only excludes "." and "_" , not "~". So "~" would be a leaf atom. Should we include "~"? For locations, maybe not needed, but we can include it; transpiler leaf locations for "~" might be expected? Not sure. We'll include "~" leaves but they are single char.
			// For now, include all leaves including "~"
			if leaf != "" {
				locs = append(locs, [2]int{contentOffset + startLocal, contentOffset + i})
			}
			continue
		}
		// Other chars (e.g., unicode letters) - treat as leaf start
		// For simplicity, skip one
		i++
	}
	return locs
}

func Mini2AST(code string) (string, error) {
	return code, nil
}

func MustMini(input string) core.Pattern {
	pat := Mini(input)
	if pat.Query == nil {
		panic(fmt.Sprintf("mini parse failed for %q", input))
	}
	return pat
}

// RegisterStringParser sets core string parser to mini
func RegisterStringParser() {
	core.SetStringParser(func(s string) core.Pattern {
		return Mini(s)
	})
}
