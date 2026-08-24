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

	mu      sync.Mutex
	conn    net.Conn
	dialErr error
	dialed  bool
}

// New creates a client. It does not dial — UDP has no handshake worth doing
// eagerly, and construction should not fail.
func New(host string, port int) *Client { return &Client{Host: host, Port: port} }

// ensure dials once, on first use.
func (c *Client) ensure() (net.Conn, error) {
	if c.Host == "" {
		return nil, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dialed {
		return c.conn, c.dialErr
	}
	c.dialed = true
	c.conn, c.dialErr = net.Dial("udp", net.JoinHostPort(c.Host, strconv.Itoa(c.Port)))
	return c.conn, c.dialErr
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

// SendSuperDirt is retained for compatibility with existing callers.
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
