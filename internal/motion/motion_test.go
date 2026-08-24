package motion

import "testing"

func TestMotionNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatalf("motion nil")
	}
}
