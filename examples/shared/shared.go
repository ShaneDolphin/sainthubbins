// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// examples/shared — small helpers shared by the tutorial templates.
//
// These exist so each template file can be almost entirely music. Nothing
// here is required to use Saint Hubbins; it is ten lines of convenience.

package shared

import (
	"fmt"
	"os"

	"codeberg.org/uzu/saint-hubbins/internal/audio"
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// SampleRate used for every offline render in the tutorial.
const SampleRate = 48000

// Tempo converts beats-per-minute into a Fast factor.
//
// One cycle of a pattern is one bar, and the offline renderer runs at a fixed
// 0.5 cycles per second — two seconds per bar. Four beats in two seconds is
// 120 BPM, so 120 is the speed of an untouched pattern and every other tempo
// is a ratio against it: Tempo(140) is 140/120, or 7/6.
func Tempo(bpm int64) core.Fraction {
	return core.NewFraction(bpm, 120)
}

// Render writes bars of p to path as a WAV file and prints what it produced.
//
// It reports the number of events and the peak level so a silent or clipping
// render is obvious without opening the file.
func Render(p core.Pattern, path string, bars int) error {
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(int64(bars)))

	samples, err := audio.RenderPatternAudio(p, bars, SampleRate)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}
	if err := audio.WriteWAV(path, samples, SampleRate); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	var peak float32
	for _, s := range samples {
		if s > peak {
			peak = s
		} else if -s > peak {
			peak = -s
		}
	}
	seconds := float64(len(samples)) / float64(SampleRate)
	fmt.Printf("%s — %d events over %d bars, %.1fs, peak %.2f\n",
		path, len(haps), bars, seconds, peak)
	if peak == 0 {
		fmt.Fprintln(os.Stderr, "warning: render is silent — check that the pattern sets note or s")
	}
	return nil
}

// Must renders and exits non-zero on failure, for use in template main().
func Must(p core.Pattern, path string, bars int) {
	if err := Render(p, path, bars); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
