# Techno — 132 BPM

```console
$ go run ./examples/techno
techno.wav — 261 events over 8 bars, 16.0s, peak 0.83
```

Source: [`examples/techno/main.go`](../../../examples/techno/main.go)

Techno keeps house's four-to-the-floor pulse and throws away its warmth. Faster,
darker, and with almost nothing on top — the interest comes from texture and
from things not quite landing on the grid.

## Line by line

### The kick

```go
kick := core.S(mini.Mini("bd*4"))
```

Identical to house. At 132 BPM rather than 125, it stops feeling like an
invitation and starts feeling like a demand.

### Two layers of hats

```go
hatsClosed := core.S(mini.Mini("hh*16")).Set(core.Gain(0.18))
hatsOpen := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.4))
```

This is worth copying into your own tracks.

The first layer is sixteen quiet hats — at `0.18` you do not hear them
individually, you hear a continuous hiss underneath everything.

The second is four louder off-beat hats on top, giving the push.

Two layers at different densities and volumes read as one textured part. One
layer cannot do this.

### The stab

```go
stab := core.Note(mini.Mini("c3")).
	Euclid(3, 8).
	Set(core.Cutoff(900)).
	Set(core.Gain(0.4))
```

A single note, placed by a Euclidean rhythm: `Euclid(3, 8)` spreads three hits
as evenly as it can across eight steps, landing them at 0, 3/8 and 3/4.

That is deliberately *not* the grid. Against a rigid four-to-the-floor kick, the
three-against-eight pull is most of what makes the track interesting.

`Cutoff(900)` keeps it dark rather than bright.

### The sub

```go
sub := core.Note(mini.Mini("c2 ~ c2 ~")).
	Set(core.Cutoff(280)).
	Set(core.Gain(0.8))
```

Two notes per bar under the kick, filtered almost shut at 280 Hz so it is felt
rather than heard.

### The clap

```go
clap := core.S(mini.Mini("~ ~ cp ~")).
	Every(2, func(p core.Pattern) core.Pattern { return core.Silence() }).
	Set(core.Gain(0.5))
```

`Every(2, silence)` replaces the pattern with nothing on every second bar, so
the clap plays on alternate bars only. A two-bar phrase from one line.

This is a useful idiom in general: `Every(n, silence)` thins any layer out.

## Change it

| Change | What you will hear |
|--------|--------------------|
| `Euclid(3, 8)` → `Euclid(5, 8)` | busier, more urgent |
| `Euclid(3, 8)` → `Euclid(7, 16)` | a long, lopsided phrase |
| `Euclid(3, 8)` → `Euclid(4, 8)` | lands on the grid — and instantly sounds ordinary |
| `Cutoff(900)` on the stab → `Cutoff(3000)` | a bright lead instead of a dark stab |
| `Cutoff(280)` on the sub → `Cutoff(600)` | the sub becomes audible as a note |
| `Gain(0.18)` on closed hats → `Gain(0.35)` | the hiss becomes a rhythm |
| `"c2 ~ c2 ~"` → `"c2*4"` | driving, no space |
| `"c3"` on the stab → `"[c3,g3]"` | a fifth — heavier |
| `Every(2, silence)` → `Every(4, silence)` | clap on three bars out of four |
| `const bpm = 132` → `140` | hard techno |
| `const bpm = 132` → `125` | back toward house |

### Add a filter sweep

Make the sub open up across the bar:

```go
sub = sub.Set(core.Cutoff(mini.Mini("280 400 600 900")))
```

### Make the stab answer itself

```go
stab = stab.Off(0.125, func(p core.Pattern) core.Pattern {
	return p.Set(core.Gain(0.2)).Set(core.Cutoff(2000))
})
```

A brighter, quieter echo an eighth later.

## Next

[Minimal dubstep](minimal-dubstep.md) — the opposite approach entirely.
