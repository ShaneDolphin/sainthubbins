# 7. A new song, in the web console

The web console is the fastest way to try a rhythm. Type, press a button, see
exactly where every event lands. Use it to work out your patterns, then move
them into a Go file to turn them into a track.

## Start it

```console
$ go run ./cmd/saint-hubbins serve
Saint Hubbins console listening on http://localhost:8080
```

Open <http://localhost:8080>. Pick a different port with `serve :9000`.

You get a text box, an **Evaluate** button, and an output panel.

## What it is for

Be clear about the console's scope before you plan a session around it:

- It accepts **mini-notation only** — the rhythm language from chapter 2.
- It **prints events**. It does not play audio.
- It has no layering, no controls, no tempo.

So it is an excellent rhythm sketchpad and not a song editor. Work out the parts
here; assemble them in Go.

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

## A session that works

Try these in order, reading where the events land each time:

| Type this | Look for |
|-----------|----------|
| `bd sd` | two events, halves of the bar |
| `bd ~ sd ~` | same two, still at 0 and 1/2 |
| `bd*4` | an even four |
| `[~ hh]*4` | four hats at 1/8, 3/8, 5/8, 7/8 — all off-beat |
| `[bd*4, hh*8]` | twelve events; the kicks and hats overlap in time |
| `bd(3,8)` | three uneven hits |
| `bd(5,8)` | five — busier, still uneven |
| `<bd sd>` | one event; press Evaluate again, it does not change |
| `[c3,eb3,g3]` | three events all starting at `0/1` — a chord |

The last one is worth dwelling on. Three events sharing a start time is exactly
what a chord is in this engine, and it is why the comma means "at the same
time".

## Seeing note timing

`/api/pianoroll` returns the same events with the times as decimals over two
cycles, which is easier to plot:

```console
$ curl -s -X POST http://localhost:8080/api/pianoroll \
    -H 'Content-Type: application/json' \
    -d '{"code":"c3 e3 g3"}'
{"haps":[{"duration":0.333,"time":0,"value":"c3",...},
         {"duration":0.333,"time":0.333,"value":"e3",...}, ...]}
```

Both endpoints accept `{"code": "..."}` by POST and allow cross-origin requests,
so you can drive them from your own page.

## When nothing seems to happen

If the output is a single event whose `value` is the text you typed, the parser
did not recognise your input and fell back to treating it as a literal:

```json
{"haps":[{"part":"0/1 → 1/1","value":"s(\"bd sd\")","whole":"0/1 → 1/1"}]}
```

That is the symptom of typing Go into the console. `s(...)`, `.fast(...)` and
`.gain(...)` are not mini-notation. Use `bd sd` and add the controls in Go.

An empty `value` means the input could not be parsed at all — check your
brackets.

## Moving it into a track

Once a rhythm works, it goes into Go unchanged, inside the quotes:

```
console:  [~ hh]*4
```

```go
hats := core.S(mini.Mini("[~ hh]*4")).Set(core.Gain(0.3))
```

That is the whole workflow: the console tells you *when*, and Go decides *what*
and *how loud*. Chapter 6 builds the file this goes into.

## Next

[Chapter 8](08-limitations.md) is the honest list of what this engine does not
do yet. It is short, and it will save you time.
