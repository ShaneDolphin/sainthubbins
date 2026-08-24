// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Minimal dubstep — 140 BPM, half-time.
//
//	go run ./examples/minimal-dubstep
//
// Dubstep runs at 140 BPM but *feels* like 70: the kick lands on beat one and
// the snare waits all the way until beat three. The space between them is the
// instrument. Resist the urge to fill it.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 140

	// Half-time skeleton: one kick, one snare, three quarters of a bar of air.
	kick := core.S(mini.Mini("bd ~ ~ ~"))
	snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))

	// A single hat before the snare — the only thing marking time.
	hat := core.S(mini.Mini("~ ~ ~ hh")).Set(core.Gain(0.25))

	// The sub is one long note per bar. <> steps through its contents one per
	// cycle, so each note is held for a whole bar and the pair makes a two-bar
	// phrase. Drop these an octave to c1/eb1 for a deeper sub than laptop
	// speakers will reproduce.
	sub := core.Note(mini.Mini("<c2 eb2>")).
		Set(core.Cutoff(200)).
		Set(core.Gain(1.0))

	// One sparse detail, and only on the last bar of each four-bar phrase.
	// LastOf(4, ...) starts from silence and swaps in the chord when the cycle
	// is the fourth of its group, so three bars of nothing pass before it
	// arrives.
	chord := core.Note(mini.Mini("~ ~ ~ [g3,bb3]")).
		Set(core.Cutoff(1600)).
		Set(core.Gain(0.2))
	detail := core.Silence().
		LastOf(4, func(core.Pattern) core.Pattern { return chord })

	song := core.Stack(kick, snare, hat, sub, detail)

	shared.Must(song.FastF(shared.Tempo(bpm)), "minimal-dubstep.wav", 8)
}
