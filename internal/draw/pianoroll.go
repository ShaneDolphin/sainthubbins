// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package draw

import (
	"encoding/json"
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// Event is a drawable event for pianoroll
type Event struct {
	Time     float64       `json:"time"`
	Duration float64       `json:"duration"`
	Value    any           `json:"value"`
	Context  map[string]any `json:"context,omitempty"`
}

// Pianoroll converts haps to drawable events (time in cycles)
func Pianoroll(haps []core.Hap) []Event {
	evs := make([]Event, 0, len(haps))
	for _, h := range haps {
		if h.Whole == nil {
			continue
		}
		evs = append(evs, Event{
			Time:     h.Whole.Begin.Float64(),
			Duration: h.Whole.Duration().Float64(),
			Value:    h.Value,
			Context:  h.Context,
		})
	}
	return evs
}

// ToJSON returns JSON for frontend
func ToJSON(haps []core.Hap) string {
	evs := Pianoroll(haps)
	b, _ := json.Marshal(evs)
	return string(b)
}

// Spiral maps haps to polar coordinates (angle = time mod 1 * 2pi, radius = cycle)
func Spiral(haps []core.Hap) []Event {
	evs := make([]Event, 0, len(haps))
	for _, h := range haps {
		if h.Whole == nil {
			continue
		}
		angle := h.Whole.Begin.CyclePos().Float64() * 6.28318530718
		radius := h.Whole.Begin.Sam().Float64() + 1
		evs = append(evs, Event{
			Time:     angle,
			Duration: radius,
			Value:    h.Value,
			Context:  h.Context,
		})
	}
	return evs
}
