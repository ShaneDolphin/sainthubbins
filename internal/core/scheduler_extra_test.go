// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestSchedulerCPS(t *testing.T) {
	c := NewClock(0.5)
	if c.CPS != 0.5 {
		t.Fatalf("cps 0.5 got %v", c.CPS)
	}
	c.SetCPS(1.0)
	if c.CPS != 1.0 {
		t.Fatalf("cps 1.0 got %v", c.CPS)
	}
}

func TestCyclistSetPattern(t *testing.T) {
	cy := NewCyclist(0.1, nil, nil)
	if cy == nil {
		t.Fatalf("cyclist nil")
	}
	p := Pure("bd")
	cy.SetPattern(p)
	cy.SetCPS(2.0)
	if cy.CPS != 2.0 {
		t.Fatalf("cyclist cps 2 got %v", cy.CPS)
	}
	if cy.Pattern == nil {
		t.Fatalf("pattern nil after SetPattern")
	}
}

func TestCyclistNow(t *testing.T) {
	now := 1000.0
	cy := NewCyclist(0.1, func(h Hap, deadline, duration, cps, targetTime float64) {}, func() float64 { return now })
	cy.SetPattern(Sequence(Pure("bd"), Pure("sd")))
	cy.SetCPS(1.0)
	if cy.Now() != 0 {
		t.Fatalf("Now before start should be 0 got %v", cy.Now())
	}
}
