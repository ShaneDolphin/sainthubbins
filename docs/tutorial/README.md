# Saint Hubbins — Tutorial

You have installed Saint Hubbins. This tutorial takes you from a first sound to
writing your own tracks.

Saint Hubbins is a **pattern engine**. A pattern is a pure function of time: you
ask it what is happening between two moments and it answers with a list of
events, called *haps*. Nothing is a recording, and nothing is a timeline —
patterns are values you build, combine and transform like any other data.

## Start here

Work through these in order. Each one is short.

| # | Page | What you get |
|---|------|--------------|
| 1 | [First sounds](01-first-sounds.md) | Check the install, make your first WAV file |
| 2 | [Mini-notation](02-mini-notation.md) | The text language for rhythm — the whole grammar |
| 3 | [Patterns in Go](03-patterns-in-go.md) | Layering, and where text stops being enough |
| 4 | [Controls](04-controls.md) | Pitch, volume, filter — what shapes the sound |
| 5 | [Transformations](05-transformations.md) | How to change a pattern you already have |
| 6 | [A new song, on the command line](06-new-song-cli.md) | Build a track from an empty file |
| 7 | [A new song, in the web console](07-new-song-web.md) | The same, using the browser — and the full text vocabulary |
| 8 | [Limitations](08-limitations.md) | What this engine does not do yet — read before you get stuck |

## The eight templates

Eight complete, runnable tracks, one per style. Start from whichever is closest
to the music you want to make, then change it.

| Style | BPM | Run it | Walkthrough |
|-------|-----|--------|-------------|
| [House](templates/house.md) | 125 | `go run ./examples/house` | Four-on-the-floor, off-beat hats |
| [Chicago House](templates/chicago-house.md) | 120 | `go run ./examples/chicago-house` | Raw 808 and an acid line |
| [Techno](templates/techno.md) | 132 | `go run ./examples/techno` | Relentless, dark, stripped back |
| [Minimal dubstep](templates/minimal-dubstep.md) | 140 | `go run ./examples/minimal-dubstep` | Half-time, built from silence |
| [Maximal dubstep](templates/maximal-dubstep.md) | 140 | `go run ./examples/maximal-dubstep` | Same skeleton, every gap filled |
| [Drum and bass](templates/drum-and-bass.md) | 174 | `go run ./examples/drum-and-bass` | Fast drums over a slow sub |
| [Electronica](templates/electronica.md) | 110 | `go run ./examples/electronica` | Broken beat, melody led |
| [Trance](templates/trance.md) | 138 | `go run ./examples/trance` | Off-beat bass, rolling arpeggio |

Each walkthrough explains every line of its template and lists specific changes
you can make, with what you will hear when you make them.

See [templates/README.md](templates/README.md) for how the eight compare.

## The one thing to know first

There are two languages, and they nest rather than compete.

**Mini-notation** is the quoted string language — `"bd*4"`, `"bd(3,8)"`. It
describes *rhythm and pitch*, and nothing else.

**Everything around it** — what the events mean, how loud, how fast, layered
with what — has two spellings that do the same thing:

```js
stack(s("bd*4"), s("hh*8").gain(0.4))        // text: eval, render, play, console
```

```go
core.Stack(                                   // Go: a song in a file
	core.S(mini.Mini("bd*4")),
	core.S(mini.Mini("hh*8")).Set(core.Gain(0.4)),
)
```

Both produce the same twelve events. The text form is bounded — twelve controls
and fourteen transforms, listed in [chapter 7](07-new-song-web.md#the-text-vocabulary)
— and is what you sketch in. The Go form is the whole engine, and is where songs
live. [Chapter 3](03-patterns-in-go.md) is about that boundary; chapters 4 and 5
are in Go because that is where the full vocabulary is.
