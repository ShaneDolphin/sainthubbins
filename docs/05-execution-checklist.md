# Execution Checklist — Saint Hubbins

Run the automated gates:

```
./scripts/check.sh
```

Ten gates run. `go vet` and the race-enabled test suite are judged on exit
status — that is what those tools report. The other **eight are behavioural
and assert on content**, so none of them can pass vacuously:

| Gate | What it asserts |
|---|---|
| `wasm builds` | the freshly built artefact starts with the `\0asm` magic number — `make wasm`'s own `\|\|` fallback swallows a failed build and still exits 0 |
| `eval produces real haps (mini-notation)` | `eval "bd sd"` prints haps valued `bd` and `sd`, and exactly 2 of them |
| `eval produces real haps (JS pattern code)` | `eval 's("bd sd")'` prints haps carrying the **control field** `"s": "bd"` / `"s": "sd"`, not one hap whose value is the source text. Grepping the quoted control field is the whole point: the literal string `s("bd sd")` contains `bd` too, so a bare-substring grep would pass either way |
| `render writes genuine audio` | RIFF/WAVE header **and** sample data that is not all zero |
| `midi writes a genuine Standard MIDI File` | `MThd` header **and** a nonzero note count in the CLI's own report line |
| `play reaches a socket (OSC to SuperDirt)` | `TestRunPlaySendsToSuperDirt` really ran (its `--- PASS:` line is grepped) and real OSC bytes carrying `bd` arrived on a real UDP socket |
| `all nine example templates build and run` | nine per-template subtests passed — the count is asserted, so adding or dropping a template fails here rather than quietly making the gate name a lie |
| `no leftover Strudel dependency/import path` | no `@strudel/`, `strudel-go` or `codeberg.org/uzu/strudel` in tracked files outside the documented exemptions |

The two gates that delegate to `go test -run` grep the `-v` output for the
named test's own `--- PASS:` line, because `go test -run` prints "no tests to
run" and exits **0** when the pattern matches nothing — a renamed or deleted
test would otherwise turn its gate green forever. See the comments in
`scripts/check.sh` for the specific bug each gate was built to catch (a
vacuously passing `eval` gate, a `make wasm` target that swallows a failed
build, and a stale import path in a generated-parser source file).

## Manual gates

A script cannot verify these — they require a human, a running SuperDirt
instance, a DAW, or a browser:

- [ ] `saint-hubbins play` is audible: with SuperDirt running and listening,
      `go run ./cmd/saint-hubbins play "bd*4"` should be heard, not just
      measured as bytes on a socket (which `./scripts/check.sh` already
      does).
- [ ] The `.mid` file `./scripts/check.sh` writes (or
      `go run ./cmd/saint-hubbins midi out.mid "bd sd"`) opens in a DAW and
      plays back the expected notes.
- [ ] `go run ./cmd/saint-hubbins serve` boots the console at
      `http://localhost:8080` and the page works in a browser. Both input
      paths, since the console accepts both: the default textarea contents
      (a JS `stack(s("bd*4"), s("hh*8").gain(0.4))` under two comment lines)
      evaluate to 12 haps carrying `s`/`gain` controls, and replacing the
      whole box with bare mini-notation
      (`[bd*4, hh*8]`) evaluates to 12 haps carrying bare string values.
      Typing something that is neither (`s("bd" +`) shows an `error` field,
      not a hap whose value is the text you typed.

Legacy checklist archived at `docs/archive/strudel-legacy/05-execution-checklist.md`.
