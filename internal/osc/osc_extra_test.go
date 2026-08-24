// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package osc

import "testing"

func TestOSCMessageBasic(t *testing.T) {
	c := New("127.0.0.1", 57120)
	if c.Host != "127.0.0.1" || c.Port != 57120 {
		t.Fatalf("client %v", c)
	}
	if err := c.SendSuperDirt([]interface{}{"bd"}); err != nil {
		t.Fatalf("send %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("close %v", err)
	}
}

func TestOSCBundle(t *testing.T) {
	c := New("", 0)
	if c == nil {
		t.Fatalf("nil")
	}
	_ = c.SendSuperDirt(nil)
	_ = c.Close()
}
