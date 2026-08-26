# 3. Patterns in Go

Mini-notation gives you one rhythm. A song is several rhythms at once, at a
chosen tempo, with the parts at different volumes.

That used to be the whole reason for this chapter: layering, volume and tempo
were Go and nothing else. They are not any more —
`stack(s("bd*4"), s("hh*8").gain(0.4))` is something you can type at `eval`, at
`render`, or into the console. So the honest question is no longer "why can't I
type this?" but "where does typing stop being enough?"

## Where the line is now

Two things are true at once, and the tutorial is clearer if you hold both.

**Text reaches further than it used to.** The same pattern vocabulary is bound
into a JavaScript evaluator (`internal/jsapi`), so the shapes that matter for
sketching — sounds, pitch, a dozen controls, layering, and fourteen transforms —
have a text spelling. That is [chapter 7](07-new-song-web.md#the-text-vocabulary).

**Go is still the whole engine, and text is a window onto part of it.**

| | Text (`eval`, `render`, `play`, console) | Go |
|---|---|---|
| Rhythm | mini-notation | mini-notation, same parser |
| Sound and pitch | `s`, `sound`, `note`, `n` | `core.S`, `core.Note`, `core.N` |
| Controls | twelve: `gain`, `cutoff`, `lpf`, `pan`, `room`, `speed`, `attack`, `release`, `shape` and the pitch three | every control in the vocabulary |
| Layering | `stack`, `cat`, `slowcat`, `fastcat`, `sequence` | those, plus `core.Arrange`, and the methods `.Jux()` / `.Superimpose()` |
| Transforms | fourteen, including `fast`, `slow`, `rev`, `ply`, `euclid`, `every` | the whole transformation core — `Chop`, `Striate`, `Struct`, `Sometimes`, `Off`, `Iter`, `Zoom`, `Compress`, `LastOf`, … |
| Modulating a control with a signal | — | `core.Sine()`, `core.Saw()`, `core.Perlin()` |
| Scales, chords, voicings | — | `internal/tonal` |
| Rendering | `render out.wav '<code>'`, fixed at 4 cycles | `shared.Must(song, "x.wav", bars)`, any length, with the event count and peak reported |

## Why songs are still Go

Three reasons, in order of how often they bite.

**1. A song is a program you keep.** Text is a line you type and lose. A file
has named parts you can reuse, a comment explaining the bar you always forget,
and a git history. Every template in this tutorial is a Go file for that reason,
not because text could not express it.

**2. Tempo is exact in Go and approximate in text.** `.fast()` takes a JS
number, and a JS number is a float. `core.FractionFromFloat` turns it into the
exact rational *of that float*, which is not the rational you meant:

```console
$ go run ./cmd/saint-hubbins eval 's("bd sd").fast(2/3)'
...
    "part": "0/1 → 4503599627370496/6004799503160661",
...
    "part": "4503599627370496/6004799503160661 → 1/1",
...
```

(Only the `part` lines are shown.) That boundary is 3/4, to about sixteen
digits. `shared.Tempo(128)` is
`core.NewFraction(128, 120)` — exactly 16/15, no float in the path. Whole
numbers (`.fast(2)`) are exact either way; it is the ratios that drift.

**3. The vocabulary runs out.** The table above is the map. When you want a
`Jux`, a `Chop`, a scale, or one of the controls the text layer does not bind,
you are in Go — and moving there is cheap, because mini-notation goes across
unchanged, inside the quotes.

In practice you write mini-notation *inside* Go calls, and the two fit together
in one line — the same nesting the text layer does with `s("bd*4")`:

```go
core.S(mini.Mini("bd*4"))
//     ^^^^^^^^^^^^^^^^^ the rhythm
// ^^^^^^ what the events mean: a sound named bd
```

## Where your code has to live

The engine lives in `internal/`, and Go only allows `internal/` packages to be
imported from inside the same module. **Your song must be a package inside this
repository** — a new directory alongside `examples/`. You cannot import Saint
Hubbins into an unrelated project.

## The smallest complete song

Create `examples/mybeat/main.go`:

```go
package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	kick := core.S(mini.Mini("bd*4"))
	hats := core.S(mini.Mini("hh*8")).Set(core.Gain(0.4))

	song := core.Stack(kick, hats)

	shared.Must(song.FastF(shared.Tempo(128)), "mybeat.wav", 4)
}
```

Run it:

```console
$ go run ./examples/mybeat
mybeat.wav — 53 events over 4 bars, 8.0s, peak 0.42
```

Five lines of music. Here is what each does.

## Line by line

### `core.S(mini.Mini("bd*4"))`

`mini.Mini` parses the string into a pattern of bare values — four events, each
holding the text `bd`. `core.S` wraps each value into a **control**: the event
becomes `map[s:bd]`, meaning "the sound named bd".

`core.Note` does the same for pitch: `core.Note(mini.Mini("c2 g2"))` produces
`map[note:c2]` and `map[note:g2]`.

Without one of these wrappers the renderer does not know whether `c2` is a drum
or a note. Chapter 4 covers the rest of the controls.

### `.Set(core.Gain(0.4))`

`Set` merges another control into every event. The hats become
`map[gain:0.4 s:hh]`.

Because layers **sum** when rendered, this is how you keep one part from
swallowing the others. It is the single most useful call in the API.

### `core.Stack(kick, hats)`

Layering. All the patterns play at once, and the result is just another pattern.

The alternatives sequence rather than layer:

| Call | Effect |
|------|--------|
| `core.Stack(a, b)` | a and b together |
| `core.FastCat(a, b)` | a then b, both inside one bar |
| `core.SlowCat(a, b)` | a for a whole bar, then b for the next |

`SlowCat` gives each argument a full cycle in turn, cycling through the list —
useful for arranging distinct sections back to back. For an *occasional* event
inside an otherwise steady pattern — a fill every fourth bar, a bar of
silence — reach for `Every` and `LastOf` from
[chapter 5](05-transformations.md) instead; that is what both dubstep templates
do.

### `.FastF(shared.Tempo(128))`

Tempo, applied to the finished stack.

One cycle is two seconds, so an untouched pattern is 120 BPM. Every other tempo
is a ratio: `Tempo(128)` is `128/120`, or `16/15`. `FastF` speeds the pattern up
by that factor.

```go
song.FastF(shared.Tempo(174))   // drum and bass
song.FastF(core.NewFraction(3, 2))  // one and a half times as fast
```

Apply it once, at the end, to the whole stack. Applying it to one layer changes
that layer's rhythm rather than the song's tempo — occasionally useful, usually
a bug.

### `shared.Must(song, "mybeat.wav", 4)`

Renders four bars and writes the WAV, reporting the event count and peak level.
It is a ten-line convenience in `examples/shared/`, not part of the engine.

Watch the peak. Above `1.00` the audio clips; the templates sit around `0.8`.
If yours clips, lower a `Gain`.

## Patterns are values

This is the idea the whole engine rests on. A pattern is not a recording or a
track — it is a value, so it can be stored, passed to a function, and reused:

```go
hats := core.S(mini.Mini("hh*8")).Set(core.Gain(0.4))

quiet := hats.Set(core.Gain(0.1))    // hats is unchanged
busy  := hats.FastF(core.NewFraction(2, 1))

song := core.Stack(kick, quiet, busy.Late(0.125))
```

Every method returns a **new** pattern and leaves the original alone, so you can
build variations from one part without copying it.

## Next

[Chapter 4](04-controls.md) covers the controls that shape the sound, and which
of them you can actually hear.
