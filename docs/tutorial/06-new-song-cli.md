# 6. A new song, on the command line

This chapter starts with an empty file and ends with a track. Every step is
something you can run and hear.

The finished file ships as [`examples/mytrack/main.go`](../../examples/mytrack/main.go)
if you want to compare, but type it out — the point is the order things arrive
in.

## Set up

```console
$ mkdir -p examples/mytrack
```

Create `examples/mytrack/main.go`:

```go
package main

import (
	"codeberg.org/uzu/saint-hubbins/examples/shared"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

func main() {
	const bpm = 128

	kick := core.S(mini.Mini("bd*4"))

	song := core.Stack(kick)

	shared.Must(song.FastF(shared.Tempo(bpm)), "mytrack.wav", 8)
}
```

```console
$ go run ./examples/mytrack
mytrack.wav — 35 events over 8 bars, 16.0s, peak 0.30
```

Kicks over eight bars at 128 BPM. Open the WAV. It is dull, and it is a
foundation.

(Thirty-five, not thirty-two: at 128 BPM the pattern runs slightly faster than
one bar per rendered bar, so a few extra kicks fit in the same sixteen seconds.)

> If you get a build error about `internal`, your directory is outside this
> repository. See chapter 3.

## Step 2 — the backbeat

A kick alone has no shape. Add something answering it:

```go
clap := core.S(mini.Mini("~ cp ~ cp")).
    Set(core.Gain(0.6))
```

and put it in the stack:

```go
song := core.Stack(kick, clap)
```

```console
$ go run ./examples/mytrack
mytrack.wav — 52 events over 8 bars, 16.0s, peak 0.48
```

The two rests are doing the work: they put the clap on beats 2 and 4 instead of
1 and 3. Try `"cp ~ cp ~"` to hear how wrong the other placement sounds.

## Step 3 — hats

```go
hats := core.S(mini.Mini("[~ hh]*4")).
    DegradeBy(0.15).
    Set(core.Gain(0.3))
```

`[~ hh]` is a rest then a hat, four times per bar, so each hat lands exactly
between two kicks. `DegradeBy(0.15)` throws away about one in seven at random,
which is the difference between a drum machine and a person.

Add it to the stack and listen — 79 events, peak 0.57. Then delete the
`DegradeBy` line and listen again: the pattern is the same, and it sounds worse.

## Step 4 — bass

```go
bass := core.Note(mini.Mini("c2 ~ c2 eb2")).
    Set(core.Cutoff(600)).
    Set(core.Gain(0.7))
```

Note the difference from the drums: `core.Note`, not `core.S`, because these are
pitches. `Cutoff(600)` rolls the top off so the bass sits under the hats instead
of fighting them.

The rest on beat 2 matters as much as the notes — the line lands on beats 1, 3 and 4. That is 105 events, peak 0.76.

## Step 5 — a chord

```go
chord := core.Note(mini.Mini("~ [c3,eb3,g3] ~ ~")).
    Set(core.Cutoff(1800)).
    Set(core.Gain(0.18))
```

`[c3,eb3,g3]` is a comma-stack: three notes at the same moment, a C minor triad.
It sits on the off-beat and nowhere else, which is why one chord per bar is
enough.

The gain looks very low because a three-note chord is three voices — it is
roughly three times as loud as the number suggests.

## Step 6 — stop it looping

Eight identical bars is a loop, not a track. One line fixes it:

```go
bass = bass.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })
```

Every fourth bar the bassline runs backwards, so the phrase is four bars long
rather than one.

```console
$ go run ./examples/mytrack
mytrack.wav — 132 events over 8 bars, 16.0s, peak 0.78
```

## Where to go from here

Change one thing at a time and re-run:

| Change | What happens |
|--------|--------------|
| `const bpm = 128` → `140` | the whole track speeds up |
| `"c2 ~ c2 eb2"` → `"c2 ~ eb2 g2"` | a busier bassline |
| `"c2"` → `"f2"` throughout | a different key |
| `Cutoff(600)` → `Cutoff(300)` | a darker, duller bass |
| `[c3,eb3,g3]` → `[c3,e3,g3]` | C **major** — the whole mood flips |
| `Every(4, Rev)` → `Every(2, Rev)` | variation twice as often |
| `[~ hh]*4` → `hh*16` | driving instead of loose |
| add `.Ply(2)` to `clap` | a double clap |

## Watch the peak

`shared.Must` prints it on every run. Above `1.00` the audio clips. If that
happens, lower the gain of whichever layer you just made louder — usually a
chord.

## Next

[Chapter 7](07-new-song-web.md) does the same exploration in the browser, which
is faster for trying rhythms.
