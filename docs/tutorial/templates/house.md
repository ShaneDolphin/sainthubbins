# House — 125 BPM

```console
$ go run ./examples/house
house.wav — 160 events over 8 bars, 16.0s, peak 0.81
```

Source: [`examples/house/main.go`](../../../examples/house/main.go)

House is built on one relationship: the kick marks every beat, and the hats sit
in the gaps between them. Get those two right and it is house before you add
anything else.

## Line by line

### The kick

```go
kick := core.S(mini.Mini("bd*4"))
```

`bd*4` is four kicks per bar, evenly spaced. Four-on-the-floor — the thing
everything else is measured against.

No `Gain`, so it plays at full level. The kick anchors the mix; everything else
is set relative to it.

### The clap

```go
clap := core.S(mini.Mini("~ cp ~ cp")).
	Set(core.Gain(0.6))
```

Four slots, two of them rests. The claps land on beats **2 and 4** — the
backbeat.

The rests are the whole trick. Write `"cp ~ cp ~"` instead and the claps move to
1 and 3, doubling the kick rather than answering it. It stops sounding like
dance music immediately. Try it.

### The hats

```go
hats := core.S(mini.Mini("[~ hh]*4")).
	Set(core.Gain(0.3))
```

Read the brackets from the inside out. `[~ hh]` is a two-step group: a rest,
then a hat. `*4` repeats that group four times per bar. Each hat therefore falls
exactly halfway between two kicks — at 1/8, 3/8, 5/8 and 7/8.

That off-beat placement is the second half of the house signature.

### The bass

```go
bass := core.Note(mini.Mini("a1 ~ a1 c2")).
	Set(core.Cutoff(700)).
	Set(core.Gain(0.7))
```

`core.Note`, not `core.S`, because these are pitches rather than drum names.

The line is in A minor and deliberately leaves beat 2 empty, so it pushes
against the kick instead of doubling it. The notes fall on beats 1, 3 and 4. `Cutoff(700)` rolls off the top so the
bass stays underneath the hats.

### The chord stabs

```go
stabs := core.Note(mini.Mini("~ [a3,c4,e4] ~ [a3,c4,e4]")).
	Set(core.Cutoff(2200)).
	Set(core.Gain(0.18))
```

`[a3,c4,e4]` is a comma-stack — three notes at the same instant, an A minor
triad. They land on beats 2 and 4, off the kick.

The gain is low because a chord is three voices at once; `0.18` here is roughly
as loud as `0.5` on a single note.

### Assembly

```go
song := core.Stack(kick, clap, hats, bass, stabs)

shared.Must(song.FastF(shared.Tempo(bpm)), "house.wav", 8)
```

`Stack` layers all five. `FastF(shared.Tempo(125))` sets the tempo on the
finished stack — 125/120, applied once at the end.

## Change it

One change at a time, then re-run.

| Change | What you will hear |
|--------|--------------------|
| `const bpm = 125` → `128` | tighter, more club-ready |
| `const bpm = 125` → `118` | deep house drag |
| `"~ cp ~ cp"` → `"cp ~ cp ~"` | the backbeat collapses onto the kick — instructive and bad |
| `"[~ hh]*4"` → `"hh*8"` | straight eighths; the off-beat push disappears |
| `"[~ hh]*4"` → `"[~ hh]*4"` with `.DegradeBy(0.2)` | looser, more human |
| `"a1 ~ a1 c2"` → `"a1 ~ c2 e2"` | a walking bassline |
| `"a1 ~ a1 c2"` → `"a1*8"` | driving, relentless — closer to techno |
| `Cutoff(700)` → `Cutoff(300)` | dark, muffled bass |
| `Cutoff(700)` → `Cutoff(4000)` | bright and thin; it stops being a bass |
| `[a3,c4,e4]` → `[a3,c#4,e4]` | A **major** — the whole mood lifts |
| `[a3,c4,e4]` → `[a3,c4,e4,g4]` | A minor 7th, jazzier (lower the gain to `0.14`) |
| `"~ [a3,c4,e4] ~ [a3,c4,e4]"` → `"~ [a3,c4,e4] ~ ~"` | one stab per bar, more space |
| remove `stabs` from `Stack` | a stripped-back dub mix |

### Add a four-bar phrase

The template loops every bar. One line gives it a longer shape:

```go
bass = bass.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

### Add an echo to the stabs

```go
stabs = stabs.Off(0.125, func(p core.Pattern) core.Pattern {
	return p.Set(core.Gain(0.07))
})
```

A quieter copy an eighth-note later.

## Next

[Chicago House](chicago-house.md) is the same tempo range and a completely
different attitude.
