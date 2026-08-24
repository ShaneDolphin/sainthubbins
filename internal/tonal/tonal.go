// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Stub for music theory (scales, chords, voicings) — uses @tonaljs/tonal equivalent via map.

package tonal

import "strings"

// Scale returns notes for a scale name (e.g., "C:major", "A:minor")
func Scale(name string) []string {
	parts := strings.Split(name, ":")
	root := "C"
	scaleType := "major"
	if len(parts) > 0 && parts[0] != "" {
		root = parts[0]
	}
	if len(parts) > 1 {
		scaleType = parts[1]
	}
	// Normalize root upper
	if len(root) > 0 {
		root = strings.ToUpper(root[:1]) + strings.ToLower(root[1:])
	}
	// Interval maps (semitones) for common scales — expanded from @tonaljs/scale
	intervals := map[string][]int{
		"major":        {0, 2, 4, 5, 7, 9, 11},
		"minor":        {0, 2, 3, 5, 7, 8, 10},
		"dorian":       {0, 2, 3, 5, 7, 9, 10},
		"mixolydian":   {0, 2, 4, 5, 7, 9, 10},
		"phrygian":     {0, 1, 3, 5, 7, 8, 10},
		"lydian":       {0, 2, 4, 6, 7, 9, 11},
		"locrian":      {0, 1, 3, 5, 6, 8, 10},
		"harmonic minor": {0, 2, 3, 5, 7, 8, 11},
		"melodic minor": {0, 2, 3, 5, 7, 9, 11},
		"whole tone":   {0, 2, 4, 6, 8, 10},
		"pentatonic":   {0, 2, 4, 7, 9},
		"blues":        {0, 3, 5, 6, 7, 10},
		"chromatic":    {0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		"aeolian":      {0, 2, 3, 5, 7, 8, 10},
		"ionian":       {0, 2, 4, 5, 7, 9, 11},
	}
	// Normalize scaleType: lower, trim spaces
	scaleKey := strings.ToLower(strings.TrimSpace(scaleType))
	ivs, ok := intervals[scaleKey]
	if !ok {
		// Try alias without spaces
		scaleKey = strings.ReplaceAll(scaleKey, " ", "")
		ivs, ok = intervals[scaleKey]
		if !ok {
			ivs = intervals["major"]
		}
	}
	chromatic := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	// Find root index
	rootIdx := 0
	for i, n := range chromatic {
		if strings.EqualFold(n, root) || strings.EqualFold(n, strings.TrimSuffix(root, "#")) {
			rootIdx = i
			if strings.HasSuffix(root, "#") && !strings.Contains(n, "#") {
				continue
			}
			if strings.EqualFold(n, root) {
				break
			}
		}
	}
	out := make([]string, len(ivs))
	for i, iv := range ivs {
		out[i] = chromatic[(rootIdx+iv)%12]
	}
	return out
}

// Chord returns notes for chord symbol (basic triad parsing) — expanded with root parsing
func Chord(name string) []string {
	low := strings.TrimSpace(name)
	if low == "" {
		return []string{"C", "E", "G"}
	}
	// Extract root: letter A-G + optional #/b + optional octave digits
	root := "C"
	rest := ""
	i := 0
	if len(low) > 0 && ((low[0] >= 'A' && low[0] <= 'G') || (low[0] >= 'a' && low[0] <= 'g')) {
		root = string(low[0])
		i = 1
		if i < len(low) && (low[i] == '#' || low[i] == 'b' || low[i] == 'B') {
			root += string(low[i])
			i++
		}
		// Normalize root case: C, C#, Db etc.
		if len(root) > 0 {
			root = strings.ToUpper(root[:1]) + strings.ToLower(root[1:])
		}
		rest = strings.ToLower(strings.TrimSpace(low[i:]))
		// Strip leading ":" or " " before quality
		rest = strings.TrimPrefix(rest, ":")
		rest = strings.TrimSpace(rest)
	} else {
		rest = strings.ToLower(low)
	}
	// Determine quality intervals relative to root major: major=0,4,7
	intervals := []int{0, 4, 7} // default major
	switch {
	case strings.Contains(rest, "maj7"):
		intervals = []int{0, 4, 7, 11}
	case strings.Contains(rest, "m7") && !strings.Contains(rest, "maj"):
		intervals = []int{0, 3, 7, 10}
	case strings.Contains(rest, "7") && !strings.Contains(rest, "maj"):
		intervals = []int{0, 4, 7, 10}
	case strings.Contains(rest, "dim7"):
		intervals = []int{0, 3, 6, 9}
	case strings.Contains(rest, "dim") || rest == "o" || strings.Contains(rest, "o7"):
		intervals = []int{0, 3, 6}
	case strings.Contains(rest, "aug") || strings.Contains(rest, "+"):
		intervals = []int{0, 4, 8}
	case strings.Contains(rest, "sus4"):
		intervals = []int{0, 5, 7}
	case strings.Contains(rest, "sus2"):
		intervals = []int{0, 2, 7}
	case strings.Contains(rest, "minor") || rest == "m" || strings.HasPrefix(rest, "m ") || strings.HasPrefix(rest, "min"):
		intervals = []int{0, 3, 7}
	case strings.Contains(rest, "m") && !strings.Contains(rest, "maj"):
		// Check if rest starts with m or contains "m" as quality
		if strings.HasPrefix(rest, "m") {
			intervals = []int{0, 3, 7}
		}
	}
	// Build notes from root + intervals
	chromatic := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	// Find root index handling enharmonics like Db -> C#
	rootNorm := root
	if len(root) == 2 && root[1] == 'b' {
		// Db -> C#, Eb -> D#, etc.
		m := map[string]string{"Db": "C#", "Eb": "D#", "Gb": "F#", "Ab": "G#", "Bb": "A#"}
		if v, ok := m[root]; ok {
			rootNorm = v
		}
	}
	rootIdx := 0
	for i, n := range chromatic {
		if strings.EqualFold(n, rootNorm) {
			rootIdx = i
			break
		}
	}
	out := make([]string, len(intervals))
	for i, iv := range intervals {
		out[i] = chromatic[(rootIdx+iv)%12]
	}
	// Preserve original root spelling for first note (e.g., Db not C# if input was Db)
	out[0] = root
	return out
}

// Voicing selects inversion
func Voicing(chord []string, name string) []string {
	low := strings.ToLower(name)
	switch {
	case strings.Contains(low, "drop2") && len(chord) >= 4:
		// drop2: move second from top down an octave — simplified as rotate
		return append(chord[1:], chord[0])
	case strings.Contains(low, "drop3") && len(chord) >= 4:
		return append(chord[2:], chord[:2]...)
	case strings.Contains(low, "first") && len(chord) > 1:
		return append(chord[1:], chord[0])
	case strings.Contains(low, "second") && len(chord) > 2:
		return append(chord[2:], chord[:2]...)
	case strings.Contains(low, "third") && len(chord) > 3:
		return append(chord[3:], chord[:3]...)
	case strings.Contains(low, "root") || strings.Contains(low, "close"):
		return chord
	}
	return chord
}

// Transpose transposes a note by semitones or interval string (e.g., "4P", "3m")
func Transpose(note string, interval any) string {
	if note == "" {
		return note
	}
	// Parse note
	n := strings.TrimSpace(note)
	// Determine semitones
	var semis int
	switch v := interval.(type) {
	case int:
		semis = v
	case int64:
		semis = int(v)
	case float64:
		semis = int(v)
	case string:
		s := strings.TrimSpace(v)
		// Try interval string like "3m", "4P", "5d"
		if iv, ok := intervalSemitones(s); ok {
			semis = iv
		} else if iv2, err := parseInt(s); err == nil {
			semis = iv2
		}
	default:
		// fallback
		semis = 0
	}
	// Convert note to midi, transpose, back to note
	midi := noteToMidi(n)
	if midi < 0 {
		return note
	}
	midi += semis
	return midiToNote(midi)
}

func intervalSemitones(s string) (int, bool) {
	// Handle leading - for descending intervals
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}
	// Map common intervals
	m := map[string]int{
		"1P": 0, "2m": 1, "2M": 2, "3m": 3, "3M": 4, "4P": 5, "4A": 6, "5d": 6, "5P": 7, "6m": 8, "6M": 9, "7m": 10, "7M": 11, "8P": 12,
		"9m": 13, "9M": 14, "11P": 17, "12P": 19,
	}
	if v, ok := m[s]; ok {
		if neg {
			v = -v
		}
		return v, true
	}
	// Try "m2", "M3" variant without number
	if len(s) == 2 {
		// e.g., "m3" -> "3m"
		alt := string(s[1]) + string(s[0])
		if v, ok := m[alt]; ok {
			if neg {
				v = -v
			}
			return v, true
		}
	}
	return 0, false
}

func parseInt(s string) (int, error) {
	var neg bool
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}
	v := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errorString("parse int")
		}
		v = v*10 + int(c-'0')
	}
	if neg {
		v = -v
	}
	return v, nil
}

func noteToMidi(note string) int {
	if len(note) < 1 {
		return -1
	}
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := strings.ToLower(strings.TrimSpace(note))
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
	octave := 4
	if idx < len(n) {
		if o, err := parseInt(n[idx:]); err == nil {
			octave = o
		}
	}
	return (octave+1)*12 + semi
}

func midiToNote(midi int) string {
	chromatic := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	oct := midi/12 - 1
	pc := chromatic[midi%12]
	if pc == "" {
		pc = "C"
	}
	return pc + itoa(oct)
}

func itoa(i int) string {
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

type errorString string

func (e errorString) Error() string { return string(e) }
