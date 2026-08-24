package core

import "testing"

func TestMJS_Port130_ClockAndSession(t *testing.T) {
	c := NewClock(0.5)
	if c.CPS == 0 {
		t.Fatalf("CPS 0")
	}
	// Test SetCPS if available
	c.SetCPS(1.0)
	if c.CPS != 1.0 {
		t.Fatalf("SetCPS 1.0 got %v", c.CPS)
	}
	// Duration/Interval check
	if c.Duration <= 0 {
		t.Fatalf("Duration <=0")
	}
}

func TestMJS_Port130_SignalPerlinWithSlow(t *testing.T) {
	s := Sine().Slow(FractionFromInt(2)).Range(0, 1)
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Sine Slow 2 Range")
	}
	p := Perlin().Slow(FractionFromInt(3))
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Perlin Slow 3")
	}
	r := Rand().Slow(FractionFromInt(2)).Range(0, 100)
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Rand Slow 2 Range 0-100")
	}
}

func TestMJS_Port130_ArpWithFastSlow(t *testing.T) {
	p := Pure("c3 e3 g3").Arp("up").FastF(FractionFromInt(2))
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp up Fast 2")
	}
	q := Pure("c3 e3 g3").Arp("down").Slow(FractionFromInt(2))
	if len(q.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Arp down Slow 2")
	}
	r := Pure("c2 c3 c4").Arp("updown").Slow(FractionFromInt(1))
	if r.QueryArc(FractionFromInt(0), FractionFromInt(1)) == nil {
		t.Fatalf("Arp updown Slow 1")
	}
}
