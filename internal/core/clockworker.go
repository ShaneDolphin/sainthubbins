// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

// ClockWorker is Web Worker clock stub (now goroutine)
type ClockWorker struct {
	Clock *Clock
}

func NewClockWorker(c *Clock) *ClockWorker { return &ClockWorker{Clock: c} }
