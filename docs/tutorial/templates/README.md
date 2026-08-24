# The eight templates

Eight complete tracks, one per style. Each is a single Go file you can run,
read, and change.

```console
$ go run ./examples/house
house.wav — 160 events over 8 bars, 16.0s, peak 0.81
```

Every template renders eight bars to a WAV in the current directory.

## How they compare

The event count is the useful column: it is a direct measure of how busy a
track is, and the two dubstep templates share a tempo and a skeleton while
differing by nearly a factor of seven.

| Template | BPM | Events | Peak | The idea |
|----------|-----|--------|------|----------|
| [House](house.md) | 125 | 160 | 0.81 | Four-on-the-floor with the hats in the gaps |
| [Chicago House](chicago-house.md) | 120 | 200 | 0.80 | A drum machine and a moving filter |
| [Techno](techno.md) | 132 | 261 | 0.83 | Relentless, dark, almost nothing on top |
| [Minimal dubstep](minimal-dubstep.md) | 140 | **42** | 0.58 | Half-time; the silence is the instrument |
| [Maximal dubstep](maximal-dubstep.md) | 140 | **279** | 0.76 | The same skeleton, every gap filled |
| [Drum and bass](drum-and-bass.md) | 174 | 319 | 0.77 | Fast drums, slow bassline |
| [Electronica](electronica.md) | 110 | 215 | 0.74 | Broken beat, melody led |
| [Trance](trance.md) | 138 | 455 | 0.82 | Off-beat bass locked to the kick |

## Read these two together

**Minimal and maximal dubstep** are the same tempo, the same kick, and the same
snare. One has 42 events, the other 279. Open both files side by side: the
difference between them is a complete lesson in arrangement, with the variable
of genre held constant.

## Suggested order

1. **House** — the clearest example of layering. Start here.
2. **Minimal dubstep** — how much you can leave out.
3. **Maximal dubstep** — what filling it back in costs.
4. **Trance** — the interlock between kick and bass.
5. Anything that sounds like the music you want to make.

## Changing them

Each walkthrough ends with a table of specific changes and what you will hear.
The general advice:

- **Change one thing, then re-run.** The feedback loop is two seconds long.
- **Watch the peak.** Over `1.00` clips. Lower a `Gain`, usually a chord's.
- **The rests matter as much as the notes.** Moving a `~` changes more than
  changing a pitch.
- **Nothing is destructive.** Patterns are values; delete a line from the
  `Stack` to mute that layer.

## Making one your own

Copy a template to a new directory and edit freely:

```console
$ cp -r examples/house examples/mytrack
$ go run ./examples/mytrack
```

Change the output filename in the last line so it does not overwrite the
original.

Chapter 6 of the tutorial builds a track from nothing, if you would rather start
from an empty file than from someone else's.
