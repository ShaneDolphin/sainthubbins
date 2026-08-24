// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package xen

// Tune represents xenharmonic tuning
type Tune struct {
	Name string
	Steps int
}

// GetTune returns tuning by name stub
func GetTune(name string) *Tune {
	return &Tune{Name: name, Steps: 12}
}
