// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

func TestRunPlaySendsToSuperDirt(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer conn.Close()
	port := conn.LocalAddr().(*net.UDPAddr).Port

	received := make(chan []byte, 8)
	go func() {
		buf := make([]byte, 2048)
		for {
			_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				close(received)
				return
			}
			cp := make([]byte, n)
			copy(cp, buf[:n])
			received <- cp
		}
	}()

	var out bytes.Buffer
	if err := runPlay("bd*4", "127.0.0.1", port, 1.0, &out); err != nil {
		t.Fatalf("runPlay: %v", err)
	}

	select {
	case msg, ok := <-received:
		if !ok {
			t.Fatal("listener closed without receiving anything")
		}
		if !bytes.Contains(msg, []byte("bd")) {
			t.Errorf("received %q, want it to carry the sound name", msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no OSC arrived within the timeout")
	}

	if !strings.Contains(out.String(), "127.0.0.1") {
		t.Errorf("runPlay should report where it is sending, got %q", out.String())
	}
}
