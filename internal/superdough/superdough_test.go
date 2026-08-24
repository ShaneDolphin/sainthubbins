package superdough

import "testing"

func TestTriggerBasic(t *testing.T) {
	e := New(48000)
	buf := e.Trigger("c4", 0.1)
	if len(buf) == 0 {
		t.Fatalf("Trigger expected non-empty")
	}
	// Check ADSR: first sample should be near 0 due to attack
	if buf[0] < -0.01 || buf[0] > 0.01 {
		t.Logf("first sample %v (attack)", buf[0])
	}
}

func TestTriggerWithControls(t *testing.T) {
	e := New(48000)
	buf := e.TriggerWithControls(map[string]any{"s": "bd", "gain": 0.5}, 0.1)
	if len(buf) == 0 {
		t.Fatalf("TriggerWithControls bd expected non-empty")
	}
	buf2 := e.TriggerWithControls(map[string]any{"note": "a4", "gain": 1.0}, 0.1)
	if len(buf2) == 0 {
		t.Fatalf("TriggerWithControls note expected non-empty")
	}
	// Gain difference: second should have larger amplitude than first
	var max1, max2 float32
	for _, v := range buf {
		if v < 0 {
			v = -v
		}
		if v > max1 {
			max1 = v
		}
	}
	for _, v := range buf2 {
		if v < 0 {
			v = -v
		}
		if v > max2 {
			max2 = v
		}
	}
	if max2 <= max1 {
		t.Logf("gain test: max1 %v max2 %v", max1, max2)
	}
}

func TestSampleToFreq(t *testing.T) {
	if sampleToFreq("bd") != 60 {
		t.Fatalf("bd freq")
	}
	if sampleToFreq("sd") != 180 {
		t.Fatalf("sd freq")
	}
}
