// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Full DSP engine (sampler, synth, effects, worklets) — interface only.

package superdough

import "math"

// Engine is superdough DSP stub
type Engine struct {
	SampleRate int
}

func New(sampleRate int) *Engine {
	if sampleRate == 0 {
		sampleRate = 48000
	}
	return &Engine{SampleRate: sampleRate}
}

func (e *Engine) Trigger(note string, dur float64) []float32 {
	return e.TriggerWithControls(map[string]any{"note": note}, dur)
}

// TriggerWithControls renders with map controls (freq, note, n, gain, cutoff etc.) mirroring superdough.
func (e *Engine) TriggerWithControls(ctrl map[string]any, dur float64) []float32 {
	n := int(dur * float64(e.SampleRate))
	if n <= 0 {
		return []float32{}
	}
	buf := make([]float32, n)
	// Determine freq from controls
	freq := 220.0
	if f, ok := ctrl["freq"]; ok {
		switch v := f.(type) {
		case float64:
			freq = v
		case int:
			freq = float64(v)
		case string:
			freq = parseFreq(v)
		}
	} else if note, ok := ctrl["note"]; ok {
		switch v := note.(type) {
		case string:
			freq = parseFreq(v)
		case float64:
			freq = 440 * pow2((v-69)/12)
		case int:
			freq = 440 * pow2((float64(v)-69)/12)
		}
	} else if s, ok := ctrl["s"]; ok {
		if sv, ok := s.(string); ok {
			freq = sampleToFreq(sv)
		}
	} else if s, ok := ctrl["sound"]; ok {
		if sv, ok := s.(string); ok {
			freq = sampleToFreq(sv)
		}
	}
	// Gain and ADSR controls
	gain := 0.2
	if g, ok := ctrl["gain"]; ok {
		switch v := g.(type) {
		case float64:
			gain = v
		case int:
			gain = float64(v)
		case int64:
			gain = float64(v)
		}
	}
	attackSec := 0.01
	if v, ok := ctrl["attack"]; ok {
		attackSec = toFloatCtrl(v)
	} else if v, ok := ctrl["att"]; ok {
		attackSec = toFloatCtrl(v)
	}
	decaySec := 0.05
	if v, ok := ctrl["decay"]; ok {
		decaySec = toFloatCtrl(v)
	} else if v, ok := ctrl["dec"]; ok {
		decaySec = toFloatCtrl(v)
	}
	sustain := 0.7
	if v, ok := ctrl["sustain"]; ok {
		sustain = toFloatCtrl(v)
	} else if v, ok := ctrl["sus"]; ok {
		sustain = toFloatCtrl(v)
	}
	releaseSec := 0.05
	if v, ok := ctrl["release"]; ok {
		releaseSec = toFloatCtrl(v)
	} else if v, ok := ctrl["rel"]; ok {
		releaseSec = toFloatCtrl(v)
	}
	if attackSec < 0 {
		attackSec = 0.01
	}
	if decaySec < 0 {
		decaySec = 0.05
	}
	if releaseSec < 0 {
		releaseSec = 0.05
	}
	attack := int(attackSec * float64(e.SampleRate))
	decay := int(decaySec * float64(e.SampleRate))
	release := int(releaseSec * float64(e.SampleRate))
	// n sample index detune: slight freq shift per n
	if nv, ok := ctrl["n"]; ok {
		nvF := toFloatCtrl(nv)
		freq = freq * (1 + nvF*0.02)
	}
	if attack > n/3 {
		attack = n / 3
	}
	if decay > n/3 {
		decay = n / 3
	}
	if release > n/3 {
		release = n / 3
	}
	// One-pole lowpass state for cutoff
	cutoffFreq := 0.0
	if v, ok := ctrl["cutoff"]; ok {
		cutoffFreq = toFloatCtrl(v)
	} else if v, ok := ctrl["lpf"]; ok {
		cutoffFreq = toFloatCtrl(v)
	} else if v, ok := ctrl["ctf"]; ok {
		cutoffFreq = toFloatCtrl(v)
	}
	// Simple one-pole filter coefficient: alpha = dt / (RC + dt), RC=1/(2pi*fc)
	var alpha float64 = 1
	var prev float64
	if cutoffFreq > 0 && cutoffFreq < 20000 {
		rc := 1 / (2 * 3.141592653589793 * cutoffFreq)
		dt := 1 / float64(e.SampleRate)
		alpha = dt / (rc + dt)
		if alpha > 1 {
			alpha = 1
		}
		if alpha < 0 {
			alpha = 0
		}
	}
	// Pan: -1 left, 1 right, 0 center — for mono we just scale gain slightly
	pan := 0.0
	if v, ok := ctrl["pan"]; ok {
		pan = toFloatCtrl(v)
		if pan < -1 {
			pan = -1
		}
		if pan > 1 {
			pan = 1
		}
	}
	panGain := 1.0 - 0.3*math.Abs(pan) // crude center louder

	for i := range buf {
		t := float64(i) / float64(e.SampleRate)
		raw := sinApprox(freq*t*6.28318530718) + 0.3*sinApprox(freq*t*6.28318530718*2) // mild harmonics
		// ADSR envelope with sustain level
		env := gain
		if attack > 0 && i < attack {
			env = gain * float64(i) / float64(attack)
		} else if decay > 0 && i < attack+decay {
			prog := float64(i-attack) / float64(decay)
			env = gain*(1-prog) + gain*sustain*prog
		} else if i >= n-release {
			if release > 0 {
				prog := float64(i-(n-release)) / float64(release)
				env = gain * sustain * (1 - prog)
			} else {
				env = 0
			}
		} else {
			env = gain * sustain
		}
		env *= panGain
		// Apply one-pole lowpass if cutoff present
		filtered := raw
		if cutoffFreq > 0 && cutoffFreq < 20000 {
			prev = prev + alpha*(raw-prev)
			filtered = prev
		} else {
			filtered = raw
		}
		buf[i] = float32(filtered * env)
	}
	return buf
}

func parseFreq(s string) float64 {
	if len(s) == 0 {
		return 220
	}
	// Try numeric
	if f, err := parseFloat(s); err == nil {
		if f > 20 && f < 20000 {
			return f
		}
	}
	// Note parsing: c, c#, db, etc + octave — mirrors audio.noteToFreq
	n := toLower(s)
	if len(n) >= 2 {
		base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
		if semi, ok := base[n[0]]; ok {
			idx := 1
			if len(n) > 2 && (n[1] == '#' || n[1] == 'b' || n[1] == 's') {
				if n[1] == '#' || n[1] == 's' {
					semi++
				} else {
					semi--
				}
				idx = 2
			}
			octave := 4
			if idx < len(n) {
				if o, err := parseInt(n[idx:]); err == nil {
					octave = o
				}
			}
			midi := (octave+1)*12 + semi
			return 440 * pow2(float64(midi-69)/12)
		}
	}
	// Fallback hash for unknown samples
	var h float64
	for _, c := range s {
		h += float64(c%7) * 5
	}
	return 220 + h
}

func parseFloat(s string) (float64, error) {
	// local without strconv import overhead — use manual
	var neg bool
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}
	var intPart float64
	var fracPart float64
	var fracDiv float64 = 1
	dot := false
	for _, c := range s {
		if c == '.' {
			dot = true
			continue
		}
		if c < '0' || c > '9' {
			return 0, errParse
		}
		d := float64(c - '0')
		if !dot {
			intPart = intPart*10 + d
		} else {
			fracDiv *= 10
			fracPart += d / fracDiv
		}
	}
	v := intPart + fracPart
	if neg {
		v = -v
	}
	return v, nil
}

func parseInt(s string) (int, error) {
	if len(s) == 0 {
		return 0, errParse
	}
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}
	v := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errParse
		}
		v = v*10 + int(c-'0')
	}
	if neg {
		v = -v
	}
	return v, nil
}

func toLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

var errParse = errorString("parse error")

type errorString string

func (e errorString) Error() string { return string(e) }

func sampleToFreq(s string) float64 {
	switch s {
	case "bd", "bassdrum":
		return 60
	case "sd", "snare":
		return 180
	case "hh", "oh", "ch", "hat":
		return 800
	case "cp", "clap":
		return 300
	case " Rim", "rim":
		return 400
	case "kick":
		return 55
	case "bass":
		return 80
	case "arp", "arpy":
		return 320
	case "drum":
		return 110
	case "hh3", "hh4":
		return 900
	default:
		// Try note-like sample "c3" etchandled by parseFreq; fallback to 220 + hash
		if len(s) >= 2 && s[0] >= 'a' && s[0] <= 'g' {
			return parseFreq(s)
		}
		return 220
	}
}

func pow2(x float64) float64 {
	return math.Pow(2, x)
}

func sinApprox(x float64) float64 {
	return math.Sin(x)
}

func toFloatCtrl(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		if f, err := parseFloat(x); err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
}
