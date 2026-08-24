package audio

import "testing"

func TestOfflineRendererSampleRate(t *testing.T) {
	r := NewOfflineRenderer(44100)
	if r.SampleRate() != 44100 {
		t.Fatalf("SampleRate expected 44100")
	}
	if r.SampleRateValue != 44100 {
		t.Fatalf("SampleRateValue")
	}
}

func TestOfflineRendererDefault(t *testing.T) {
	r := NewOfflineRenderer(0)
	if r.SampleRate() != 48000 {
		t.Fatalf("default 48000")
	}
}
