// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package serial

import "testing"

func TestSerialExtra(t *testing.T) {
	c := New("dummy")
	if c.Port != "dummy" {
		t.Fatalf("port %q", c.Port)
	}
	if n, err := c.Write([]byte("test")); err != nil || n != 4 {
		t.Fatalf("write %v %d", err, n)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close %v", err)
	}
	c2 := New("")
	if c2 == nil {
		t.Fatalf("nil")
	}
	_ = c2.Close()
}
