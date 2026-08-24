# Trance — 138 BPM

```console
$ go run ./examples/trance
trance.wav — 455 events over 8 bars, 16.0s, peak 0.82
```

Source: [`examples/trance/main.go`](../../../examples/trance/main.go)

The busiest of the eight — 455 events, nearly twelve times the minimal dubstep
template. It is built on one trick: the kick lands on the beat, and the bass
answers in every single gap between kicks.

## The interlock

```go
kick := core.S(mini.Mini("bd*4"))

bass := core.Note(mini.Mini("[~ a1]*4")).
	Set(core.Cutoff(600)).
	Set(core.Gain(0.6))
```

Read those two together, because they are one idea.

`bd*4` puts kicks at 0, 1/4, 1/2 and 3/4. `[~ a1]` is a rest then a note, and
`*4` repeats it four times — so bass notes land at 1/8, 3/8, 5/8 and 7/8.

Every bass note falls exactly halfway between two kicks. Nothing ever collides,
and the two parts together read as a single rolling engine. That is the trance
bassline, and it is why it is written as `[~ a1]*4` rather than as four separate
notes.

Try `"a1*4"` instead: the bass now lands *on* the kicks, everything muddies, and
the engine stops.

## Line by line

### Two hat layers

```go
hats := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.3))
hatsClosed := core.S(mini.Mini("hh*16")).Set(core.Gain(0.1))
```

The same two-layer technique as [techno](techno.md): loud off-beat hats
reinforcing the bass placement, over a quiet sixteenth hiss at `0.1` that fills
the space between everything.

### The arpeggio

```go
arp := core.Note(mini.Mini("a3 c4 e4 a4 e4 c4 a3 c4")).
	FastF(core.NewFraction(2, 1)).
	Set(core.Cutoff(2600)).
	Set(core.Gain(0.2))
```

Eight notes climbing an A minor triad and coming back down, then `FastF(2)`
doubles the speed to sixteen notes per bar.

Writing it as eight notes and doubling is easier to read and edit than writing
sixteen — and it means one number controls the whole arpeggio's rate.

`Cutoff(2600)` keeps it bright and on top of the mix.

### The chords

```go
chords := core.Note(mini.Mini("<[a2,c3,e3] [f2,a2,c3]>")).
	Set(core.Cutoff(1400)).
	Set(core.Gain(0.15))
```

Angle brackets around two comma-stacks: A minor for one bar, F major for the
next, each sustained through its whole cycle.

That two-bar alternation is what makes trance feel like it is going somewhere
rather than repeating. The gain is only `0.15` because each chord is three
voices at once.

### The clap

```go
clap := core.S(mini.Mini("~ cp ~ cp")).Set(core.Gain(0.4))
```

Backbeat, well behind the kick at `0.4`. With this much else happening it is
reinforcement rather than a feature.

## Change it

This template is close to its ceiling at peak `0.82`. If you make something
louder, make something else quieter.

| Change | What you will hear |
|--------|--------------------|
| `"[~ a1]*4"` → `"a1*4"` | the interlock breaks; bass and kick collide |
| `"[~ a1]*4"` → `"[~ a1]*8"` | a double-time bass — harder, more driving |
| `FastF(2, 1)` on arp → `FastF(1, 1)` | eight notes instead of sixteen; much calmer |
| the same → `FastF(4, 1)` | frantic (drop the gain to `0.12`) |
| `"a3 c4 e4 a4 e4 c4 a3 c4"` → `"a3 c4 e4 a4 c5 a4 e4 c4"` | the arpeggio reaches higher |
| `[f2,a2,c3]` → `[e2,g#2,b2]` | an E major turn — the classic uplifting change |
| `<[a2,c3,e3] [f2,a2,c3]>` → `<[a2,c3,e3] [f2,a2,c3] [c3,e3,g3] [g2,b2,d3]>` | a four-bar chord progression |
| `Cutoff(600)` on bass → `Cutoff(1200)` | the bass gets a growl |
| `Gain(0.1)` on closed hats → `Gain(0.2)` | the sixteenth layer becomes audible as rhythm |
| `const bpm = 138` → `142` | harder, more euphoric |
| `const bpm = 138` → `132` | closer to progressive |

### Build a breakdown

Trance is arrangement as much as pattern. Drop the drums for the first two bars
of every four:

```go
drums := core.Stack(kick, clap, hats, hatsClosed)
drums = drums.Every(4, func(core.Pattern) core.Pattern { return core.Silence() })

song := core.Stack(drums, bass, arp, chords)
```

One bar in four with no kit at all — arpeggio and chords alone — then the drums
return. Use `Every(2, ...)` to alternate every other bar instead.

### Open the filter across the phrase

```go
arp = arp.Set(core.Cutoff(mini.Mini("<1200 1800 2400 3200>")))
```

One cutoff value per bar, brightening over four bars — a filter build.

## Next

Back to [the template index](README.md), or the
[tutorial](../README.md) if you skipped ahead.
