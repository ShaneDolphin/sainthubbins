// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Techno — 132 BPM.
//
//	go run ./examples/techno
//
// Techno trades the warmth of house for relentlessness: the same four-to-the-
// floor pulse, but darker, tighter, and with almost nothing on top.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 132

	kick := core.S(mini.Mini("bd*4"))

	// Sixteenth hats sit under everything as a constant hiss, with the
	// off-beat hats on top of them for the push.
	hatsClosed := core.S(mini.Mini("hh*16")).Set(core.Gain(0.18))
	hatsOpen := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.4))

	// A single dark stab placed by a Euclidean rhythm: three hits spread as
	// evenly as possible across eight steps, which lands them off the grid.
	stab := core.Note(mini.Mini("c3")).
		Euclid(3, 8).
		Set(core.Cutoff(900)).
		Set(core.Gain(0.4))

	// Rumbling sub under the kick, filtered almost shut.
	sub := core.Note(mini.Mini("c2 ~ c2 ~")).
		Set(core.Cutoff(280)).
		Set(core.Gain(0.8))

	// Clap only every other bar, so the loop has a two-bar shape.
	clap := core.S(mini.Mini("~ ~ cp ~")).
		Every(2, func(p core.Pattern) core.Pattern { return core.Silence() }).
		Set(core.Gain(0.5))

	song := core.Stack(kick, hatsClosed, hatsOpen, stab, sub, clap)

	shared.Must(song.FastF(shared.Tempo(bpm)), "techno.wav", 8)
}
