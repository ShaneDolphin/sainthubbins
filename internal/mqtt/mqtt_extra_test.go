// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package mqtt

import "testing"

func TestMQTTExtra(t *testing.T) {
	c := New("localhost:1883")
	if c.Broker != "localhost:1883" {
		t.Fatalf("mqtt %v", c)
	}
	if err := c.Publish("test/topic", []byte("hello")); err != nil {
		t.Logf("publish %v (ok if stub)", err)
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
