// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Maximal dubstep — 140 BPM.
//
//	go run ./examples/maximal-dubstep
//
// Same 140 BPM half-time skeleton as the minimal template, and the same kick
// and snare placement — then everything the minimal version leaves out is
// crammed into the gaps. Compare the two files side by side: the difference
// between them is the entire lesson.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 140

	// The skeleton is unchanged — kick on one, snare on three.
	kick := core.S(mini.Mini("bd ~ ~ ~"))
	snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))

	// ...but now there are ghost kicks crowding the second half of the bar.
	ghostKick := core.S(mini.Mini("~ ~ ~ [bd bd]")).Set(core.Gain(0.55))

	// Busy sixteenth hats, with a third dropped so they stay ragged.
	hats := core.S(mini.Mini("hh*16")).
		DegradeBy(0.3).
		Set(core.Gain(0.22))

	// The wobble. Eight sub notes per bar, each with its own cutoff value —
	// that sweeping filter is what makes a dubstep bass talk.
	wobble := core.Note(mini.Mini("c1*8")).
		Set(core.Cutoff(mini.Mini("180 900 260 1500 200 1100 320 1800"))).
		Set(core.Gain(1.0))

	// Every second bar the wobble doubles in speed — the classic half/double
	// switch. Ply(2) repeats each event twice inside its own time slot.
	wobble = wobble.Every(2, func(p core.Pattern) core.Pattern { return p.Ply(2) })

	// A shrieking lead answering the snare.
	lead := core.Note(mini.Mini("~ ~ [g4,bb4] ~")).
		Set(core.Cutoff(3000)).
		Set(core.Gain(0.3))

	// A snare fill on the last bar of each four-bar phrase. LastOf(4, ...)
	// starts from silence and swaps in the fill on the fourth cycle of every
	// group — which is where a fill belongs.
	fill := core.Silence().
		LastOf(4, func(core.Pattern) core.Pattern {
			return core.S(mini.Mini("~ ~ ~ [sd sd sd sd]")).Set(core.Gain(0.6))
		})

	song := core.Stack(kick, snare, ghostKick, hats, wobble, lead, fill)

	shared.Must(song.FastF(shared.Tempo(bpm)), "maximal-dubstep.wav", 8)
}
