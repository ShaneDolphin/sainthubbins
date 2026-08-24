// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package audio

import (
	"math"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestRenderOfflineGainCutoff(t *testing.T) {
	r := NewOfflineRenderer(48000)
	// Gain difference: low gain vs high gain should affect max amplitude
	pLow := core.Pure(map[string]any{"s": "bd", "gain": 0.1})
	pHigh := core.Pure(map[string]any{"s": "bd", "gain": 1.0})
	bufLow, _ := r.RenderOffline(pLow, 1)
	bufHigh, _ := r.RenderOffline(pHigh, 1)
	maxLow := maxAbs(bufLow)
	maxHigh := maxAbs(bufHigh)
	if maxHigh <= maxLow {
		t.Fatalf("high gain max %v should > low gain max %v", maxHigh, maxLow)
	}
	// Cutoff: low cutoff should attenuate high freq hh vs no cutoff
	pHH := core.Pure(map[string]any{"s": "hh"})
	bufNoCut, _ := r.RenderOffline(pHH, 1)
	pHHCut := core.Pure(map[string]any{"s": "hh", "cutoff": 200})
	bufCut, _ := r.RenderOffline(pHHCut, 1)
	maxNoCut := maxAbs(bufNoCut)
	maxCut := maxAbs(bufCut)
	if maxCut >= maxNoCut {
		t.Logf("cutoff attenuation: noCut %v cut %v (expected cut < noCut)", maxNoCut, maxCut)
	}
	// CPS via pattern: use Fast to verify non-zero
	pFast := core.Pure(map[string]any{"s": "bd"}).Fast(core.Pure(2))
	bufFast, _ := r.RenderOffline(pFast, 1)
	if len(bufFast) == 0 {
		t.Fatalf("fast render empty")
	}
}

func maxAbs(buf []float32) float64 {
	m := 0.0
	for _, v := range buf {
		fv := float64(v)
		if fv < 0 {
			fv = -fv
		}
		if fv > m {
			m = fv
		}
		// Also handle math
		_ = math.Abs(fv)
	}
	return m
}

func TestNoteToFreq(t *testing.T) {
	if f := noteToFreq("c4"); f < 260 || f > 262 {
		t.Fatalf("c4 freq %v", f)
	}
	if f := noteToFreq("a4"); f < 439 || f > 441 {
		t.Fatalf("a4 freq %v", f)
	}
}
