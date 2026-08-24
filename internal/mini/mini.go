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
	// Handle range operator "a .. b" at sequence level (spaced "..")
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == ".." && i > 0 && i+1 < len(tokens) {
			leftStr := tokens[i-1]
			rightStr := tokens[i+1]
			if l, err1 := strconv.Atoi(leftStr); err1 == nil {
				if r, err2 := strconv.Atoi(rightStr); err2 == nil {
					var pats []core.Pattern
					if l <= r {
						for v := l; v <= r; v++ {
							pats = append(pats, core.Pure(v))
						}
					} else {
						for v := l; v >= r; v-- {
							pats = append(pats, core.Pure(v))
						}
					}
					rangePat := core.FastCat(pats...)
					newPats := make([]core.Pattern, 0, len(tokens)-2)
					for j := 0; j < i-1; j++ {
						newPats = append(newPats, parseToken(tokens[j]))
					}
					newPats = append(newPats, rangePat)
					for j := i + 2; j < len(tokens); j++ {
						newPats = append(newPats, parseToken(tokens[j]))
					}
					if len(newPats) == 0 {
						return core.Silence()
					}
					if len(newPats) == 1 {
						return newPats[0]
					}
					return core.FastCat(newPats...)
				}
			}
		}
	}
	// Handle choice operator "|" at sequence level: "a | b | c" -> Choose
	// Detect isolated "|" tokens
	hasPipe := false
	for _, t := range tokens {
		if t == "|" {
			hasPipe = true
			break
		}
	}
	if hasPipe {
		// Split by "|"
		var groups [][]string
		var cur []string
		for _, t := range tokens {
			if t == "|" {
				groups = append(groups, cur)
				cur = nil
			} else {
				cur = append(cur, t)
			}
		}
		groups = append(groups, cur)
		choices := make([]any, 0, len(groups))
		for _, g := range groups {
			if len(g) == 0 {
				continue
			}
			// Reconstruct mini string for group and parse
			choicePat := Mini(strings.Join(g, " "))
			choices = append(choices, choicePat)
		}
		if len(choices) > 0 {
			return core.Pure(0).Choose(choices)
		}
	}
	// Handle stack with commas at sequence level without curly? "a,b" already handled as token "a,b" but spaced "a , b" not
	// For consistency, also handle isolated "," tokens as stack
	hasComma := false
	for _, t := range tokens {
		if t == "," {
			hasComma = true
			break
		}
	}
	if hasComma {
		var groups [][]string
		var cur []string
		for _, t := range tokens {
			if t == "," {
				groups = append(groups, cur)
				cur = nil
			} else {
				cur = append(cur, t)
			}
		}
		groups = append(groups, cur)
		pats := make([]core.Pattern, 0, len(groups))
		for _, g := range groups {
			if len(g) == 0 {
				continue
			}
			pats = append(pats, Mini(strings.Join(g, " ")))
		}
		if len(pats) > 0 {
			return core.Stack(pats...)
		}
	}
	if len(tokens) == 1 {
		return parseToken(tokens[0])
	}
	pats := make([]core.Pattern, 0, len(tokens))
	for _, tok := range tokens {
		if tok == ".." || tok == "|" || tok == "," {
			continue
		}
		pats = append(pats, parseToken(tok))
	}
	return core.FastCat(pats...)
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

// unwrapGroup returns the inner text of a token that is wholly enclosed by a
// single bracket pair — "[a b]" -> "a b" — so the group is parsed as a unit
// before any operator scanning happens.
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
			toks := splitMiniTokens(inner)
			pats := make([]core.Pattern, len(toks))
			for i, t := range toks {
				pats[i] = parseToken(t)
			}
			return core.SlowCat(pats...)
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
		// Find closing }
		closeIdx := strings.LastIndex(tok, "}")
		if closeIdx > 0 {
			inner := tok[1:closeIdx]
			suffix := tok[closeIdx+1:]
			// inner is comma-separated sequences
			seqStrs := strings.Split(inner, ",")
			seqPats := make([]core.Pattern, 0, len(seqStrs))
			for _, s := range seqStrs {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				seqPats = append(seqPats, Mini(s))
			}
			var base core.Pattern
			if len(seqPats) == 0 {
				base = core.Silence()
			} else if len(seqPats) == 1 {
				base = seqPats[0]
			} else {
				// Polymeter: stack with steps alignment via Polymeter
				// Polymeter requires Steps; FastCat has nil Steps so falls back to Stack
				allPats := seqPats
				base = core.Polymeter(allPats...)
				// Fallback to Stack if Polymeter returned Silence due to missing Steps
				if len(base.FirstCycle()) == 0 {
					base = core.Stack(seqPats...)
				}
			}
			// Handle suffix: %n (steps-per-cycle), *n (fast), /n (slow)
			suffix = strings.TrimSpace(suffix)
			if suffix != "" {
				if strings.HasPrefix(suffix, "%") {
					if v, err := strconv.Atoi(strings.TrimSpace(suffix[1:])); err == nil && v > 0 {
						// steps-per-cycle: repeat to lcm? For "%3", JS expects [a b a, c d e] style — approximate via FastCat repeated?
						// Simplified: Fast by factor len(seqPats)?? Use Polymeter with explicit steps: not perfect, but ensure non-empty
						return base
					}
				} else if strings.HasPrefix(suffix, "*") {
					if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
						return base.FastF(core.FractionFromFloat(v))
					}
				} else if strings.HasPrefix(suffix, "/") {
					if v, err := strconv.ParseFloat(strings.TrimSpace(suffix[1:]), 64); err == nil {
						return base.SlowF(core.FractionFromFloat(v))
					}
				}
			}
			return base
		}
	}
	// Handle weight @: "bd@2" or "bd:3@2" etc — handle suffix @N and _N and ! and ?
	// Strip trailing modifiers: @/_ weight, ! replicate, ? degrade
	// For full PEG parity, these would be ElementStub ops; here we simplify to repeats/weights
	if containsAtDepth0(tok, "@_!?") {
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
		// Handle replicate !: "bd!2" or "bd!!" etc
		if containsAtDepth0(tok, "!") {
			baseEnd := indexAtDepth0(tok, "!")
			base := tok[:baseEnd]
			rest := tok[baseEnd:]
			reps := strings.Count(rest, "!") + 0
			// Also handle !N
			pat := parseToken(base)
			// Count ! and number after
			numStr := strings.Trim(rest, "!")
			if numStr != "" {
				if v, err := strconv.Atoi(numStr); err == nil {
					reps = v
				}
			}
			if reps <= 1 {
				reps = 2
			}
			// Fast replicate: repeat event reps times via Ply? For mini, ! is weight-like
			return pat.Ply(reps)
		}
		// Handle weight @ or _: "bd@2" means weight 2 (elongate) via TimeCatWeighted
		if containsAtDepth0(tok, "@_") {
			sep := byte('@')
			if containsAtDepth0(tok, "_") {
				sep = '_'
			}
			parts := splitAtDepth0(tok, sep)
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
