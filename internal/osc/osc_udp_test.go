// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// listener binds a loopback UDP socket and returns it with its port. Tests
// stay hermetic: nothing leaves the machine and no SuperDirt is required.
func listener(t *testing.T) (*net.UDPConn, int) {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn, conn.LocalAddr().(*net.UDPAddr).Port
}

func TestClientSendReachesTheWire(t *testing.T) {
	conn, port := listener(t)
	c := New("127.0.0.1", port)
	defer c.Close()

	if err := c.Send("/dirt/play", "s", "bd"); err != nil {
		t.Fatalf("Send: %v", err)
	}

	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("nothing received: %v", err)
	}
	if !bytes.HasPrefix(buf[:n], []byte("/dirt/play\x00\x00")) {
		t.Errorf("received %q, want a /dirt/play message", buf[:n])
	}
	if !bytes.Contains(buf[:n], []byte("bd\x00")) {
		t.Errorf("received %q, want it to carry the sound name", buf[:n])
	}
}

func TestClientSendAtSendsABundle(t *testing.T) {
	conn, port := listener(t)
	c := New("127.0.0.1", port)
	defer c.Close()

	if err := c.SendAt(time.Now().Add(time.Second), "/dirt/play", "s", "sd"); err != nil {
		t.Fatalf("SendAt: %v", err)
	}
	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("nothing received: %v", err)
	}
	if !bytes.HasPrefix(buf[:n], []byte("#bundle\x00")) {
		t.Errorf("received %q, want a bundle", buf[:n])
	}
}

// A client with no host is a sink, so tests and offline use never need a peer.
func TestClientWithoutHostIsANoOp(t *testing.T) {
	c := New("", 0)
	if err := c.Send("/dirt/play", "s", "bd"); err != nil {
		t.Errorf("Send on a hostless client should be a no-op, got %v", err)
	}
	if err := c.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}

// Connect must dial eagerly and report failure for a host it cannot reach,
// so a caller (runPlay) can fail fast instead of finding out on the first
// event. A host string containing a null byte fails inside net.Dial itself
// ("invalid argument") without performing any DNS lookup or other network
// access, which keeps this deterministic — the same reasoning applied to
// TestClientRecoversFromDialFailure above.
func TestClientConnectFailsForUnreachableHost(t *testing.T) {
	c := New("127.0.0.1\x00evil", 57120)
	defer c.Close()

	if err := c.Connect(); err == nil {
		t.Fatal("Connect: want an error for a host net.Dial cannot use, got nil")
	}
}

// Connect must succeed for a reachable loopback peer, and a hostless client
// must remain a successful no-op sink.
func TestClientConnectSucceedsForLoopback(t *testing.T) {
	_, port := listener(t)
	c := New("127.0.0.1", port)
	defer c.Close()

	if err := c.Connect(); err != nil {
		t.Fatalf("Connect: %v", err)
	}
}

func TestClientConnectHostlessIsANoOp(t *testing.T) {
	c := New("", 0)
	defer c.Close()

	if err := c.Connect(); err != nil {
		t.Errorf("Connect on a hostless client should be a no-op, got %v", err)
	}
}

// A dial failure must not brick the client permanently: once the underlying
// condition clears, a later Send should succeed rather than replaying a
// cached error forever. The first dial targets an out-of-range port, which
// net.Dial rejects immediately and locally ("address 999999: invalid port")
// without any network access — unlike a DNS lookup for a reserved hostname
// such as "invalid.invalid", which a captive-portal or corporate resolver
// that hijacks NXDOMAIN can resolve anyway, making that approach flaky.
func TestClientRecoversFromDialFailure(t *testing.T) {
	c := New("127.0.0.1", 999999)
	defer c.Close()

	if err := c.Send("/dirt/play", "s", "bd"); err == nil {
		t.Fatalf("Send: expected the first dial (to the out-of-range port) to fail")
	}

	conn, port := listener(t)
	c.Host, c.Port = "127.0.0.1", port

	if err := c.Send("/dirt/play", "s", "bd"); err != nil {
		t.Fatalf("Send after recovery: %v", err)
	}

	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("nothing received: %v", err)
	}
	if !bytes.HasPrefix(buf[:n], []byte("/dirt/play\x00\x00")) {
		t.Errorf("received %q, want a /dirt/play message", buf[:n])
	}
}
