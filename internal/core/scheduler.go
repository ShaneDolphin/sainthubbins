// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/zyklus.mjs, cyclist.mjs, schedulerState.mjs — scheduler/clock.

package core

import (
	"context"
	"sync"
	"time"
)

// Clock mirrors zyklus createClock — ticks at interval, advances phase.
type Clock struct {
	mu       sync.Mutex
	CPS      float64
	Duration float64 // seconds per tick (like JS duration)
	Interval float64 // seconds between callbacks
	Overlap  float64
	Phase    float64
	Tick     int
	ticker   *time.Ticker
	stopCh   chan struct{}
}

// NewClock creates a clock with CPS.
func NewClock(cps float64) *Clock {
	return &Clock{
		CPS:      cps,
		Duration: 0.05,
		Interval: 0.1,
		Overlap:  0.1,
	}
}

// SetCPS updates CPS.
func (c *Clock) SetCPS(cps float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CPS = cps
}

// GetPhase returns current phase.
func (c *Clock) GetPhase() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Phase
}

// Start runs clock ticks until context done.
func (c *Clock) Start(ctx context.Context, callback func(phase, duration float64, tick int, now time.Time)) {
	c.mu.Lock()
	if c.stopCh != nil {
		c.mu.Unlock()
		return
	}
	stop := make(chan struct{})
	c.stopCh = stop
	c.mu.Unlock()

	ticker := time.NewTicker(time.Duration(c.Interval * float64(time.Second)))
	c.ticker = ticker
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-stop:
				return
			case t := <-ticker.C:
				now := t
				c.mu.Lock()
				phase := c.Phase
				if phase == 0 {
					phase = float64(now.UnixNano())/1e9 + 0.01
					c.Phase = phase
				}
				lookahead := float64(now.UnixNano())/1e9 + c.Interval + c.Overlap
				dur := c.Duration
				tick := c.Tick
				c.mu.Unlock()
				for phase < lookahead {
					callback(phase, dur, tick, now)
					phase += dur
					tick++
					c.mu.Lock()
					c.Phase = phase
					c.Tick = tick
					c.mu.Unlock()
				}
			}
		}
	}()
}

// Stop halts clock.
func (c *Clock) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.stopCh != nil {
		close(c.stopCh)
		c.stopCh = nil
	}
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
	c.Tick = 0
	c.Phase = 0
}

func (c *Clock) Pause() { c.Stop() }

// SchedulerState mirrors schedulerState.mjs
type SchedulerState struct {
	CPS                      float64
	NumTicksSinceCPSChange   int
	LastTick                 float64
	LastBegin                float64 // cycles
	LastEnd                  float64
	NumCyclesAtCPSChange     float64
	SecondsAtCPSChange       float64
}

// Cyclist is Go port of JS Cyclist — event scheduler.
type Cyclist struct {
	mu                         sync.Mutex
	Started                    bool
	CPS                        float64
	NumTicksSinceCPSChange     int
	LastTick                   float64
	LastBegin                  float64
	LastEnd                    float64
	SecondsAtCPSChange         float64
	NumCyclesAtCPSChange       float64
	Latency                    float64
	Pattern                    *Pattern
	Clock                      *Clock
	OnTrigger                  func(hap Hap, deadline, duration, cps, targetTime float64)
	OnToggle                   func(bool)
	OnError                    func(error)
	getTime                    func() float64
}

// NewCyclist creates a cyclist.
func NewCyclist(interval float64, onTrigger func(Hap, float64, float64, float64, float64), getTime func() float64) *Cyclist {
	if getTime == nil {
		getTime = func() float64 { return float64(time.Now().UnixNano()) / 1e9 }
	}
	c := &Cyclist{
		CPS:       0.5,
		Latency:   0.1,
		getTime:   getTime,
		OnTrigger: onTrigger,
	}
	c.Clock = NewClock(0.5)
	c.Clock.Interval = interval
	if interval == 0 {
		c.Clock.Interval = 0.1
	}
	return c
}

func (c *Cyclist) SetPattern(pat Pattern) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Pattern = &pat
}

func (c *Cyclist) SetCPS(cps float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.CPS == cps {
		return
	}
	c.CPS = cps
	c.NumTicksSinceCPSChange = 0
	c.Clock.SetCPS(cps)
}

func (c *Cyclist) Now() float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.Started {
		return 0
	}
	secondsSinceLastTick := c.getTime() - c.LastTick - c.Clock.Duration
	return c.LastBegin + secondsSinceLastTick*c.CPS
}

func (c *Cyclist) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.Pattern == nil {
		c.mu.Unlock()
		return ErrNoPattern
	}
	c.Started = true
	if c.OnToggle != nil {
		c.OnToggle(true)
	}
	c.NumTicksSinceCPSChange = 0
	c.NumCyclesAtCPSChange = 0
	c.mu.Unlock()

	c.Clock.Start(ctx, func(phase, duration float64, tick int, t time.Time) {
		c.mu.Lock()
		if c.NumTicksSinceCPSChange == 0 {
			c.NumCyclesAtCPSChange = c.LastEnd
			c.SecondsAtCPSChange = phase
		}
		c.NumTicksSinceCPSChange++
		secondsSinceCPSChange := float64(c.NumTicksSinceCPSChange) * duration
		numCyclesSinceCPSChange := secondsSinceCPSChange * c.CPS

		begin := c.LastEnd
		c.LastBegin = begin
		end := c.NumCyclesAtCPSChange + numCyclesSinceCPSChange
		c.LastEnd = end
		c.LastTick = phase
		pat := c.Pattern
		cps := c.CPS
		latency := c.Latency
		numCyclesAtCPSChange := c.NumCyclesAtCPSChange
		secondsAtCPSChange := c.SecondsAtCPSChange
		c.mu.Unlock()

		if phase < float64(t.UnixNano())/1e9 {
			return
		}
		if pat == nil {
			return
		}
		haps := pat.QueryArc(FractionFromFloat(begin), FractionFromFloat(end), map[string]any{"_cps": cps})
		for _, hap := range haps {
			if !hap.HasOnset() {
				continue
			}
			targetTime := (hap.Whole.Begin.Float64()-numCyclesAtCPSChange)/cps + secondsAtCPSChange + latency
			duration := hap.Duration().Float64() / cps
			deadline := targetTime - phase
			if c.OnTrigger != nil {
				c.OnTrigger(hap, deadline, duration, cps, targetTime)
			}
			if v, ok := hap.Value.(map[string]any); ok {
				if cpsVal, exists := v["cps"]; exists {
					var newCPS float64
					switch x := cpsVal.(type) {
					case float64:
						newCPS = x
					case int:
						newCPS = float64(x)
					case Fraction:
						newCPS = x.Float64()
					}
					if newCPS != 0 && newCPS != cps {
						c.SetCPS(newCPS)
					}
				}
			}
		}
	})
	return nil
}

func (c *Cyclist) Pause() {
	c.mu.Lock()
	c.Started = false
	if c.OnToggle != nil {
		c.OnToggle(false)
	}
	c.mu.Unlock()
	c.Clock.Pause()
}

func (c *Cyclist) Stop() {
	c.mu.Lock()
	c.Started = false
	c.LastEnd = 0
	if c.OnToggle != nil {
		c.OnToggle(false)
	}
	c.mu.Unlock()
	c.Clock.Stop()
}

var ErrNoPattern = errNoPattern("scheduler: no pattern set! call SetPattern first.")

type errNoPattern string

func (e errNoPattern) Error() string { return string(e) }
