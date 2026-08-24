package mqtt

import "testing"

func TestMQTTNew(t *testing.T) {
	c := New("tcp://localhost:1883")
	if c.Broker != "tcp://localhost:1883" {
		t.Fatalf("mqtt broker")
	}
	if err := c.Publish("test", []byte("hi")); err != nil {
		t.Fatalf("publish")
	}
}
