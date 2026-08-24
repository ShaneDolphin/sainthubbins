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
| 3 | [Patterns in Go](03-patterns-in-go.md) | Layering, and why real songs are written in Go |
| 4 | [Controls](04-controls.md) | Pitch, volume, filter — what shapes the sound |
| 5 | [Transformations](05-transformations.md) | How to change a pattern you already have |
| 6 | [A new song, on the command line](06-new-song-cli.md) | Build a track from an empty file |
| 7 | [A new song, in the web console](07-new-song-web.md) | The same, using the browser |
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

There are two ways to write patterns, and they are not equally powerful.

**Mini-notation** is the quoted string language — `"bd*4"`, `"bd(3,8)"`. It
describes *rhythm* and is what the `eval` command and the web console accept.

**The Go API** is everything else: layering, volume, filters, and every
transformation. A song with more than one part is written in Go.

You will use both, usually in the same line — mini-notation inside a Go call:

```go
core.S(mini.Mini("bd*4"))
```

Chapter 3 explains why the split exists.
