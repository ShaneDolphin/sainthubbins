package gamepad

import "testing"

func TestGamepadNew(t *testing.T) {
	g := New()
	if g == nil {
		t.Fatalf("gamepad nil")
	}
}
