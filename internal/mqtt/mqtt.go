// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package mqtt

// Client is paho-mqtt stub
type Client struct {
	Broker string
}

func New(broker string) *Client { return &Client{Broker: broker} }

func (c *Client) Publish(topic string, payload []byte) error { return nil }
func (c *Client) Close() error { return nil }
