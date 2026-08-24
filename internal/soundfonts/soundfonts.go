// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package soundfonts

// Font represents a loaded SoundFont
type Font struct {
	Name string
	URL  string
}

// List is subset of GM soundfonts
var List = []Font{
	{Name: "piano", URL: "https://example.com/piano.sf2"},
	{Name: "strings", URL: "https://example.com/strings.sf2"},
}
