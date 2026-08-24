// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package audio

import (
	"os"
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestOfflineRenderWAV(t *testing.T) {
	r := NewOfflineRenderer(48000)
	p := core.FastCat(core.Pure(map[string]any{"s": "bd"}), core.Pure(map[string]any{"s": "sd"}))
	buf, err := r.RenderOffline(p, 2)
	if err != nil {
		t.Fatalf("render %v", err)
	}
	// 2 sec per cycle at 0.5cps => 2 cycles = 4 sec => 192000 at 48k
	expected := 2 * 48000 * 2
	if len(buf) != expected {
		t.Fatalf("buf len %d expected %d", len(buf), expected)
	}
	tmp := "/tmp/test_audio_go_phase3.wav"
	if err := WriteWAV(tmp, buf, 48000); err != nil {
		t.Fatalf("write %v", err)
	}
	info, err := os.Stat(tmp)
	if err != nil {
		t.Fatalf("stat %v", err)
	}
	if info.Size() < 1000 {
		t.Fatalf("wav too small %d", info.Size())
	}
	// Verify header starts with RIFF
	f, _ := os.Open(tmp)
	hdr := make([]byte, 4)
	_, _ = f.Read(hdr)
	_ = f.Close()
	if string(hdr) != "RIFF" {
		t.Fatalf("header %q", string(hdr))
	}
}
