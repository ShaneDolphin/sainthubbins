// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/webaudio/*, superdough/* — AudioContext interface (native/WASM).

package audio

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// AudioContext is the Go abstraction over Web Audio / native (oto/malgo).
type AudioContext interface {
	SampleRate() int
	CurrentTime() float64
	RenderOffline(pattern core.Pattern, cycles int) ([]float32, error)
	Play(pattern core.Pattern) error
	Stop() error
}

// OfflineRenderer renders patterns to float32 samples (mono for now).
type OfflineRenderer struct {
	SampleRateValue int
	BPM             float64
}

func NewOfflineRenderer(sampleRate int) *OfflineRenderer {
	if sampleRate == 0 {
		sampleRate = 48000
	}
	return &OfflineRenderer{SampleRateValue: sampleRate, BPM: 120}
}

func (r *OfflineRenderer) SampleRate() int { return r.SampleRateValue }
func (r *OfflineRenderer) CurrentTime() float64 { return 0 }
func (r *OfflineRenderer) Play(pattern core.Pattern) error { return nil }
func (r *OfflineRenderer) Stop() error { return nil }

// RenderOffline renders pattern for cycles to mono float32.
// Each hap's duration determines a simple sine tone (for demo).
func (r *OfflineRenderer) RenderOffline(pattern core.Pattern, cycles int) ([]float32, error) {
	// Query pattern for cycles 0..cycles
	haps := pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(int64(cycles)))
	totalSamples := r.SampleRateValue * cycles * 2 // 2 sec per cycle at 0.5cps? simplified 1 sec per cycle
	if totalSamples <= 0 {
		totalSamples = r.SampleRateValue
	}
	buf := make([]float32, totalSamples)
	cps := 0.5 // default
	for _, h := range haps {
		if h.Whole == nil {
			continue
		}
		startSec := h.Whole.Begin.Float64() / cps
		endSec := h.Whole.End.Float64() / cps
		durSec := endSec - startSec
		if durSec <= 0 {
			continue
		}
		startSamp := int(startSec * float64(r.SampleRateValue))
		endSamp := int(endSec * float64(r.SampleRateValue))
		if startSamp < 0 {
			startSamp = 0
		}
		if endSamp > len(buf) {
			endSamp = len(buf)
		}
		// Frequency from hap: freq/n/note/s → freq (superdough offline fidelity)
		freq := 220.0
		gain := 0.3
		cutoff := 0.0
		mVal, isMap := h.Value.(map[string]any)
		if isMap {
			if f, ok := mVal["freq"]; ok {
				switch v := f.(type) {
				case float64:
					freq = v
				case int:
					freq = float64(v)
				case int64:
					freq = float64(v)
				case string:
					if fv, err := strconv.ParseFloat(v, 64); err == nil {
						freq = fv
					} else {
						freq = noteToFreq(v)
					}
				}
			}
			if n, ok := mVal["n"]; ok && freq == 220.0 {
				switch v := n.(type) {
				case string:
					if v != "" {
						freq = noteToFreq(v)
					}
				case float64:
					freq = 440 * math.Pow(2, (v-69)/12)
				case int:
					freq = 440 * math.Pow(2, (float64(v)-69)/12)
				case int64:
					freq = 440 * math.Pow(2, (float64(v)-69)/12)
				}
			}
			if note, ok := mVal["note"]; ok && freq == 220.0 {
				switch v := note.(type) {
				case string:
					freq = noteToFreq(v)
				case float64:
					freq = 440 * math.Pow(2, (v-69)/12)
				case int:
					freq = 440 * math.Pow(2, (float64(v)-69)/12)
				case int64:
					freq = 440 * math.Pow(2, (float64(v)-69)/12)
				}
			}
			if s, ok := mVal["s"]; ok && freq == 220.0 {
				switch v := s.(type) {
				case string:
					switch v {
					case "bd":
						freq = 60
					case "sd":
						freq = 180
					case "hh", "oh", "ch":
						freq = 800
					case "cp":
						freq = 300
					case "cicada", "birds":
						freq = 1200
					default:
						// Try note parsing for sample names like "c4"
						if nf := noteToFreq(v); nf != 220 {
							freq = nf
						} else {
							freq = 220 + float64(len(v))*10
						}
					}
				}
			}
			if g, ok := mVal["gain"]; ok {
				switch v := g.(type) {
				case float64:
					gain = v * 0.3
				case int:
					gain = float64(v) * 0.3
				case int64:
					gain = float64(v) * 0.3
				case string:
					if fv, err := strconv.ParseFloat(v, 64); err == nil {
						gain = fv * 0.3
					}
				}
				if gain < 0 {
					gain = 0
				}
				if gain > 2 {
					gain = 2
				}
			}
			if cf, ok := mVal["cutoff"]; ok {
				switch v := cf.(type) {
				case float64:
					cutoff = v
				case int:
					cutoff = float64(v)
				case int64:
					cutoff = float64(v)
				case string:
					if fv, err := strconv.ParseFloat(v, 64); err == nil {
						cutoff = fv
					}
				}
			} else if cf, ok := mVal["lpf"]; ok {
				switch v := cf.(type) {
				case float64:
					cutoff = v
				case int:
					cutoff = float64(v)
				}
			}
		} else if s, ok := h.Value.(string); ok {
			switch s {
			case "bd":
				freq = 60
			case "sd":
				freq = 180
			case "hh":
				freq = 800
			default:
				freq = noteToFreq(s)
				if freq == 220 && s != "sd" && s != "bd" {
					freq = 220
				}
			}
		}
		// One-pole low-pass for cutoff if present
		var alpha float64 = 1
		var lpPrev float64
		if cutoff > 0 && cutoff < 20000 {
			rc := 1 / (2 * math.Pi * cutoff)
			dt := 1 / float64(r.SampleRateValue)
			alpha = dt / (rc + dt)
		}
		for i := startSamp; i < endSamp; i++ {
			t := float64(i-startSamp) / float64(r.SampleRateValue)
			// ADSR envelope simplified
			env := 1.0
			if durSec > 0.01 {
				attack := 0.01
				release := 0.05
				if t < attack {
					env = t / attack
				} else if t > durSec-release {
					env = (durSec - t) / release
					if env < 0 {
						env = 0
					}
				}
			}
			raw := math.Sin(2*math.Pi*freq*t) * env * gain
			if cutoff > 0 && cutoff < 20000 {
				lpPrev = lpPrev + alpha*(raw-lpPrev)
				raw = lpPrev
			}
			buf[i] += float32(raw)
		}
	}
	return buf, nil
}

// WriteWAV writes float32 mono buffer to WAV file.
func WriteWAV(path string, samples []float32, sampleRate int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// WAV header
	dataSize := len(samples) * 2 // 16-bit
	// RIFF
	_, _ = f.Write([]byte("RIFF"))
	_ = binary.Write(f, binary.LittleEndian, uint32(36+dataSize))
	_, _ = f.Write([]byte("WAVE"))
	_, _ = f.Write([]byte("fmt "))
	_ = binary.Write(f, binary.LittleEndian, uint32(16))
	_ = binary.Write(f, binary.LittleEndian, uint16(1)) // PCM
	_ = binary.Write(f, binary.LittleEndian, uint16(1)) // mono
	_ = binary.Write(f, binary.LittleEndian, uint32(sampleRate))
	_ = binary.Write(f, binary.LittleEndian, uint32(sampleRate*2))
	_ = binary.Write(f, binary.LittleEndian, uint16(2))
	_ = binary.Write(f, binary.LittleEndian, uint16(16))
	_, _ = f.Write([]byte("data"))
	_ = binary.Write(f, binary.LittleEndian, uint32(dataSize))
	for _, s := range samples {
		v := int16(math.Max(-32768, math.Min(32767, float64(s*32767))))
		_ = binary.Write(f, binary.LittleEndian, v)
	}
	return nil
}

// noteToFreq maps note string like "c4", "a3" to frequency (A4=440).
func noteToFreq(note string) float64 {
	// Simple parser: letter + optional #/b + octave
	if len(note) < 2 {
		return 220
	}
	// Use tonal-ish mapping: c=0, c#=1, d=2, d#=3, e=4, f=5, f#=6, g=7, g#=8, a=9, a#=10, b=11
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := strings.ToLower(note)
	letter := n[0]
	semi, ok := base[letter]
	if !ok {
		return 220
	}
	idx := 1
	if len(n) > 2 && (n[1] == '#' || n[1] == 'b') {
		if n[1] == '#' {
			semi++
		} else {
			semi--
		}
		idx = 2
	}
	octave := 4
	if idx < len(n) {
		if o, err := strconv.Atoi(n[idx:]); err == nil {
			octave = o
		}
	}
	midi := (octave+1)*12 + semi // C4 = 60
	return 440 * math.Pow(2, float64(midi-69)/12)
}

// Ensure interfaces
var _ AudioContext = (*OfflineRenderer)(nil)

// RenderPatternAudio is the Go port of JS renderPatternAudio (OfflineAudioContext).
func RenderPatternAudio(pattern core.Pattern, cycles int, sampleRate int) ([]float32, error) {
	r := NewOfflineRenderer(sampleRate)
	return r.RenderOffline(pattern, cycles)
}
