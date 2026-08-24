# 2. Mini-notation

Mini-notation is the string language for rhythm. It goes inside quotes, and it
is what `eval`, `render` and the web console accept directly.

Everything below was checked against this engine. Where the behaviour differs
from Strudel or TidalCycles, it says so.

## Try things as you read

```console
$ go run ./cmd/saint-hubbins eval "bd ~ sd ~"
```

## The grammar

A bar is divided evenly between whatever you write.

### Sequence — space

```
"bd sd"          bd 0 → 1/2,  sd 1/2 → 1
"bd sd hh cp"    four equal quarters
```

More items means each is shorter. The bar is always full.

### Weight — `@`

```
"bd@3 sd"        bd for three quarters, sd for the last quarter
```

`@n` gives a step `n` times the share of the bar it would otherwise get,
taken from its siblings — the other steps in the same sequence still split
whatever is left evenly. `"bd@3 sd"` divides the bar 3:1, so `bd` holds from
0 to 3/4 and `sd` holds from 3/4 to 1. A step with no `@` counts as weight 1.

### Rest — `~`

```
"bd ~ sd ~"      bd at 0, sd at 1/2, silence between
```

A rest holds its slot open. This is how you place something on beat 3 rather
than beat 2 — and, in dubstep, it is most of the composition.

### Repeat — `*`

```
"bd*4"           four kicks, at 0, 1/4, 1/2, 3/4
"hh*16"          sixteen hats
```

`*` divides the item's own slot, so `"bd*4"` fills the bar but `"bd*2 sd"`
squeezes two kicks into the first half only.

### Slow — `/`

```
"bd/2"           one kick every two bars
```

### Group — `[ ]`

Brackets make several items act as one, so operators apply to the whole group:

```
"bd [sd sd]"     bd for the first half, two snares sharing the second
"[bd sd]*2"      the pair, twice per bar
```

Groups nest.

### Stack — `,`

A comma layers things so they sound **at the same time**:

```
"[bd, hh]"       kick and hat together on beat 1
"[bd*4, hh*8]"   four kicks and eight hats, simultaneously
"[c3,e3,g3]"     a C major chord
```

This is how you write a chord. At the top level the brackets are optional —
`"bd*4, hh*8"` stacks too.

### Alternate — `< >`

Angle brackets play **one item per cycle**, in turn:

```
"<bd sd>"        bd on bar 1, sd on bar 2, repeating
"<c2 eb2>"       one note held for a whole bar each
```

This is how you write something longer than a bar. It is also the cleanest way
to hold a note for a full bar, since the single item fills the cycle.

### Euclidean rhythm — `( )`

```
"bd(3,8)"        3 hits spread as evenly as possible over 8 steps
"bd(3,8,2)"      the same, rotated 2 steps
```

`bd(3,8)` lands at 0, 3/8 and 3/4 — the off-kilter pulse behind a great deal of
electronic music. `(5,8)`, `(3,4)` and `(7,16)` are all worth trying.

### Sample index — `:`

```
"bd:1 sd:2"      selects variant 1 of bd, variant 2 of sd
```

Produces a value carrying both the name and the index.

### Random choice — `|`

```
"bd|sd|hh"       picks one, per cycle
```

### Degrade — `?`

```
"hh*16?"         drops hats at random (50% by default)
"hh*16?0.3"      drops about 30%
```

The fastest way to stop a hat pattern sounding mechanical.

### Polymeter — `{ }`

```
"{bd sd, hh hh hh}"    a 2-step and a 3-step cycle running together
```

The sequences have different lengths and drift against each other.

### Number range — `..`

```
"0 .. 3"         0 1 2 3 as four steps
```

## One difference from Strudel

If you are coming from Strudel or Tidal, this does not behave as you expect.

**`!` subdivides in place rather than adding steps.** `"bd!3 sd"` puts three
kicks inside the *first half* of the bar (at 0, 1/6, 1/3), where Strudel would
give four equal quarters. If you want four equal steps, write `"bd bd bd sd"`.

This is listed in [Limitations](08-limitations.md).

## Reference

| Syntax | Meaning | Example |
|--------|---------|---------|
| ` ` | sequence | `bd sd` |
| `@n` | weight (share of the bar) | `bd@3 sd` |
| `~` | rest | `bd ~ sd ~` |
| `*n` | repeat n times | `bd*4` |
| `/n` | play every n cycles | `bd/2` |
| `[ ]` | group | `bd [sd sd]` |
| `,` | stack / chord | `[bd, hh]`, `[c3,e3,g3]` |
| `< >` | one per cycle | `<bd sd>` |
| `(p,s)` | euclidean | `bd(3,8)` |
| `(p,s,r)` | euclidean, rotated | `bd(3,8,2)` |
| `:n` | sample index | `bd:1` |
| `\|` | random choice | `bd\|sd` |
| `?` | random drop | `hh*16?` |
| `{ }` | polymeter | `{bd sd, hh hh hh}` |
| `..` | number range | `0 .. 3` |

## What mini-notation cannot do

It describes rhythm and pitch, and stops there. It has no volume, no filter, no
tempo, and no way to transform a pattern.

```console
$ go run ./cmd/saint-hubbins eval 's("bd sd")'
```

That does **not** work. `s(...)`, `.fast(...)` and friends are not part of this
language; typing them gets you an event whose value is the literal text. Those
are Go, which is [chapter 3](03-patterns-in-go.md).
