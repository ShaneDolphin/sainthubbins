// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/util.mjs
package core

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Note helpers

var isNoteRe = regexp.MustCompile(`^[a-gA-G][#bsf]*-?[0-9]*$`)
var isNoteWithOctaveRe = regexp.MustCompile(`^[a-gA-G][#bsf]*[0-9]*$`)
var tokenizeRe = regexp.MustCompile(`^([a-gA-G])([#bsf]*)(-?[0-9]*)$`)

var chromas = map[string]int{"c": 0, "d": 2, "e": 4, "f": 5, "g": 7, "a": 9, "b": 11}
var accs = map[rune]int{'#': 1, 'b': -1, 's': 1, 'f': -1}

func IsNote(name string) bool { return isNoteRe.MatchString(name) }
func IsNoteWithOctave(name string) bool { return isNoteWithOctaveRe.MatchString(name) }

func TokenizeNote(note string) (pc string, acc string, oct *int) {
	m := tokenizeRe.FindStringSubmatch(note)
	if m == nil {
		return "", "", nil
	}
	pc = m[1]
	acc = m[2]
	if m[3] != "" {
		n, err := strconv.Atoi(m[3])
		if err == nil {
			oct = &n
		}
	}
	return pc, acc, oct
}

func GetAccidentalsOffset(accidentals string) int {
	offset := 0
	for _, ch := range accidentals {
		if v, ok := accs[ch]; ok {
			offset += v
		}
	}
	return offset
}

func NoteToMidi(note string, defaultOctave ...int) (int, error) {
	dOct := 3
	if len(defaultOctave) > 0 {
		dOct = defaultOctave[0]
	}
	pc, acc, octPtr := TokenizeNote(note)
	if pc == "" {
		return 0, fmt.Errorf("not a note: %q", note)
	}
	chroma := chromas[strings.ToLower(pc)]
	offset := GetAccidentalsOffset(acc)
	oct := dOct
	if octPtr != nil {
		oct = *octPtr
	}
	return (oct+1)*12 + chroma + offset, nil
}

func MidiToFreq(n float64) float64 { return math.Pow(2, (n-69)/12) * 440 }
func FreqToMidi(freq float64) float64 { return (12*math.Log(freq/440))/math.Ln2 + 69 }

func Mod(n, m int) int { return ((n % m) + m) % m }
func ModFloat(n, m float64) float64 { return math.Mod(math.Mod(n, m)+m, m) }

// _mod for Fraction
func ModFrac(n, m Fraction) Fraction { return n.Mod(m) }

func NanFallback(value any, fallback float64) float64 {
	switch v := value.(type) {
	case float64:
		if math.IsNaN(v) {
			return fallback
		}
		return v
	case int:
		return float64(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil || math.IsNaN(f) {
			return fallback
		}
		return f
	default:
		return fallback
	}
}

func GetSoundIndex(n float64, numSounds int) int {
	if math.IsNaN(n) {
		n = 0
	}
	return Mod(int(math.Round(n)), numSounds)
}

func RemoveUndefineds[T any](xs []*T) []*T {
	var out []*T
	for _, x := range xs {
		if x != nil {
			out = append(out, x)
		}
	}
	return out
}

func Flatten[T any](arr [][]T) []T {
	var out []T
	for _, inner := range arr {
		out = append(out, inner...)
	}
	return out
}

func ListRange(min, max int) []int {
	if max < min {
		return []int{}
	}
	out := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		out = append(out, i)
	}
	return out
}

func ParseNumeral(s string) (float64, error) {
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f, nil
	}
	if IsNote(s) {
		midi, err := NoteToMidi(s)
		if err == nil {
			return float64(midi), nil
		}
	}
	return 0, fmt.Errorf("cannot parse as numeral: %q", s)
}

func NumeralArgs(fn func(...float64) float64) func(...string) (float64, error) {
	return func(args ...string) (float64, error) {
		floats := make([]float64, len(args))
		for i, a := range args {
			f, err := ParseNumeral(a)
			if err != nil {
				return 0, err
			}
			floats[i] = f
		}
		return fn(floats...), nil
	}
}

func SplitAt(index int, value string) (string, string) {
	if index < 0 {
		index = 0
	}
	if index > len(value) {
		index = len(value)
	}
	return value[:index], value[index:]
}

func ZipWith[A, B, C any](f func(A, B) C, xs []A, ys []B) []C {
	n := len(xs)
	if len(ys) < n {
		n = len(ys)
	}
	out := make([]C, n)
	for i := 0; i < n; i++ {
		out[i] = f(xs[i], ys[i])
	}
	return out
}

func Pairs[T any](xs []T) [][2]T {
	var out [][2]T
	for i := 0; i < len(xs)-1; i++ {
		out = append(out, [2]T{xs[i], xs[i+1]})
	}
	return out
}

func Clamp(num, min, max float64) float64 { return math.Min(math.Max(num, min), max) }

func StringifyValues(value any, compact bool) string {
	switch v := value.(type) {
	case map[string]any:
		parts := []string{}
		for k, val := range v {
			parts = append(parts, fmt.Sprintf("%s:%v", k, val))
		}
		if compact {
			return strings.Join(parts, " ")
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", value)
	}
}
