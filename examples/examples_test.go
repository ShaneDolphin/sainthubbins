// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Guards the tutorial templates: every one must still build, run, and produce
// audible, non-clipping audio. Without this a change to the pattern engine can
// break the documented examples silently.

package examples_test

import (
	"encoding/binary"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// templates lists every tutorial template and the WAV file it writes.
var templates = []struct {
	dir string
	wav string
}{
	{"house", "house.wav"},
	{"chicago-house", "chicago-house.wav"},
	{"techno", "techno.wav"},
	{"minimal-dubstep", "minimal-dubstep.wav"},
	{"maximal-dubstep", "maximal-dubstep.wav"},
	{"drum-and-bass", "drum-and-bass.wav"},
	{"electronica", "electronica.wav"},
	{"trance", "trance.wav"},
	// The worked example built in tutorial chapter 6.
	{"mytrack", "mytrack.wav"},
}

func TestTemplatesRenderAudio(t *testing.T) {
	if testing.Short() {
		t.Skip("builds and runs nine binaries; skipped under -short")
	}
	for _, tc := range templates {
		tc := tc
		t.Run(tc.dir, func(t *testing.T) {
			t.Parallel()
			tmp := t.TempDir()
			bin := filepath.Join(tmp, "tmpl")

			// The test's working directory is examples/, so ./<dir> is the
			// template package.
			build := exec.Command("go", "build", "-o", bin, "./"+tc.dir)
			if out, err := build.CombinedOutput(); err != nil {
				t.Fatalf("build %s: %v\n%s", tc.dir, err, out)
			}

			run := exec.Command(bin)
			run.Dir = tmp // the template writes its WAV into the working directory
			out, err := run.CombinedOutput()
			if err != nil {
				t.Fatalf("run %s: %v\n%s", tc.dir, err, out)
			}

			peak, frames := readWAV(t, filepath.Join(tmp, tc.wav))
			if peak == 0 {
				t.Errorf("%s rendered silence", tc.dir)
			}
			// Layers sum in the renderer, so a template that stacks too loudly
			// clips. Full scale is 32767.
			if peak >= 32767 {
				t.Errorf("%s clips (peak %d); reduce a Gain", tc.dir, peak)
			}
			// Every template renders 8 bars at 2s per bar, 48kHz.
			if want := 8 * 2 * 48000; frames != want {
				t.Errorf("%s: got %d frames, want %d", tc.dir, frames, want)
			}
		})
	}
}

// readWAV returns the peak sample magnitude and frame count of a 16-bit mono WAV.
func readWAV(t *testing.T, path string) (peak int, frames int) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	const headerLen = 44
	if len(b) <= headerLen {
		t.Fatalf("%s is too short to contain audio (%d bytes)", path, len(b))
	}
	data := b[headerLen:]
	frames = len(data) / 2
	for i := 0; i+1 < len(data); i += 2 {
		v := int(int16(binary.LittleEndian.Uint16(data[i : i+2])))
		if v < 0 {
			v = -v
		}
		if v > peak {
			peak = v
		}
	}
	return peak, frames
}

// TestDocumentedIdioms exercises the Go-side idioms the tutorial teaches, so a
// change to the engine cannot leave the prose describing behaviour that no
// longer happens. Each case asserts on a whole-span query, which is what the
// offline renderer performs.
func TestDocumentedIdioms(t *testing.T) {
	span := func(p core.Pattern, cycles int64) int {
		return len(p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(cycles)))
	}
	perCycle := func(p core.Pattern, c int64) int {
		return len(p.QueryArc(core.FractionFromInt(c), core.FractionFromInt(c+1)))
	}

	t.Run("chord is simultaneous", func(t *testing.T) {
		// Chapter 2: a comma-stack is a chord.
		haps := core.Note(mini.Mini("[c3,eb3,g3]")).
			QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(haps) != 3 {
			t.Fatalf("chord: got %d notes, want 3", len(haps))
		}
		for _, h := range haps {
			if h.Whole == nil || h.Whole.Begin.Float64() != 0 {
				t.Errorf("chord note %v does not start at 0", h.Value)
			}
		}
	})

	t.Run("Set merges controls", func(t *testing.T) {
		// Chapter 4: Set merges a control into every event.
		haps := core.S(mini.Mini("hh")).Set(core.Gain(0.4)).
			QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		m, ok := haps[0].Value.(map[string]any)
		if !ok {
			t.Fatalf("want a control bag, got %T", haps[0].Value)
		}
		if m["s"] != "hh" || m["gain"] != 0.4 {
			t.Errorf("got %v, want s:hh and gain:0.4", m)
		}
	})

	t.Run("patterned control varies per step", func(t *testing.T) {
		// Chapter 4: the wobble technique.
		haps := core.Note(mini.Mini("c1*4")).
			Set(core.Cutoff(mini.Mini("200 900 400 1600"))).
			QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		seen := map[any]bool{}
		for _, h := range haps {
			seen[h.Value.(map[string]any)["cutoff"]] = true
		}
		if len(seen) != 4 {
			t.Errorf("want 4 distinct cutoff values, got %d: %v", len(seen), seen)
		}
	})

	t.Run("Every fires on matching cycles only", func(t *testing.T) {
		// Chapter 5. This is the case that silently broke the templates.
		p := core.Silence().Every(4, func(core.Pattern) core.Pattern {
			return core.S(mini.Mini("sd*4"))
		})
		for c := int64(0); c < 8; c++ {
			want := 0
			if c%4 == 0 {
				want = 4
			}
			if got := perCycle(p, c); got != want {
				t.Errorf("Every: cycle %d got %d, want %d", c, got, want)
			}
		}
		if got, want := span(p, 8), 8; got != want {
			t.Errorf("Every over one 8-cycle query: got %d, want %d", got, want)
		}
	})

	t.Run("LastOf fires at the end of each phrase", func(t *testing.T) {
		// Chapter 5, and both dubstep fills.
		p := core.Silence().LastOf(4, func(core.Pattern) core.Pattern {
			return core.S(mini.Mini("~ ~ ~ [sd sd sd sd]"))
		})
		for c := int64(0); c < 8; c++ {
			want := 0
			if c%4 == 3 {
				want = 4
			}
			if got := perCycle(p, c); got != want {
				t.Errorf("LastOf: cycle %d got %d, want %d", c, got, want)
			}
		}
	})

	t.Run("Off adds a delayed copy", func(t *testing.T) {
		// Chapter 5: the echo idiom.
		base := core.S(mini.Mini("bd sd"))
		echoed := base.Off(0.125, func(p core.Pattern) core.Pattern {
			return p.Set(core.Gain(0.3))
		})
		// The shifted copy wraps at the cycle boundary, so the count is not
		// exactly double. What matters is that the original survives and a
		// quieter, later copy exists alongside it.
		haps := echoed.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(haps) <= span(base, 1) {
			t.Errorf("Off: got %d haps, want more than the original %d", len(haps), span(base, 1))
		}
		var dry, wet int
		for _, h := range haps {
			if m, ok := h.Value.(map[string]any); ok {
				if _, hasGain := m["gain"]; hasGain {
					wet++
				} else {
					dry++
				}
			}
		}
		if dry == 0 || wet == 0 {
			t.Errorf("Off: want both original (%d) and echoed (%d) events", dry, wet)
		}
	})

	t.Run("Struct needs Go booleans", func(t *testing.T) {
		// Chapter 8 documents this exact trap.
		mask := core.FastCat(core.Pure(true), core.Pure(false), core.Pure(true), core.Pure(true))
		if got := span(core.Note(mini.Mini("c3")).Struct(mask), 1); got != 3 {
			t.Errorf("Struct with booleans: got %d, want 3", got)
		}
		if got := span(core.Note(mini.Mini("c3")).Struct(mini.Mini("t ~ t t")), 1); got != 0 {
			t.Errorf("Struct with mini strings: got %d, want 0 (documented as not working)", got)
		}
	})

	t.Run("transpose before wrapping", func(t *testing.T) {
		// Chapter 5: Add on a bare number still works, unwrapped or wrapped
		// after the fact — both land in the same place.
		haps := core.Note(core.Pure(60).Add(12)).
			QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		m, ok := haps[0].Value.(map[string]any)
		if !ok || m["note"] != 72.0 {
			t.Errorf("got %v, want map[note:72]", haps[0].Value)
		}
	})

	t.Run("transpose an already-wrapped pattern", func(t *testing.T) {
		// Chapter 5: Add now transposes a wrapped pattern directly instead
		// of flattening it to a bare number.
		haps := core.Note(mini.Mini("0 4 7")).Add(12).
			QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
		if len(haps) != 3 {
			t.Fatalf("got %d haps, want 3", len(haps))
		}
		want := []float64{12, 16, 19}
		for i, w := range want {
			m, ok := haps[i].Value.(map[string]any)
			if !ok {
				t.Fatalf("hap %d: value is %T, want a control bag", i, haps[i].Value)
			}
			if m["note"] != w {
				t.Errorf("hap %d note = %v, want %v", i, m["note"], w)
			}
		}
	})

	t.Run("tempo is a ratio against 120", func(t *testing.T) {
		// Chapter 3: BPM = 120 * factor.
		f := shared.Tempo(140)
		if got, want := f.Float64(), 140.0/120.0; got != want {
			t.Errorf("Tempo(140) = %v, want %v", got, want)
		}
	})
}
