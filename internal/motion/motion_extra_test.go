// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package motion

import "testing"

func TestMotionExtra(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatalf("nil")
	}
	m2 := New()
	if m2 == nil {
		t.Fatalf("nil2")
	}
}
