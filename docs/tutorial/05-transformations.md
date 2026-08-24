# 5. Transformations

This is the chapter to reach for when you have a pattern and want it to sound
different. Every method returns a new pattern and leaves the original alone.

All examples start from:

```go
b := mini.Mini("bd sd hh cp")   // events at 0, 1/4, 1/2, 3/4
```

## Time

### `Rev()` — backwards

```go
b.Rev()     // cp@0  hh@1/4  sd@1/2  bd@3/4
```

### `FastF(f)` / `SlowF(f)` — speed

```go
b.FastF(core.NewFraction(2, 1))   // the whole bar twice: 8 events
b.SlowF(core.NewFraction(2, 1))   // stretched over two bars
```

Use these on a single layer to change its rhythm, or on the finished stack to
set the tempo (chapter 3).

### `Late(t)` / `Early(t)` — nudge

```go
b.Late(0.25)    // everything a quarter-bar later, wrapping round
b.Early(0.125)  // an eighth earlier
```

Small values are the difference between a stiff pattern and one with a groove.
`Late(0.02)` on a hat layer is a swing.

### `Ply(n)` — repeat each event in place

```go
b.Ply(2)    // bd bd sd sd hh hh cp cp, each pair inside the original slot
```

A drum roll without changing where the beats fall. This is the half-time /
double-time switch in the maximal dubstep template.

### `Segment(n)` — chop a held value into n events

```go
mini.Mini("c3").Segment(2)   // two c3 events instead of one long one
```

Turns a sustained note into a repeated one.

## Structure

### `Every(n, fn)` — vary on a cycle

```go
b.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

Applies `fn` on every nth bar. The single most effective way to stop a loop
sounding like a loop: reverse every fourth bar, or double-time every second.

### `LastOf(n, fn)` — vary on the *last* cycle of each group

```go
core.Silence().LastOf(4, func(core.Pattern) core.Pattern { return fill })
```

Where `Every(4, ...)` fires on cycles 0, 4, 8, `LastOf(4, ...)` fires on 3, 7,
11 — the end of each four-bar phrase, which is where a fill belongs.

Starting from `core.Silence()` and swapping *in* a pattern is the idiom for
"play this only occasionally". Both dubstep templates use it.

### `Superimpose(fn)` — layer a variation on top

```go
b.Superimpose(func(p core.Pattern) core.Pattern { return p.FastF(core.NewFraction(2, 1)) })
```

Keeps the original **and** adds the transformed copy — twelve events here, not
eight: the four originals plus eight from the doubled-speed copy.

### `Off(t, fn)` — a delayed copy

```go
b.Off(0.125, func(p core.Pattern) core.Pattern { return p.Set(core.Gain(0.3)) })
```

Superimpose, shifted later. A quieter copy an eighth behind is an echo built out
of pattern operations. This is the arpeggio delay in the electronica template.

### `Palindrome()` — forwards then backwards

```go
b.Palindrome()   // bar 1 forwards, bar 2 reversed
```

### `Jux(fn)` — the stereo split

```go
b.Jux(func(p core.Pattern) core.Pattern { return p.Rev() })
```

Intended to put the original in one ear and the transformed copy in the other.
The offline renderer is **mono and ignores pan**, so here it behaves like
`Superimpose`. It is correct in the event data, and will matter through a stereo
backend.

## Rhythm

### `Euclid(pulses, steps)`

```go
core.Note(mini.Mini("c3")).Euclid(3, 8)   // hits at 0, 3/8, 3/4
```

Spreads `pulses` hits as evenly as possible across `steps`. `(3,8)`, `(5,8)` and
`(7,16)` all produce rhythms that sound composed rather than programmed. Same as
`"c3(3,8)"` in mini-notation.

### `Struct(binary)` — take the rhythm from somewhere else

Plays a pattern only where a boolean pattern is true:

```go
onOffOn := core.FastCat(core.Pure(true), core.Pure(false), core.Pure(true), core.Pure(true))
core.Note(mini.Mini("c3")).Struct(onOffOn)   // 3 events: at 0, 1/2, 3/4
```

It needs **Go booleans**. Writing the mask in mini-notation as `"t ~ t t"` does
not work, because that produces the strings `"t"`, not `true`, and you get no
events at all.

### `DegradeBy(p)` — drop events at random

```go
core.S(mini.Mini("hh*16")).DegradeBy(0.3)   // about 30% removed
core.S(mini.Mini("hh*16")).Degrade()        // 50%, the default
```

Sixteen identical hats sound like a machine. Drop a third and they sound played.
Measured: 64 events over four bars becomes 45 at `0.3`, 37 at `0.5`.

## Pitch

To transpose, change the note names — `"c2 eb2 g2"` to `"d2 f2 a2"`. That always
works and is the clearest thing to read.

To transpose arithmetically, work in MIDI numbers or scale degrees. `Add` on an
already-wrapped pattern transposes it in place — it adds into the bag's
`note` field and leaves every other control untouched:

```go
core.Note(mini.Mini("0 4 7")).Add(12)   // map[note:12], map[note:16], map[note:19]
core.Note(core.Pure(60)).Add(12)        // map[note:72] — C4 up an octave
```

This works because the values above are numbers (or numeric mini-notation
strings like `"0"`). Note *names* such as `"c3"` are kept as strings in the
control bag until the sound engine parses them, so `Add` cannot do arithmetic
on them — it leaves a named note exactly as it was rather than guessing:

```go
core.Note(mini.Mini("c3 e3 g3")).Add(12)   // unchanged: "c3", "e3", "g3"
```

Add does not raise an error here, and it does not merge the three notes into
one — each stays its own untouched string. It just does nothing to them. Use
note names when you want to write pitches directly, and numbers when you want
to transpose them.

## Choosing one

| You want | Use |
|----------|-----|
| less repetitive over time | `Every(4, ...)` |
| a busier version of a part | `Superimpose` or `Ply` |
| an echo | `Off` |
| a less mechanical hat line | `DegradeBy(0.3)` |
| an off-grid pulse | `Euclid(3, 8)` |
| a groove | `Late(0.02)` |
| a fill at the end of each phrase | `core.Silence().LastOf(4, ...)` |
| a bar of silence every few bars | `Every(4, ...)` returning `core.Silence()` |
| one rhythm driving another | `Struct(booleanPattern)` |
| a different tempo | `FastF(shared.Tempo(bpm))` on the stack |

## Next

You have the whole vocabulary. [Chapter 6](06-new-song-cli.md) builds a track
from an empty file.
