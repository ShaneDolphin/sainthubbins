// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Drum and bass — 174 BPM.
//
//	go run ./examples/drum-and-bass
//
// Drum and bass is fast drums over a slow bassline. The kit runs at 174 BPM
// while the sub moves at half that or less, and the tension between those two
// speeds is the genre.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 174

	// The two-step break: kick on one, snare on three. Simple on paper — the
	// speed is what makes it drum and bass.
	kick := core.S(mini.Mini("bd ~ ~ ~"))
	snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))

	// Ghost notes are the difference between a drum machine and a break:
	// quiet extra hits filling the spaces around the main snare.
	ghosts := core.S(mini.Mini("~ [~ bd] [~ sd] [bd ~]")).Set(core.Gain(0.4))

	// Rolling sixteenth hats carry the tempo.
	hats := core.S(mini.Mini("hh*16")).
		DegradeBy(0.2).
		Set(core.Gain(0.2))

	// The sub moves far slower than the drums: <> holds one note for a whole
	// bar, so the bassline takes two bars to say what the kit says in one.
	sub := core.Note(mini.Mini("<c2 g1>")).
		Set(core.Cutoff(220)).
		Set(core.Gain(1.0))

	// A stab on a Euclidean rhythm, offset so it never lines up with the snare.
	stab := core.Note(mini.Mini("[c3,eb3,g3]")).
		Euclid(3, 8).
		Late(0.125).
		Set(core.Cutoff(1800)).
		Set(core.Gain(0.22))

	// Every fourth bar, reverse the ghost notes for a turnaround.
	ghosts = ghosts.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })

	song := core.Stack(kick, snare, ghosts, hats, sub, stab)

	shared.Must(song.FastF(shared.Tempo(bpm)), "drum-and-bass.wav", 8)
}
