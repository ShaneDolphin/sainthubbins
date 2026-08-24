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

## Only five controls reach the audio

Measured on this build. These change the WAV:

`Note`, `N`, `S` (pitch), `Gain` (volume), `Cutoff` / `Lpf` (low-pass).

Every other control — `Pan`, `Room`, `Speed`, `Shape`, `CRush`, `Attack`,
`Release`, `Resonance`, `Hpf`, and the remaining 280-odd — is carried in the
event data and ignored by the renderer. Setting them is not an error and does
nothing you can hear.

## There is no live audio

Saint Hubbins renders to a file. There is no command that plays a pattern out of
your speakers in real time, and the web console does not produce sound. The
workflow is: change the file, run it, open the WAV.

## The console evaluates mini-notation only

`s("bd sd")`, `.fast(2)`, `.gain(0.5)` are **not** implemented as text. There is
no script evaluator wired up, so anything the mini-notation parser does not
recognise comes back as a literal string value.

Layering, controls and transformations are the Go API. This is the single most
common source of confusion — see [chapter 3](03-patterns-in-go.md).

## Your song must live inside this repository

The engine is in `internal/`, which Go only allows the same module to import.
A song is a new package beside `examples/`. You cannot `go get` Saint Hubbins
into a separate project.

## Mini-notation differences from Strudel

| Syntax | Expected | Actual here |
|--------|----------|-------------|
| `bd!3 sd` | four equal steps | three bds inside the *first half* |
| `{a b, c d e}%4` | four steps per cycle | **no effect** — the suffix is ignored |

Workaround: write repeats out in full (`"bd bd bd sd"`) instead of using `!`.

## `Struct` needs Go booleans

```go
core.Note(mini.Mini("c3")).Struct(mini.Mini("t ~ t t"))   // 0 events
```

Mini-notation produces the *string* `"t"`, not `true`. Build the mask in Go:

```go
mask := core.FastCat(core.Pure(true), core.Pure(false), core.Pure(true), core.Pure(true))
core.Note(mini.Mini("c3")).Struct(mask)                   // 3 events
```

## `Add` does not transpose a wrapped pattern

```go
core.Note(mini.Mini("0 4 7")).Add(12)   // every event becomes plain 12
```

`Add` on a pattern already wrapped in a control replaces the control bag with a
bare number. Add first, wrap second:

```go
core.Note(core.Pure(60).Add(12))        // map[note:72]
```

Or just write the note names you want.

## Tempo is a ratio, not a setting

The renderer runs at a fixed two seconds per cycle. There is no BPM parameter —
`shared.Tempo(bpm)` computes `bpm/120` and `FastF` scales the pattern. This
works, and it means "tempo" and "play this pattern faster" are the same
operation. Applying `FastF` to one layer changes that layer's rhythm rather than
the song's tempo.

## The WASM build is not wired to the console

`make wasm` produces `web/static/saint-hubbins.wasm`, and the console footer
mentions it, but the page never loads it — the console calls the Go server over
HTTP instead. The WASM target builds and is unused.

## What is solid

So the list above does not leave the wrong impression, these are dependable:

- Exact rational timing. Events that should coincide always do; nothing drifts.
- The mini-notation grammar in [chapter 2](02-mini-notation.md), minus the three
  rows above.
- Every transformation in [chapter 5](05-transformations.md), including
  `Every` and `LastOf` across multi-cycle renders.
- Layering with `Stack`, and control merging with `Set`.
- Pattern data as the real output — the full event stream, including the
  controls the sine renderer ignores, is there for MIDI, OSC and visual
  backends.
