package audio

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestRenderOffline(t *testing.T) {
	pat := core.Pure(map[string]any{"s": "bd"})
	samples, err := RenderPatternAudio(pat, 1, 48000)
	if err != nil {
		t.Fatalf("render err %v", err)
	}
	if len(samples) == 0 {
		t.Fatalf("no samples")
	}
	// Check not all zero
	nonZero := 0
	for _, s := range samples {
		if s != 0 {
			nonZero++
			break
		}
	}
	if nonZero == 0 {
		t.Fatalf("all zero")
	}
}

func TestWriteWAV(t *testing.T) {
	pat := core.Pure(map[string]any{"freq": 440.0})
	samples, _ := RenderPatternAudio(pat, 1, 8000)
	err := WriteWAV("/tmp/test_audio_go.wav", samples, 8000)
	if err != nil {
		t.Fatalf("write wav %v", err)
	}
}
