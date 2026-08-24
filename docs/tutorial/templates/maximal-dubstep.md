# Maximal dubstep — 140 BPM

```console
$ go run ./examples/maximal-dubstep
maximal-dubstep.wav — 279 events over 8 bars, 16.0s, peak 0.76
```

Source: [`examples/maximal-dubstep/main.go`](../../../examples/maximal-dubstep/main.go)

The same 140 BPM, the same half-time skeleton, and the same kick and snare as
[minimal dubstep](minimal-dubstep.md) — with everything the minimal version
leaves out crammed into the gaps. 42 events becomes 279.

Open both files side by side. The variable of genre is held constant, so what
you are looking at is purely arrangement.

## The unchanged skeleton

```go
kick := core.S(mini.Mini("bd ~ ~ ~"))
snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))
```

Identical to the minimal template, line for line. Everything below is addition.

## Line by line

### Ghost kicks

```go
ghostKick := core.S(mini.Mini("~ ~ ~ [bd bd]")).Set(core.Gain(0.55))
```

`[bd bd]` is a group of two kicks sharing the final quarter of the bar, so they
arrive twice as fast as the grid. At `0.55` they sit behind the main kick as a
stumble into the next bar.

This is where the minimal template has its single quiet hat.

### Busy hats

```go
hats := core.S(mini.Mini("hh*16")).
	DegradeBy(0.3).
	Set(core.Gain(0.22))
```

Sixteen per bar against the minimal version's one. `DegradeBy(0.3)` keeps them
ragged instead of mechanical.

### The wobble

```go
wobble := core.Note(mini.Mini("c1*8")).
	Set(core.Cutoff(mini.Mini("180 900 260 1500 200 1100 320 1800"))).
	Set(core.Gain(1.0))
```

The centrepiece, and the same technique as the Chicago house acid line put to a
different use.

Eight sub notes per bar, all the same pitch. The movement is entirely in the
**patterned cutoff**: each note gets its own filter value, alternating low and
high — 180, 900, 260, 1500. That alternation between closed and open is what
makes a dubstep bass sound like it is talking.

Compare with the minimal template's single sustained `<c2 eb2>` at a fixed
`Cutoff(200)`. Same instrument, opposite philosophy.

### The half/double switch

```go
wobble = wobble.Every(2, func(p core.Pattern) core.Pattern { return p.Ply(2) })
```

`Ply(2)` repeats each event twice *inside its own time slot*, so every second
bar the wobble runs at double speed without any of the notes moving. This is the
standard dubstep gear-change, in one line.

### The lead

```go
lead := core.Note(mini.Mini("~ ~ [g4,bb4] ~")).
	Set(core.Cutoff(3000)).
	Set(core.Gain(0.3))
```

A two-note chord an octave above the minimal template's, at `Cutoff(3000)`
instead of `1600` — bright and shrieking, landing on the snare.

### The fill

```go
fill := core.Silence().
	LastOf(4, func(core.Pattern) core.Pattern {
		return core.S(mini.Mini("~ ~ ~ [sd sd sd sd]")).Set(core.Gain(0.6))
	})
```

The same `LastOf` structure the minimal template uses for its one quiet chord —
here it delivers four rapid snares at the end of every four-bar phrase.

## Change it

| Change | What you will hear |
|--------|--------------------|
| `"180 900 260 1500 200 1100 320 1800"` → `"180 1800"` | a slow two-step wobble |
| the same → `"180 400 700 1000 1300 1600 1800 2000"` | a rising sweep instead of a wobble |
| `"c1*8"` → `"c1*16"` | twice the wobble rate (halve the cutoff pattern's steps to match) |
| `Ply(2)` → `Ply(4)` | a quadruple-time gear change |
| `Every(2, Ply)` → `Every(4, Ply)` | the switch arrives half as often |
| `[sd sd sd sd]` → `[sd sd sd sd sd sd]` | a longer, more frantic fill |
| `LastOf(4, ...)` → `LastOf(2, ...)` | a fill every other bar |
| remove `ghostKick` and `lead` from the `Stack` | most of the way back to minimal |
| `Gain(0.22)` on hats → `Gain(0.4)` | the hats take over the top end |
| `[g4,bb4]` → `[g4,bb4,d5]` | a fuller lead (lower the gain to `0.22`) |

### The exercise worth doing

Delete layers from the `Stack` one at a time:

```go
song := core.Stack(kick, snare, ghostKick, hats, wobble, lead, fill)
```

Remove `fill`, then `lead`, then `hats`, then `ghostKick`, re-running each time.
You will arrive at the minimal template by subtraction, and hear exactly what
each layer was contributing.

## Next

[Drum and bass](drum-and-bass.md) takes the tempo up to 174.
