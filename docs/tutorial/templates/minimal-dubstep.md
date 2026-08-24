# Minimal dubstep — 140 BPM

```console
$ go run ./examples/minimal-dubstep
minimal-dubstep.wav — 42 events over 8 bars, 16.0s, peak 0.58
```

Source: [`examples/minimal-dubstep/main.go`](../../../examples/minimal-dubstep/main.go)

**42 events over eight bars.** The [maximal](maximal-dubstep.md) template has
279 at the same tempo with the same kick and snare. Read the two together — the
difference between them is the whole lesson.

## The half-time feel

Dubstep runs at 140 BPM and *feels* like 70. The trick is placement: the kick
lands on beat one and the snare waits until beat three, so a bar that is
technically fast reads as slow.

```go
kick := core.S(mini.Mini("bd ~ ~ ~"))
snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))
```

Four slots each, one filled. Six rests across two lines. Almost everything here
is silence, and that is the composition — not a gap waiting to be filled.

## Line by line

### The hat

```go
hat := core.S(mini.Mini("~ ~ ~ hh")).Set(core.Gain(0.25))
```

One hat per bar, on beat four, quiet. The only thing marking time between the
snare and the next kick. Remove it and the bar loses its edge; double it and the
half-time feel starts to slip.

### The sub

```go
sub := core.Note(mini.Mini("<c2 eb2>")).
	Set(core.Cutoff(200)).
	Set(core.Gain(1.0))
```

`<c2 eb2>` plays **one item per bar**: c2 through the whole of bar one, eb2
through bar two, repeating. Because the angle brackets hold a single item for
the full cycle, the note is sustained rather than struck — one long tone under
everything.

`Cutoff(200)` removes everything but the fundamental. This is felt more than
heard, which is correct for a sub.

To go deeper, drop it an octave to `<c1 eb1>`. Most laptop speakers will not
reproduce it, which is also correct for a sub.

### The one detail

```go
chord := core.Note(mini.Mini("~ ~ ~ [g3,bb3]")).
	Set(core.Cutoff(1600)).
	Set(core.Gain(0.2))
detail := core.Silence().
	LastOf(4, func(core.Pattern) core.Pattern { return chord })
```

`LastOf(4, ...)` starts from silence and swaps in the chord when the cycle is
the fourth of its group. Three bars pass with nothing at all, the chord lands on
the fourth, and the phrase begins again.

This is how you write an arrangement rather than a loop. It is also the entire
top end of the track — one two-note chord, once every four bars.

## Change it

The temptation is to add. Resist it long enough to hear what removal does.

| Change | What you will hear |
|--------|--------------------|
| remove `hat` from the `Stack` | starker; the bar loses its marker |
| `"~ ~ ~ hh"` → `"~ hh ~ hh"` | time starts moving; the half-time feel weakens |
| `"bd ~ ~ ~"` → `"bd ~ bd ~"` | no longer half-time — it becomes a straight beat |
| `<c2 eb2>` → `<c2 eb2 f2 d2>` | a four-bar bass phrase instead of two |
| `<c2 eb2>` → `"c2*8"` | eight sub hits: the space vanishes entirely |
| `Cutoff(200)` → `Cutoff(800)` | the sub becomes a bass line you can hear pitch in |
| `LastOf(4, ...)` → `LastOf(2, ...)` | the detail every two bars — busier |
| `[g3,bb3]` → `[g3,bb3,d4]` | a fuller chord (drop the gain to `0.15`) |
| `const bpm = 140` → `150` | toward drum and bass territory |

### The one addition worth making

A reversed sub tail leading into each phrase:

```go
sub = sub.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

### Then compare

```console
$ go run ./examples/maximal-dubstep
```

Same tempo. Same kick. Same snare. Nearly seven times the events.

## Next

[Maximal dubstep](maximal-dubstep.md).
