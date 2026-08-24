// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// House — 125 BPM.
//
//	go run ./examples/house
//
// The defining sound of house is the four-on-the-floor kick with the hats
// pushed into the gaps between the beats. Everything else here is decoration.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 125

	// Kick on every beat — the floor everyone dances on.
	kick := core.S(mini.Mini("bd*4"))

	// Clap answers on beats 2 and 4. The rests are what make it an answer.
	clap := core.S(mini.Mini("~ cp ~ cp")).
		Set(core.Gain(0.6))

	// Hats land on the off-beats: [~ hh] is a rest then a hat, repeated four
	// times, so each hat sits exactly halfway between two kicks.
	hats := core.S(mini.Mini("[~ hh]*4")).
		Set(core.Gain(0.3))

	// Bassline in A minor, following the kick but leaving beat 3 open.
	bass := core.Note(mini.Mini("a1 ~ a1 c2")).
		Set(core.Cutoff(700)).
		Set(core.Gain(0.7))

	// Off-beat chord stabs — the other half of the house signature.
	// [a3,c4,e4] is a comma-stack: three notes sounding together.
	stabs := core.Note(mini.Mini("~ [a3,c4,e4] ~ [a3,c4,e4]")).
		Set(core.Cutoff(2200)).
		Set(core.Gain(0.18))

	song := core.Stack(kick, clap, hats, bass, stabs)

	shared.Must(song.FastF(shared.Tempo(bpm)), "house.wav", 8)
}
