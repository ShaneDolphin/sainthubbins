# 7. A new song, in the web console

The web console is the fastest way to try a pattern. Type, press a button, see
exactly where every event lands. Use it to work things out, then move them into
a Go file to turn them into a track.

## Start it

```console
$ go run ./cmd/saint-hubbins serve
Saint Hubbins console listening on http://localhost:8080
```

Open <http://localhost:8080>. Pick a different port with `serve :9000`.

You get a text box, an **Evaluate** button, and an output panel.

## What it is for

Be clear about the console's scope before you plan a session around it:

- It accepts **both** input languages. Anything you type is run as JavaScript
  first; if that fails, it is parsed as mini-notation. `bd sd` is not valid JS,
  so bare mini-notation lands in the second path and works exactly as it always
  did.
- Its JS vocabulary is a **curated subset** of the Go API — twelve controls and
  fourteen transforms, listed below. It is not everything the engine can do.
- It **prints events**. It does not play audio. For sound you need `render` (to
  a file) or `play` (to SuperDirt).

So it is a very good sketchpad, and still not a song editor. Work out the parts
here; assemble them in Go ([chapter 3](03-patterns-in-go.md)).

## Reading the output

Type `bd(3,8)` and press Evaluate:

```json
{"haps":[
  {"part":"0/1 → 1/8","value":"bd","whole":"0/1 → 1/8"},
  {"part":"3/8 → 1/2","value":"bd","whole":"3/8 → 1/2"},
  {"part":"3/4 → 7/8","value":"bd","whole":"3/4 → 7/8"}
]}
```

Three hits, at 0, 3/8 and 3/4 of the bar. You can see the rhythm is lopsided
without having to hear it — that is the point of the console.

Compare `bd*4`, which gives a hit at 0, 1/4, 1/2 and 3/4. Even. Dull. Reliable.

Now type `s("bd sd")`:

```json
{"haps":[
  {"part":"0/1 → 1/2","value":{"s":"bd"},"whole":"0/1 → 1/2"},
  {"part":"1/2 → 1/1","value":{"s":"sd"},"whole":"1/2 → 1/1"}
]}
```

Same timing, different `value`. Bare mini-notation gives you a **string** —
`"bd"` is a token, and nothing has said what kind of thing it is. Wrapping it in
`s(...)` gives you a **control bag** — `{"s":"bd"}`, "the sound named bd". That
is the same distinction `core.S` makes in Go, and it is what MIDI, OSC and the
renderer actually read.

Add a control and you can watch it land in the bag:

```json
s("bd*4").gain(0.8)

{"haps":[
  {"part":"0/1 → 1/4","value":{"gain":0.8,"s":"bd"},"whole":"0/1 → 1/4"},
  ... three more, identical but for the time ...
]}
```

## A session that works

Try these in order, reading where the events land and what the values look like.
The first half is rhythm; the second half is everything mini-notation cannot
say.

| Type this | Look for |
|-----------|----------|
| `bd sd` | two events, halves of the bar, values are bare strings |
| `bd ~ sd ~` | same two, still at 0 and 1/2 |
| `bd*4` | an even four |
| `[~ hh]*4` | four hats at 1/8, 3/8, 5/8, 7/8 — all off-beat |
| `[bd*4, hh*8]` | twelve events; the kicks and hats overlap in time |
| `bd(3,8)` | three uneven hits |
| `bd(5,8)` | five — busier, still uneven |
| `<bd sd>` | one event; press Evaluate again, it does not change |
| `[c3,eb3,g3]` | three events all starting at `0/1` — a chord |
| `s("bd sd")` | the same two events, values now `{"s":"bd"}` |
| `s("bd*4").gain(0.8)` | four events, each carrying `gain` as well as `s` |
| `stack(s("bd*4"), s("hh*8").gain(0.4))` | twelve — two named layers instead of one bracketed string |
| `s("bd sd").fast(2)` | four; the pair, twice |
| `s("bd sd hh cp").ply(2)` | eight; each hit doubled in its own slot |
| `s("hh*8").degradeBy(0.5)` | about four or five of the eight, different each time |
| `cat(s("bd*4"), s("sd*4"))` | four kicks; Evaluate queries cycle 0, and `cat` gives each argument a whole cycle |
| `s("bd*4").every(2, x => x.rev())` | four kicks — identical here, because `bd*4` reversed is `bd*4`. Try `s("bd sd hh cp").every(2, x => x.rev())` and read the order |
| `note("c3 eb3 g3")` | three events valued `{"note":"c3"}` — pitch, not drums |

Two of those are worth dwelling on.

`[c3,eb3,g3]` — three events sharing a start time is exactly what a chord is in
this engine, and it is why the comma means "at the same time".

`stack(s("bd*4"), s("hh*8").gain(0.4))` — this is the same twelve events as
`[bd*4, hh*8]`, but each layer now carries what it *means* and how loud it is.
That difference is the whole reason the text layer exists.

## The text vocabulary

Everything the JS evaluator binds. There is nothing else — a name not on this
list raises an error rather than doing nothing, *except* inside an
`every(n, fn)` callback, where a typo cannot raise anything the console can
show you and silently produces no events for the cycles it runs on instead —
see [chapter 8](08-limitations.md#a-typo-inside-an-every-callback-is-silent).

| Kind | Names |
|---|---|
| Controls — usable as a constructor (`gain(0.8)`) or a chained setter (`s("bd").gain(0.8)`) | `s` / `sound`, `note`, `n`, `gain`, `cutoff`, `lpf`, `pan`, `room`, `speed`, `attack`, `release`, `shape` |
| Combinators | `stack`, `cat`, `slowcat`, `fastcat`, `sequence`, `silence()`, `mini("...")` |
| Methods, no argument | `.rev()`, `.palindrome()`, `.degrade()`, `.hush()` |
| Methods, one number | `.fast()`, `.slow()`, `.ply()`, `.segment()`, `.late()`, `.early()`, `.degradeBy()`, `.add()` |
| Methods, two arguments | `.euclid(pulses, steps)`, `.every(n, fn)` |

Three things about it:

**Strings are mini-notation, everywhere.** `s("bd*4")` is chapter 2's grammar
inside chapter 7's function call. So is `.gain("0.2 0.8")` (a gain that changes
halfway through the bar) and `stack("bd sd", s("hh*8"))`. The two languages
nest; they do not compete.

**It is a real JavaScript engine.** Variables, arithmetic, arrow functions and
several statements separated by `;` all work — the value of the last expression
is the pattern:

```js
const kick = s("bd*4");
const hats = s("hh*8").gain(0.4);
stack(kick, hats)
```

That is 12 haps, the same as writing it on one line.

**It is a subset, on purpose.** `jux`, `chop`, `striate`, `struct`, `sometimes`,
`off`, `iter`, `zoom`, `arrange`, `lastOf`, the signal generators (`sine`,
`saw`), the `tonal` scale and chord helpers, and every control outside the
twelve above are Go only. When you want one, you are writing a Go file — which
is where the song was going anyway.

## When something is wrong

Type `s("bd" +` and Evaluate:

```json
{"haps":[],"error":"jsapi: SyntaxError: SyntaxError: (anonymous): Line 1:9 Unexpected end of input (and 1 more errors)"}
```

An `error` field, an empty `haps` list, and a message naming the column. The
same for a name that does not exist (`s("bd").nope()` → `TypeError: Object has
no member 'nope'`) and a bad argument (`s("bd").fast("banana")` → `fast:
argument "banana" is not a number`).

This is worth knowing because it did not always work this way. The console used
to answer unrecognised input with a single hap whose `value` was the text you
had typed — a result that looked like success. If you see that shape in an old
note or an old screenshot, that is what it was.

An empty `haps` with no `error` means the pattern is genuinely silent on this
cycle — `~ ~`, or a `<~ bd>` on the bar where the rest falls. That is a correct
answer, not a failure.

## Seeing note timing

`/api/pianoroll` returns the same events with the times as decimals over two
cycles, which is easier to plot:

```console
$ curl -s -X POST http://localhost:8080/api/pianoroll \
    -H 'Content-Type: application/json' \
    -d '{"code":"note(\"c3 e3 g3\")"}'
{"haps":[{"duration":0.3333333333333333,"part":"0/1 → 1/3","time":0,"value":{"note":"c3"},"whole":"0/1 → 1/3"},
         {"duration":0.3333333333333333,"part":"1/3 → 2/3","time":0.3333333333333333,"value":{"note":"e3"},"whole":"1/3 → 2/3"},
         ... four more, running to time 2.0 ...]}
```

Both endpoints accept `{"code": "..."}` by POST, resolve it the same way (JS
first, mini-notation second), and allow cross-origin requests, so you can drive
them from your own page.

## Moving it into a track

A console line maps onto Go call for call. The rhythm strings cross unchanged;
the function names change case:

```
console:  stack(s("[~ hh]*4").gain(0.3), s("bd*4"))
```

```go
song := core.Stack(
	core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.3)),
	core.S(mini.Mini("bd*4")),
)
```

| Console | Go |
|---|---|
| `s("bd*4")` | `core.S(mini.Mini("bd*4"))` |
| `.gain(0.3)` | `.Set(core.Gain(0.3))` |
| `stack(a, b)` | `core.Stack(a, b)` |
| `.fast(2)` | `.FastF(core.FractionFromInt(2))` |
| `.every(4, x => x.rev())` | `.Every(4, func(p core.Pattern) core.Pattern { return p.Rev() })` |

One thing does **not** carry across identically: tempo. In the console you would
write `.fast(128/120)`; in Go you write `.FastF(shared.Tempo(128))`, which is
the exact fraction 16/15 rather than a float's nearest rational. See
[chapter 3](03-patterns-in-go.md#why-songs-are-still-go).

That is the workflow: the console tells you *when* and *what*, and the Go file
is where it becomes a track you can keep. Chapter 6 builds the file this goes
into.

## Next

[Chapter 8](08-limitations.md) is the honest list of what this engine does not
do yet. It is short, and it will save you time.
