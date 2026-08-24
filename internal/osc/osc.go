// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package osc

import "fmt"

// Client is OSC SuperDirt stub
type Client struct {
	Host string
	Port int
}

func New(host string, port int) *Client { return &Client{Host: host, Port: port} }

func (c *Client) SendSuperDirt(haps []interface{}) error {
	_ = fmt.Sprintf("%v", haps)
	return nil
}
func (c *Client) Close() error { return nil }
