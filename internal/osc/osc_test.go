package osc

import "testing"

func TestOSCNew(t *testing.T) {
	c := New("localhost", 57120)
	if c.Host != "localhost" || c.Port != 57120 {
		t.Fatalf("osc")
	}
	if err := c.SendSuperDirt(nil); err != nil {
		t.Fatalf("send")
	}
}
