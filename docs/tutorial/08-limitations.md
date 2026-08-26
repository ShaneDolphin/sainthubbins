# 8. Limitations

Everything here was verified against this build. Knowing these will save you
from debugging things that are not broken.

## The offline renderer is a sketch pad

`render` and `shared.Must` produce a **mono, single-sine-oscillator** mix. Each
event becomes one sine tone with a short attack and release.

- There are no samples. `bd` is a 60 Hz sine, not a kick drum.
- There is no stereo. Anything about panning is silently ignored.
- There is no reverb, distortion, saturation or delay effect.

It is for checking that a rhythm and a set of pitches work. Judge your patterns
on their timing, not their timbre.

## Only seven controls reach the audio

Measured on this build. These change the WAV:

`Freq`, `Note`, `N`, `S` (pitch), `Gain` (volume), `Cutoff` / `Lpf`
(low-pass).

`Freq` is the one to know about, because it is read *first* and it wins:
give an event a `Freq` and the renderer takes that frequency and skips `N`,
`Note` and `S` entirely. `Gain` and the filter are applied either way. The
one exception is `Freq(220)`, which is indistinguishable from setting no
`Freq` at all — 220 Hz is the renderer's default and it is the value the
skip test looks for, so `core.Note("c3").Set(core.Freq(220))` renders as c3.

Every other control — `Pan`, `Room`, `Speed`, `Shape`, `CRush`, `Attack`,
`Release`, `Resonance`, `Hpf`, and the 300-plus others — is carried in the
event data and ignored by the renderer. Setting them is not an error and does
nothing you can hear.

## Live audio exists, but only through `play` and only with SuperDirt running

`saint-hubbins play <code> [host] [port] [secs]` streams a pattern to
SuperDirt over OSC in real time — see [the README](../../README.md) for the
command table. It needs SuperCollider with SuperDirt already running and
listening on `127.0.0.1:57120`; if nothing is listening, `play` still exits
cleanly (OSC over UDP has no receiver to fail against) but you will hear
nothing.

What's still true: the **web console** (`saint-hubbins serve`) does not
produce sound by itself — it evaluates patterns and returns hap data, nothing
more. `render` still renders to a file rather than speakers. The workflow for
those two remains: change the file, run it, open the WAV.

## The text vocabulary is a subset of the engine

This section used to say `s("bd sd")` and `.fast(2)` were not implemented as
text. They are now — `eval`, `render`, `midi`, `play` and the console all run
your input as JavaScript first and fall back to the mini-notation parser. What
remains true is that the JS vocabulary is **curated, not complete**:

- **Twelve controls** are bound: `s`/`sound`, `note`, `n`, `gain`, `cutoff`,
  `lpf`, `pan`, `room`, `speed`, `attack`, `release`, `shape`. Every other
  control in the vocabulary — `Freq`, `Hpf`, `Crush`, `Resonance` and the
  hundreds beyond them — is Go only.
- **Fourteen transforms** are bound: `fast`, `slow`, `ply`, `segment`, `late`,
  `early`, `degradeBy`, `add`, `rev`, `palindrome`, `degrade`, `hush`, `euclid`,
  `every`. `Jux`, `Chop`, `Striate`, `Struct`, `Sometimes`, `Off`, `Iter`,
  `Zoom`, `Compress`, `LastOf` and the rest are Go only.
- **No signal generators and no music theory.** `.gain(sine)` is a
  `ReferenceError`; the `tonal` scale, chord and voicing helpers have no text
  spelling at all.

The full list is [chapter 7](07-new-song-web.md#the-text-vocabulary); the
binding tables themselves are `internal/jsapi/registry.go`. A name outside it
raises an error rather than doing nothing quietly, so you will find the edge by
hitting it, not by wondering.

One difference in kind, not just coverage: `.fast()` takes a JavaScript number,
so a ratio like `.fast(128/120)` is converted from a *float*, and the resulting
event boundaries are the exact rationals of that float rather than of the ratio
you wrote. `eval 's("bd sd").fast(2/3)'` puts the boundary at
`4503599627370496/6004799503160661`; the same speed built in Go from
`core.NewFraction(2, 3)` puts it at `3/4`. Whole numbers (`.fast(2)`) are exact
either way. See [chapter 3](03-patterns-in-go.md#why-songs-are-still-go).

## Your song must live inside this repository

The engine is in `internal/`, which Go only allows the same module to import.
A song is a new package beside `examples/`. You cannot `go get` Saint Hubbins
into a separate project.

## `Struct` needs Go booleans

```go
core.Note(mini.Mini("c3")).Struct(mini.Mini("t ~ t t"))   // 0 events
```

Mini-notation produces the *string* `"t"`, not `true`. Build the mask in Go:

```go
mask := core.FastCat(core.Pure(true), core.Pure(false), core.Pure(true), core.Pure(true))
core.Note(mini.Mini("c3")).Struct(mask)                   // 3 events
```

## Only `Add` understands control bags

`Add` on an already-wrapped pattern adds into the bag's numeric field
(`note`/`n`/`freq`, in that priority order) and leaves every other control
untouched — see [chapter 5](05-transformations.md#pitch):

```go
core.Note(mini.Mini("0 4 7")).Add(12)   // map[note:12], map[note:16], map[note:19]
```

`Sub`, `Mul`, `Div`, `Mod` and `Pow` do not — each still calls `toFloat` on
the whole value, so any control bag collapses to a bare number and every
other control (gain, pan, whatever else was set) is discarded:

```go
core.Note(60).Sub(12)   // -12, not map[note:48]
```

Do the arithmetic before wrapping in a control, or use `Add` with a negative,
reciprocal, etc. where that's workable.

## Tempo is a ratio, not a setting

The renderer runs at a fixed two seconds per cycle. There is no BPM parameter —
`shared.Tempo(bpm)` computes `bpm/120` and `FastF` scales the pattern. This
works, and it means "tempo" and "play this pattern faster" are the same
operation. Applying `FastF` to one layer changes that layer's rhythm rather than
the song's tempo.

## The WASM build is for embedders, not the console

`make wasm` produces `web/static/saint-hubbins.wasm`, served at
`/static/saint-hubbins.wasm` for anyone embedding the engine in their own
page. The live console does not load it: the console's JavaScript calls
`POST /api/evaluate` over HTTP against the running Go server, the same as
`curl` would. This is a deliberate choice, not an oversight — the console is
one Go binary talking HTTP to itself, with no second (WASM) evaluation path
to keep in sync with the first. Nothing in this repository loads the WASM
binary today; it exists only as a build target for embedding scenarios.

## What is solid

So the list above does not leave the wrong impression, these are dependable:

- Exact rational timing. Events that should coincide always do; nothing drifts.
- The mini-notation grammar in [chapter 2](02-mini-notation.md).
- Every transformation in [chapter 5](05-transformations.md), including
  `Every` and `LastOf` across multi-cycle renders.
- Layering with `Stack`, and control merging with `Set`.
- Pattern data as the real output — the full event stream, including the
  controls the sine renderer ignores, is there for MIDI, OSC and visual
  backends.
- Bad input is reported. `eval`, `render`, `midi` and `play` exit non-zero with
  a real message on unparseable text, and the console returns an `error` field,
  rather than handing back a hap whose value is your source text.
