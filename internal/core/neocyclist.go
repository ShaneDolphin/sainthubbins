// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

// NeoCyclist is multi-instance scheduler stub
type NeoCyclist struct {
	Cyclists []*Cyclist
}

func NewNeoCyclist() *NeoCyclist { return &NeoCyclist{} }

func (n *NeoCyclist) Add(c *Cyclist) { n.Cyclists = append(n.Cyclists, c) }
