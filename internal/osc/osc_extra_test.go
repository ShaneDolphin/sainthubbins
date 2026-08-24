// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package osc

import "testing"

func TestOSCMessageBasic(t *testing.T) {
	// Was New("127.0.0.1", 57120): that dials the default SuperDirt port, so
	// SendSuperDirt below sent a real /dirt/play packet to whatever is
	// listening there on every `go test ./...` run — including on a machine
	// that has a live SuperDirt/synth for verifying this feature. Point it
	// at a loopback listener bound to an OS-assigned port instead, the way
	// osc_udp_test.go's tests do, so the test stays hermetic.
	_, port := listener(t)
	c := New("127.0.0.1", port)
	if c.Host != "127.0.0.1" || c.Port != port {
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
