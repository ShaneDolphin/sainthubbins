# Electronica — 110 BPM

```console
$ go run ./examples/electronica
electronica.wav — 215 events over 8 bars, 16.0s, peak 0.74
```

Source: [`examples/electronica/main.go`](../../../examples/electronica/main.go)

Slower and more melodic than the dancefloor templates. The kick stops marking
every beat, which leaves room for the melody to carry the rhythm instead.

## Line by line

### The broken kick

```go
kick := core.S(mini.Mini("bd")).Euclid(3, 8)
```

This one line is why the track is not house.

`Euclid(3, 8)` spreads three kicks as evenly as possible across eight steps,
landing them at 0, 3/8 and 3/4. Three hits, unevenly spaced, instead of four
square beats — so the pulse is implied rather than stated.

Change it to `Euclid(4, 8)` and you get four-on-the-floor and a completely
different genre.

### Soft backbeat

```go
snare := core.S(mini.Mini("~ ~ sd ~")).Set(core.Gain(0.5))
hats := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.22))
```

Both quieter than in the dance templates — `0.5` and `0.22`. In a melody-led
track the drums support rather than drive.

### The arpeggio

```go
arp := core.Note(mini.Mini("<a3 c4 e4 g4>*8")).
	Set(core.Cutoff(2400)).
	Set(core.Gain(0.3))
```

Worth reading carefully, because it combines two operators.

`<a3 c4 e4 g4>` on its own plays **one note per bar**, taking four bars to get
through the sequence. `*8` speeds that alternation up eightfold, so the four
notes run twice inside every bar — eight arpeggio notes per bar, identical each
time.

That is the general shape of `<...>*n`: the angle brackets choose the order, the
`*n` chooses the speed.

### The echo

```go
arp = arp.Off(0.125, func(p core.Pattern) core.Pattern {
	return p.Set(core.Gain(0.12))
})
```

`Off` keeps the original and adds a copy shifted an eighth-note later, here at
`0.12` gain — well under half the original's `0.3`.

That is a delay effect built entirely out of pattern operations. No effect unit
is involved; there are simply twice as many events now, and the extra ones are
quieter and late.

### The pad

```go
pad := core.Note(mini.Mini("<[a2,c3,e3] [f2,a2,c3]>")).
	Set(core.Cutoff(900)).
	Set(core.Gain(0.3))
```

Angle brackets around two comma-stacks: an A minor triad for one whole bar, then
an F major triad for the next. One chord per bar, each sustained across the
cycle.

`Cutoff(900)` keeps the pad soft so it sits behind the arpeggio.

### The bass

```go
bass := core.Note(mini.Mini("a1 ~ ~ f1")).
	Set(core.Cutoff(500)).
	Set(core.Gain(0.85))
```

Two notes per bar, following the pad's chord changes — A under the A minor, F
under the F major.

## Change it

| Change | What you will hear |
|--------|--------------------|
| `Euclid(3, 8)` → `Euclid(4, 8)` | four-on-the-floor; it becomes house |
| `Euclid(3, 8)` → `Euclid(5, 8)` | busier and more urgent |
| `<a3 c4 e4 g4>*8` → `<a3 c4 e4 g4>*4` | a slower, calmer arpeggio |
| the same → `<a3 c4 e4 g4>*16` | frantic, glittering |
| the same → `<a3 c4 e4 g4>` | one note per bar — a slow melody instead of an arpeggio |
| `Off(0.125, ...)` → `Off(0.25, ...)` | a wider, more obvious echo |
| `Gain(0.12)` in the echo → `Gain(0.25)` | the echo competes with the original |
| remove the `Off` block | dry and much emptier |
| `[f2,a2,c3]` → `[d2,f2,a2]` | a D minor turn — sadder |
| `"a1 ~ ~ f1"` → `"a1 ~ f1 ~"` | the chord change arrives earlier |
| `Cutoff(2400)` on arp → `Cutoff(800)` | the arpeggio moves behind the pad |
| `const bpm = 110` → `95` | downtempo |

### Make the arpeggio wander

Give it a longer cycle so it does not repeat every bar:

```go
arp := core.Note(mini.Mini("<a3 c4 e4 g4 b3 d4>*8"))
```

Six notes at eight per bar: the phrase now takes three bars to come round.

## Next

[Trance](trance.md) — the tightest kick-and-bass interlock of the eight.
