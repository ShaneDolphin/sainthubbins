// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

import "strings"

// DrawLine renders pattern as string with cycle separators (|), holds (-), silence (.)
func DrawLine(pat Pattern, chars int) string {
	if chars <= 0 {
		chars = 60
	}
	cycle := FractionFromInt(0)
	pos := FractionFromInt(0)
	lines := []string{""}
	emptyLine := ""
	for len(lines[0]) < chars {
		haps := pat.QueryArc(cycle, cycle.Add(FractionFromInt(1)))
		var durations []Fraction
		for _, h := range haps {
			if h.HasOnset() {
				durations = append(durations, h.Duration())
			}
		}
		var charFrac Fraction
		if len(durations) == 0 {
			charFrac = FractionFromInt(1)
		} else {
			// gcd of durations
			charFrac = durations[0]
			for _, d := range durations[1:] {
				charFrac = GcdFraction(charFrac, d)
			}
		}
		if charFrac.Equals(FractionFromInt(0)) {
			charFrac = FractionFromInt(1)
		}
		// inverse = 1/charFrac
		totalSlots := FractionFromInt(1).Div(charFrac).Num // simplified, assume denominator 1?
		// Actually charFrac inverse numerator/denominator
		if charFrac.Num != 0 {
			totalSlots = charFrac.Den / charFrac.Num
			if totalSlots <= 0 {
				totalSlots = int64(1)
			}
		}
		for i := range lines {
			lines[i] += "|"
		}
		emptyLine += "|"
		for i := int64(0); i < totalSlots; i++ {
			begin := pos
			end := pos.Add(charFrac)
			span := NewTimeSpan(begin, end)
			// Find hap covering this slot
			found := "."
			for _, h := range haps {
				if h.Whole != nil && h.Whole.Begin.Equals(begin) {
					// onset
					found = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(begin.Show(), " ", ""), "/", ""), ".", ""), "0", ""), "1", "")
					// Simpler: use value string first char
					s := ""
					switch val := h.Value.(type) {
					case string:
						if len(val) > 0 {
							s = string(val[0])
						} else {
							s = "?"
						}
					default:
						s = strings.TrimSpace(strings.Split(strings.ReplaceAll(span.Show(), " ", ""), "")[0])
						if s == "" {
							s = "x"
						}
						if len(s) > 1 {
							s = string(s[0])
						}
					}
					found = s
					break
				} else if h.Part.Intersection(span) != nil && h.Whole != nil && !h.Whole.Begin.Equals(begin) && !h.HasOnset() {
					// hold
					found = "-"
					break
				}
			}
			for idx := range lines {
				// Only first line for now
				if idx == 0 {
					lines[idx] += found
				} else {
					lines[idx] += " "
				}
			}
			pos = end
		}
		cycle = cycle.Add(FractionFromInt(1))
		if len(lines[0]) >= chars {
			break
		}
		// Prevent infinite
		if cycle.Float64() > 100 {
			break
		}
	}
	result := strings.Join(lines, "\n")
	if len(result) > chars {
		result = result[:chars]
	}
	return result
}
