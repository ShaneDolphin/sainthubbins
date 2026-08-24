// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package embed

// Embed is iframe embed stub
type Embed struct {
	URL string
}

func New(url string) *Embed { return &Embed{URL: url} }
