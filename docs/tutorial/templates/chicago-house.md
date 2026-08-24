# Chicago House — 120 BPM

```console
$ go run ./examples/chicago-house
chicago-house.wav — 200 events over 8 bars, 16.0s, peak 0.80
```

Source: [`examples/chicago-house/main.go`](../../../examples/chicago-house/main.go)

Where the [house](house.md) template is polished, this is a drum machine and an
attitude: a hard kick, sharp claps, rattling sixteenth hats, and a bassline
whose filter never sits still.

Two techniques carry the whole track — **DegradeBy** and a **patterned control**.

## Line by line

### Kick and clap

```go
kick := core.S(mini.Mini("bd*4"))

clap := core.S(mini.Mini("~ cp ~ cp")).
	Set(core.Gain(0.7))
```

The same skeleton as house, with the clap louder — `0.7` against house's `0.6`.
That alone shifts it from smooth to hard.

### The hats

```go
hats := core.S(mini.Mini("hh*16")).
	DegradeBy(0.3).
	Set(core.Gain(0.3))
```

Sixteen hats per bar is a wall of sound and, played exactly, unmistakably a
machine. `DegradeBy(0.3)` deletes about three in ten at random, and the gaps are
what make it sound played rather than programmed.

The order matters: degrade first, then set the gain, so the gain applies to
whatever survived.

Delete the `DegradeBy` line and listen. That is the entire difference between a
groove and a metronome.

### The acid line

```go
acid := core.Note(mini.Mini("c2 c2 eb2 c2 g2 c2 bb1 c2")).
	Set(core.Cutoff(mini.Mini("300 1400 500 2000 700 1800 400 1100"))).
	Set(core.Gain(0.7))
```

This is the important one.

The first line is eight notes, mostly the root with a few departures — a classic
303 shape, where the movement comes from the filter rather than the melody.

The second line is a **patterned control**. Instead of one cutoff for the whole
part, `core.Cutoff` receives a pattern of eight values, so each note gets its
own filter setting: dull, bright, dull, very bright, and so on. That sweep is
the acid sound.

Both patterns have eight steps, so they line up one-to-one. Give the cutoff
pattern a different number of steps and it will drift against the notes, which
is sometimes exactly what you want.

### The turnaround

```go
acid = acid.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

Every fourth bar the whole line plays backwards — notes, filter movement and
all. The phrase becomes four bars instead of one.

## Change it

| Change | What you will hear |
|--------|--------------------|
| `DegradeBy(0.3)` → `DegradeBy(0.6)` | sparse, skittering hats |
| remove `DegradeBy` | a rigid machine |
| `"300 1400 500 2000 700 1800 400 1100"` → `"300 400 500 600 700 800 900 1000"` | a smooth rising sweep instead of a jumping one |
| the same → `"2000 300 2000 300 2000 300 2000 300"` | a hard alternating squelch |
| all cutoff values → `"800"` | a single value: the acid character vanishes entirely |
| `"c2 c2 eb2 c2 g2 c2 bb1 c2"` → `"c2*8"` | one note; now the filter is doing *all* the work |
| `Every(4, Rev)` → `Every(2, Rev)` | turnaround twice as often |
| `Gain(0.7)` on clap → `Gain(0.4)` | pulls it back toward house |
| `const bpm = 120` → `127` | later, harder Chicago |

### Make the filter move independently

Give the cutoff a different number of steps from the notes:

```go
Set(core.Cutoff(mini.Mini("300 1400 500 2000 700")))
```

Five filter values against eight notes: the pairing shifts every bar and takes
forty steps to repeat.

## Next

[Techno](techno.md) takes the same relentless pulse somewhere darker.
