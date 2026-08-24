// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// UDP client for SuperDirt and other OSC receivers.

package osc

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// DirtAddress is the OSC address SuperDirt listens on for note events.
const DirtAddress = "/dirt/play"

// Client sends OSC over UDP. A Client with an empty Host is a sink: every send
// succeeds and goes nowhere, so offline rendering and tests need no peer.
type Client struct {
	Host string
	Port int

	mu     sync.Mutex
	conn   net.Conn
	dialed bool
}

// New creates a client. It does not dial — UDP has no handshake worth doing
// eagerly, and construction should not fail.
func New(host string, port int) *Client { return &Client{Host: host, Port: port} }

// Connect dials the peer eagerly and reports whether it succeeded. Callers
// that want to fail fast — rather than discovering a bad host or port only
// when the first event silently goes nowhere — should call this once before
// starting a scheduler, so the dial (including any DNS resolution) happens
// off the tick goroutine. It reuses ensure()'s lazy-dial machinery, so a
// hostless client remains a successful no-op sink.
func (c *Client) Connect() error {
	_, err := c.ensure()
	return err
}

// ensure dials on first use and again after any dial that failed, so a
// transient failure (a DNS hiccup, SuperDirt not up yet) does not brick the
// client for the rest of the process. Only a successful dial is cached.
func (c *Client) ensure() (net.Conn, error) {
	if c.Host == "" {
		return nil, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dialed {
		return c.conn, nil
	}
	conn, err := net.Dial("udp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
	if err != nil {
		return nil, err
	}
	c.conn, c.dialed = conn, true
	return c.conn, nil
}

func (c *Client) write(b []byte) error {
	conn, err := c.ensure()
	if err != nil {
		return err
	}
	if conn == nil {
		return nil // hostless sink
	}
	_, err = conn.Write(b)
	return err
}

// Send transmits one OSC message immediately.
func (c *Client) Send(addr string, args ...any) error {
	msg, err := EncodeMessage(addr, args...)
	if err != nil {
		return err
	}
	return c.write(msg)
}

// SendAt transmits one OSC message inside a bundle timestamped for at, so the
// receiver plays it at that moment rather than on arrival.
func (c *Client) SendAt(at time.Time, addr string, args ...any) error {
	msg, err := EncodeMessage(addr, args...)
	if err != nil {
		return err
	}
	return c.write(EncodeBundle(at, msg))
}

// SendSuperDirt is a compatibility shim for callers migrating from an older
// API; nothing in this codebase calls it anymore (runPlay sends through
// SendAt with osc.DirtArgs's alternating key/value list). It forwards haps
// as raw positional arguments rather than that key/value form, so most
// calls produce an odd-length /dirt/play payload that real SuperDirt does
// not expect — treat it as compat-only, not a production path.
func (c *Client) SendSuperDirt(haps []interface{}) error {
	if len(haps) == 0 {
		return nil
	}
	return c.Send(DirtAddress, haps...)
}

// Close releases the socket if one was opened.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	c.dialed = false
	return err
}
