# 1. First sounds

## Check the install

From the repository root:

```console
$ go build ./...
```

No output means everything compiled. Saint Hubbins needs Go 1.25 or newer and
nothing else — no Node, no audio server.

Now ask the binary what it can do:

```console
$ go run ./cmd/saint-hubbins
Usage: saint-hubbins <eval|serve|render|play|query|midi> [args]
  eval <code>        — evaluate pattern string
  query              — demo query: Stack(s("bd"), s("sd"))
  serve [addr]       — start live console server (default :8080)
  render <out.wav> <code> — offline render to WAV
  play <code> [host] [port] [secs] — stream to SuperDirt over OSC
  midi <out.mid> <code> [cycles] — render to a Standard MIDI File
  (also available as 'hubbins' — these go to eleven)
```

(It exits with status 1, since you gave it no command.)

Six commands. You will use `eval` to look at patterns, `render` to hear them,
and `serve` for the browser console.

## Look at a pattern

`eval` takes a pattern and prints the events it produces over one cycle. The
simplest thing to hand it is a mini-notation string — the rhythm language of
[chapter 2](02-mini-notation.md):

```console
$ go run ./cmd/saint-hubbins eval "bd sd"
[
  {
    "part": "0/1 → 1/2",
    "value": "bd",
    "whole": "0/1 → 1/2"
  },
  {
    "part": "1/2 → 1/1",
    "value": "sd",
    "whole": "1/2 → 1/1"
  }
]

2 haps
```

Read that carefully, because it is the whole model:

- A **cycle** is one bar. It runs from `0/1` to `1/1`.
- Each event is a **hap**: a value, and the span of time it occupies.
- `bd` occupies `0/1 → 1/2` — the first half of the bar. `sd` takes the second.
- Times are **exact fractions**, not floating-point. Two events written to land
  together always land together; timing never drifts.

Nothing said how long a bar lasts, or what a `bd` sounds like. A pattern is only
*what happens when*. Turning that into sound is a separate step, which is why
the same pattern can drive audio, MIDI or graphics.

Add a rest with `~` and watch the timing hold:

```console
$ go run ./cmd/saint-hubbins eval "bd ~ sd ~"
```

Four slots, two of them silent. `sd` still starts at `1/2`.

## Hear it

`render` writes a WAV file:

```console
$ go run ./cmd/saint-hubbins render first.wav "c3 e3 g3 c4"
wrote first.wav (384000 samples)
```

Open `first.wav` in any audio player. Four rising tones over eight seconds.

That is four cycles at two seconds each, which matters:

> **One cycle is two seconds.** Four beats in two seconds is **120 BPM**, so an
> untouched pattern plays at 120 BPM. Every other tempo is a ratio against it.
> Chapter 3 shows how to set one.

Try a drum pattern instead:

```console
$ go run ./cmd/saint-hubbins render beat.wav "bd ~ sd ~"
```

### What you are hearing

Be clear about this early so nothing later disappoints you: the offline renderer
is a **sine-wave sketch pad**, not a drum machine. Each event becomes a single
sine tone with a short envelope. `bd` is 60 Hz, `sd` is 180 Hz, `hh` is 800 Hz.
It is enough to check that your rhythm works, and it is not going to sound like a
record.

The real output of Saint Hubbins is the pattern data. Audio is one thing you can
do with it. See [Limitations](08-limitations.md).

## Playing it live

`render` bounces sine tones offline; `play` streams the real pattern data
over OSC to [SuperDirt](https://github.com/musikinformatik/SuperDirt), which
is what actually turns it into sound with real samples and synths:

```console
$ go run ./cmd/saint-hubbins play "bd sd" 127.0.0.1 57120 4
```

This only produces sound if **SuperCollider, with SuperDirt already
started**, is listening on that host and port (57120 by default) — if you
hear nothing, that is almost always why.

One trap worth knowing before you type a pattern like `"0 1 2 3"`: a bare
number in mini-notation is stored as a string, so `play` sends it as a
**sample name** (`s "0"`), not a note number. For a sample index use `bd:3`
syntax; for real note numbers, name the control — `play 'n("0 1 2 3")'`, or
`core.N` in Go.

That `n(...)` is the other thing `eval`, `render` and `play` accept: JS pattern
code, tried before the mini-notation parser. `s("bd sd").fast(2)` and
`stack(...)` work anywhere a pattern string does. [Chapter 7](07-new-song-web.md)
has the whole vocabulary; for now, mini-notation is enough.

## Next

You have seen mini-notation do two things: sequence (`bd sd`) and rest (`~`).
[Chapter 2](02-mini-notation.md) covers the rest of the grammar.
