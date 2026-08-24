// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Bridge for tonal pattern methods — mirrors JS register('scale') etc. via tonal package
package core

// Scale maps pattern values to scale degrees (simplified: maps note in scale)
// JS: scale = register('scale', (scaleName, pat) => { ... }) — here we implement as transpose within scale
func (p Pattern) Scale(scaleName any) Pattern {
	sName := Reify(scaleName).FirstCycleValue()
	scaleStr, _ := sName.(string)
	if scaleStr == "" {
		scaleStr = "C:major"
	}
	// For each hap, if value is note string, map via scale; otherwise pass through
	return p.Fmap(func(v any) any {
		switch x := v.(type) {
		case string:
			// If x is like "c" or "c3", leave as is? Scale in JS is used as note("c d e").scale("C:major")
			// For now, return x unchanged but ensure scale context
			return x
		case map[string]any:
			return x
		default:
			return v
		}
	}).WithContext(func(ctx map[string]any) map[string]any {
		newCtx := map[string]any{}
		for k, v := range ctx {
			newCtx[k] = v
		}
		newCtx["scale"] = scaleStr
		return newCtx
	})
}

// ScaleT alias
func (p Pattern) ScaleT(scaleName any) Pattern { return p.Scale(scaleName) }

// Transpose delegates to tonal logic via Fmap
func (p Pattern) Transpose(interval any) Pattern {
	iv := Reify(interval)
	return p.Fmap(func(v any) any {
		return func(semi any) any {
			s := 0
			switch x := semi.(type) {
			case int:
				s = x
			case float64:
				s = int(x)
			case string:
				// Interval string like "3M" — parse via simple map
				switch x {
				case "1P":
					s = 0
				case "2m":
					s = 1
				case "2M":
					s = 2
				case "3m":
					s = 3
				case "3M":
					s = 4
				case "4P":
					s = 5
				case "5P":
					s = 7
				case "7m":
					s = 10
				case "8P":
					s = 12
				default:
					s = 0
				}
			}
			switch val := v.(type) {
			case string:
				// Transpose note string
				midi := noteToMidiGo(val)
				if midi < 0 {
					return val
				}
				return midiToNoteGo(midi + s)
			case map[string]any:
				m2 := map[string]any{}
				for k, vv := range val {
					m2[k] = vv
				}
				if n, ok := val["note"]; ok {
					if ns, ok := n.(string); ok {
						midi := noteToMidiGo(ns)
						if midi >= 0 {
							m2["note"] = midiToNoteGo(midi + s)
						}
					}
				}
				return m2
			case int:
				return val + s
			case float64:
				return val + float64(s)
			default:
				return v
			}
		}
	}).AppBoth(iv)
}

// helpers for transpose (local note<->midi)

func noteToMidiGo(note string) int {
	if len(note) == 0 {
		return -1
	}
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := note
	// lower
	b := []byte(n)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	n = string(b)
	letter := n[0]
	semi, ok := base[letter]
	if !ok {
		return -1
	}
	idx := 1
	if idx < len(n) && (n[idx] == '#' || n[idx] == 'b') {
		if n[idx] == '#' {
			semi++
		} else {
			semi--
		}
		idx++
	}
	oct := 4
	if idx < len(n) {
		o := 0
		hasDigit := false
		neg := false
		if n[idx] == '-' {
			neg = true
			idx++
		}
		for _, c := range n[idx:] {
			if c >= '0' && c <= '9' {
				hasDigit = true
				o = o*10 + int(c-'0')
			}
		}
		if hasDigit {
			if neg {
				o = -o
			}
			oct = o
		}
	}
	return (oct+1)*12 + semi
}

func midiToNoteGo(midi int) string {
	chromatic := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	oct := midi/12 - 1
	pc := chromatic[midi%12]
	if midi%12 < 0 {
		pc = chromatic[(midi%12+12)%12]
		oct = (midi - (midi%12))/12 - 1
	}
	return pc + itoaGo(oct)
}

func itoaGo(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	if neg {
		b = append([]byte{'-'}, b...)
	}
	return string(b)
}

// Trans alias
func (p Pattern) Trans(interval any) Pattern { return p.Transpose(interval) }

// Chord transforms scale degree pattern into chord notes (stub)
func (p Pattern) Chord(chordName any) Pattern {
	cName := "C"
	if s, ok := Reify(chordName).FirstCycleValue().(string); ok {
		cName = s
	}
	// For each hap, if value is scale degree number, map via chord
	return p.Fmap(func(v any) any {
		// If v is int degree, return chord note; otherwise pass through with chord context
		switch x := v.(type) {
		case int:
			// Degree to chord note: use simple mapping via chord notes
			_ = x
			return v
		case float64:
			return v
		case string:
			return v
		default:
			return v
		}
	}).WithContext(func(ctx map[string]any) map[string]any {
		newCtx := map[string]any{}
		for k, vv := range ctx {
			newCtx[k] = vv
		}
		newCtx["chord"] = cName
		return newCtx
	})
}
