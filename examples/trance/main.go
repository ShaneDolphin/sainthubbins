// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Trance — 138 BPM.
//
//	go run ./examples/trance
//
// Trance is built on one trick: the kick lands on the beat and the bass answers
// in every gap between kicks. That interlock is the engine. Over it goes a
// sixteenth-note arpeggio that never stops moving.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 138

	kick := core.S(mini.Mini("bd*4"))

	// The off-beat bass. [~ a1] is a rest then a note, four times per bar, so
	// every bass note falls precisely in a gap between two kicks.
	bass := core.Note(mini.Mini("[~ a1]*4")).
		Set(core.Cutoff(600)).
		Set(core.Gain(0.6))

	// Open hats double the off-beat to reinforce it.
	hats := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.3))
	hatsClosed := core.S(mini.Mini("hh*16")).Set(core.Gain(0.1))

	// Sixteenth arpeggio climbing an A minor triad and back down.
	arp := core.Note(mini.Mini("a3 c4 e4 a4 e4 c4 a3 c4")).
		FastF(core.NewFraction(2, 1)).
		Set(core.Cutoff(2600)).
		Set(core.Gain(0.2))

	// One chord per bar, alternating, so the phrase is two bars long — that
	// longer unit is what gives trance its sense of a build rather than a loop.
	chords := core.Note(mini.Mini("<[a2,c3,e3] [f2,a2,c3]>")).
		Set(core.Cutoff(1400)).
		Set(core.Gain(0.15))

	// Clap on the backbeat, sitting behind the kick.
	clap := core.S(mini.Mini("~ cp ~ cp")).Set(core.Gain(0.4))

	song := core.Stack(kick, bass, hats, hatsClosed, arp, chords, clap)

	shared.Must(song.FastF(shared.Tempo(bpm)), "trance.wav", 8)
}
