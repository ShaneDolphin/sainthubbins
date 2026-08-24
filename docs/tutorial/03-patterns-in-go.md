# 3. Patterns in Go

Mini-notation gives you one rhythm. A song is several rhythms at once, at a
chosen tempo, with the parts at different volumes. That is Go.

## Why the split exists

The `eval` command hands your string to the mini-notation parser and prints the
result. There is no interpreter for method calls, so `s("bd sd")` is not code —
it is just text the parser does not recognise, and you get an event whose value
is the string `s("bd sd")`.

Everything beyond rhythm is the Go API. In practice you write mini-notation
*inside* Go calls, and the two fit together in one line:

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
