# 4. Controls

A control turns a bare value into an instruction. `mini.Mini("bd*4")` gives four
events holding the text `bd`; `core.S` makes them four events meaning *play the
sound bd*.

```go
core.S(mini.Mini("bd*4"))       // map[s:bd]
core.Note(mini.Mini("c2 g2"))   // map[note:c2]
```

## Combining controls

`Set` merges controls into every event of a pattern:

```go
core.S(mini.Mini("hh*8")).
    Set(core.Gain(0.4)).
    Set(core.Cutoff(3000))
// map[cutoff:3000 gain:0.4 s:hh]
```

The argument can itself be a pattern, which is how you make a value move:

```go
core.Note(mini.Mini("c1*8")).
    Set(core.Cutoff(mini.Mini("180 900 260 1500 200 1100 320 1800")))
```

Eight bass notes, each with its own filter setting. That is a dubstep wobble,
and it is the technique behind the maximal dubstep template.

## The four that make sound

The offline renderer is a single sine oscillator per event. Only these change
what you hear. This was measured, not assumed:

| Control | Effect | Range |
|---------|--------|-------|
| `core.Note(v)` | pitch, by note name or MIDI number | `c1`–`c6`, `0`–`127` |
| `core.N(v)` | pitch as a MIDI number | `0`–`127` |
| `core.S(v)` | pitch, by drum name | see below |
| `core.Gain(v)` | volume | `0`–`2`, default `1` |
| `core.Cutoff(v)` | low-pass filter, in Hz (`core.Lpf` is identical) | `20`–`20000` |

When more than one is present, `n` wins over `note`, which wins over `s`.

### Drum names

The renderer maps a few names to fixed frequencies:

| Name | Frequency | Use |
|------|-----------|-----|
| `bd` | 60 Hz | kick |
| `sd` | 180 Hz | snare |
| `hh`, `oh`, `ch` | 800 Hz | hats |
| `cp` | 300 Hz | clap |

Anything else becomes a tone derived from the name. These are sine tones, not
samples — `bd` is a low thud, not a kick drum.

### Note names

`c2`, `f#3`, `eb4`, `a1`. Letter, optional `#` or `b`, then octave. Middle C is
`c4`. Bass parts live around `c1`–`c2`, leads around `c4`–`c5`.

## The ones that are carried but silent

Saint Hubbins defines 295 controls, and your events can carry any of them. The
offline renderer ignores everything except the five above.

Measured as having **no effect on the WAV**: `Pan`, `Room`, `Speed`, `Shape`,
`CRush`, `Attack`, `Release`, `Resonance`, `Hpf`, and the rest.

They are not useless. They travel with the event data, which is what MIDI and
OSC backends consume, and they will matter if the renderer grows. But if you set
`core.Room(0.8)` and hear no reverb, the pattern is fine — the sine renderer
simply does not implement it. See [Limitations](08-limitations.md).

## Watch the level

Layers **sum**. Each voice reaches about `0.3`, so seven layers hitting together
can exceed full scale and clip.

`shared.Must` prints the peak:

```
trance.wav — 455 events over 8 bars, 16.0s, peak 0.82
```

Keep it below `1.00`. The templates sit near `0.8`. If yours is over, lower a
`Gain` rather than every gain — usually one loud layer is responsible, and
chords are the usual culprit because a three-note chord is three voices.

## A worked example

```go
// Kick: no gain set, so it plays at full level and anchors the mix.
kick := core.S(mini.Mini("bd*4"))

// Hats: quiet, so they sit behind the kick.
hats := core.S(mini.Mini("hh*16")).Set(core.Gain(0.2))

// Bass: filtered down so it stays out of the way of the hats.
bass := core.Note(mini.Mini("<c2 eb2>")).
    Set(core.Cutoff(400)).
    Set(core.Gain(0.8))

// Chord: three voices at once, so its gain is lower than it looks.
stab := core.Note(mini.Mini("~ [c3,eb3,g3] ~ ~")).
    Set(core.Cutoff(2000)).
    Set(core.Gain(0.2))

song := core.Stack(kick, hats, bass, stab)
```

## Next

[Chapter 5](05-transformations.md) is about changing a pattern once you have it.
