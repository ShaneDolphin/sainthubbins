// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Chicago House — 120 BPM.
//
//	go run ./examples/chicago-house
//
// Where house is polished, Chicago house is a drum machine and an attitude:
// a hard 808 kick, sharp claps, rattling sixteenth hats, and a squelching
// bassline that never sits still.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 120

	kick := core.S(mini.Mini("bd*4"))

	// Clap on the backbeat, louder and drier than the house template.
	clap := core.S(mini.Mini("~ cp ~ cp")).
		Set(core.Gain(0.7))

	// Sixteenth hats, then DegradeBy drops a third of them at random so the
	// machine breathes instead of marching.
	hats := core.S(mini.Mini("hh*16")).
		DegradeBy(0.3).
		Set(core.Gain(0.3))

	// The acid line: eight sixteenth notes with a cutoff that moves on every
	// step. The moving filter is the whole character of the sound.
	acid := core.Note(mini.Mini("c2 c2 eb2 c2 g2 c2 bb1 c2")).
		Set(core.Cutoff(mini.Mini("300 1400 500 2000 700 1800 400 1100"))).
		Set(core.Gain(0.7))

	// Every fourth bar, run the acid line backwards for a turnaround.
	acid = acid.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })

	song := core.Stack(kick, clap, hats, acid)

	shared.Must(song.FastF(shared.Tempo(bpm)), "chicago-house.wav", 8)
}
