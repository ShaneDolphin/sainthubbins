// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// mytrack — the finished result of docs/tutorial/06-new-song-cli.md.
//
//	go run ./examples/mytrack
//
// Chapter 6 builds this file one layer at a time. It is here so you can run the
// finished version, or compare it against your own as you go.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 128

	// Step 1 — the pulse.
	kick := core.S(mini.Mini("bd*4"))

	// Step 2 — the backbeat, answering on 2 and 4.
	clap := core.S(mini.Mini("~ cp ~ cp")).
		Set(core.Gain(0.6))

	// Step 3 — off-beat hats, thinned out so they breathe.
	hats := core.S(mini.Mini("[~ hh]*4")).
		DegradeBy(0.15).
		Set(core.Gain(0.3))

	// Step 4 — a bassline in C minor that leaves beat 2 open.
	bass := core.Note(mini.Mini("c2 ~ c2 eb2")).
		Set(core.Cutoff(600)).
		Set(core.Gain(0.7))

	// Step 5 — a chord on the off-beat. Three notes at once, so keep it quiet.
	chord := core.Note(mini.Mini("~ [c3,eb3,g3] ~ ~")).
		Set(core.Cutoff(1800)).
		Set(core.Gain(0.18))

	// Step 6 — stop it looping. Every fourth bar the bass runs backwards.
	bass = bass.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })

	song := core.Stack(kick, clap, hats, bass, chord)

	shared.Must(song.FastF(shared.Tempo(bpm)), "mytrack.wav", 8)
}
