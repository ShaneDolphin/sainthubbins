# Drum and bass — 174 BPM

```console
$ go run ./examples/drum-and-bass
drum-and-bass.wav — 319 events over 8 bars, 16.0s, peak 0.77
```

Source: [`examples/drum-and-bass/main.go`](../../../examples/drum-and-bass/main.go)

Fast drums over a slow bassline. The kit runs at 174 BPM while the sub changes
note once a bar, and the tension between those two speeds is the genre.

## Line by line

### The two-step

```go
kick := core.S(mini.Mini("bd ~ ~ ~"))
snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.9))
```

Kick on one, snare on three. On paper this is identical to the
[dubstep](minimal-dubstep.md) skeleton — the same four slots, the same two hits.

The tempo is the entire difference. At 140 it reads as half-time and heavy; at
174 the same placement reads as fast and driving. Change `bpm` to `140` and the
track becomes dubstep without touching a single pattern.

### Ghost notes

```go
ghosts := core.S(mini.Mini("~ [~ bd] [~ sd] [bd ~]")).Set(core.Gain(0.4))
```

This is what separates a break from a drum machine.

Each `[...]` group splits its quarter of the bar into two halves. `[~ bd]` is a
rest then a kick, so the kick lands on the second eighth. `[bd ~]` is the
reverse. The result is quiet extra hits scattered off the main beats, at `0.4`
so they sit underneath.

Set the gain to `0.9` for a moment to hear where they actually are.

### Rolling hats

```go
hats := core.S(mini.Mini("hh*16")).
	DegradeBy(0.2).
	Set(core.Gain(0.2))
```

Sixteenths at 174 BPM are very fast, which is what carries the tempo.
`DegradeBy(0.2)` drops one in five so they roll rather than buzz.

### The slow sub

```go
sub := core.Note(mini.Mini("<c2 g1>")).
	Set(core.Cutoff(220)).
	Set(core.Gain(1.0))
```

`<c2 g1>` holds one note for an entire bar, then the other. Against sixteenth
hats, the bass is moving thirty-two times more slowly than the top of the kit.

That ratio is the point. A fast bassline under fast drums is noise; a slow one
gives the ear something to hold onto.

### The offset stab

```go
stab := core.Note(mini.Mini("[c3,eb3,g3]")).
	Euclid(3, 8).
	Late(0.0625).
	Set(core.Cutoff(1800)).
	Set(core.Gain(0.22))
```

A C minor chord on a three-against-eight Euclidean rhythm, then `Late(0.0625)`
nudges the whole thing a sixteenth later.

The offset is doing real work. `Euclid(3, 8)` alone puts hits at 0, 3/8 and 3/4 —
and the first of those lands exactly on the kick. Shifted a sixteenth, the hits
move to 1/16, 7/16 and 13/16, which misses the kick at 0 and the snare at 1/2
entirely. `Late` is worth remembering: small offsets are the difference between
parts that collide and parts that interlock. Try `Late(0.125)` to hear the stab
double the snare instead.

### The turnaround

```go
ghosts = ghosts.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

Every fourth bar the ghost pattern reverses — a fill that costs one line.

## Change it

| Change | What you will hear |
|--------|--------------------|
| `const bpm = 174` → `140` | the identical patterns become dubstep |
| `const bpm = 174` → `160` | jungle territory |
| `Gain(0.4)` on ghosts → `Gain(0.9)` | the break moves to the front; useful for hearing it |
| `"~ [~ bd] [~ sd] [bd ~]"` → `"~ [~ sd] [~ sd] [sd ~]"` | a snare-led break |
| `<c2 g1>` → `<c2 g1 eb2 f1>` | a four-bar bassline |
| `<c2 g1>` → `"c2 ~ g1 ~"` | the sub now moves within the bar — busier, less anchored |
| `Late(0.0625)` → `Late(0.125)` | the stab collides with the snare on every bar |
| remove `Late` entirely | the stab lands on the kick |
| `Euclid(3, 8)` → `Euclid(5, 16)` | a longer, sparser stab phrase |
| `DegradeBy(0.2)` → `DegradeBy(0.5)` | half the hats gone; more space |
| `Cutoff(220)` → `Cutoff(500)` | the sub becomes an audible bassline |

### Add a reese

Two detuned notes together give the classic drum and bass bass:

```go
sub := core.Note(mini.Mini("<[c2,c2] [g1,g1]>")).
	Set(core.Cutoff(mini.Mini("220 400 220 600"))).
	Set(core.Gain(0.8))
```

The moving cutoff does the growling.

## Next

[Electronica](electronica.md) slows everything down to 110.
