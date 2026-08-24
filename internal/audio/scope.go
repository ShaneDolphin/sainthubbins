// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package audio

// Scope is audio analyser stub
type Scope struct {
	Enabled bool
}

func NewScope() *Scope { return &Scope{} }

func (s *Scope) Analyse(samples []float32) []float32 {
	if len(samples) == 0 {
		return nil
	}
	// RMS per 1024 block + spectrum stub (peak)
	block := 1024
	if len(samples) < block {
		block = len(samples)
	}
	var sum float64
	for i := 0; i < block; i++ {
		sum += float64(samples[i]) * float64(samples[i])
	}
	rms := float32(sum / float64(block))
	if rms < 0 {
		rms = -rms
	}
	peak := float32(0)
	for i := 0; i < block; i++ {
		if samples[i] > peak {
			peak = samples[i]
		}
		if -samples[i] > peak {
			peak = -samples[i]
		}
	}
	return []float32{rms, peak, float32(len(samples))}
}
