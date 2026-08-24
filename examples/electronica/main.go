// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Electronica — 110 BPM.
//
//	go run ./examples/electronica
//
// Slower and more melodic than the dancefloor templates. The kick no longer
// marks every beat, which leaves room for the melody to carry the rhythm.

package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 110

	// A broken kick: three hits spread evenly across eight steps rather than
	// four square beats. This is what stops it sounding like house.
	kick := core.S(mini.Mini("bd")).Euclid(3, 8)

	// Soft backbeat.
	snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.5))

	hats := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.22))

	// A four-note arpeggio in A minor. <> alternates one note per cycle, and
	// *8 speeds that alternation up eightfold — so the four notes run twice
	// inside every bar instead of taking four bars to arrive.
	arp := core.Note(mini.Mini("<a3 c4 e4 g4>*8")).
		Set(core.Cutoff(2400)).
		Set(core.Gain(0.3))

	// Off makes a copy shifted an eighth later and quieter — a delay built
	// out of pattern operations rather than an effect.
	arp = arp.Off(0.125, func(p core.Pattern) core.Pattern {
		return p.Set(core.Gain(0.12))
	})

	// A held pad underneath, one chord per bar.
	pad := core.Note(mini.Mini("<[a2,c3,e3] [f2,a2,c3]>")).
		Set(core.Cutoff(900)).
		Set(core.Gain(0.3))

	bass := core.Note(mini.Mini("a1 ~ ~ f1")).
		Set(core.Cutoff(500)).
		Set(core.Gain(0.85))

	song := core.Stack(kick, snare, hats, arp, pad, bass)

	shared.Must(song.FastF(shared.Tempo(bpm)), "electronica.wav", 8)
}
