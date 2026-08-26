# Execution Checklist — Saint Hubbins

Run the automated gates:

```
./scripts/check.sh
```

This runs `go vet`, the race-enabled test suite, the WASM build (verified by
its magic number, not just `make`'s exit code — see the script's comments),
the CLI's `eval`/`render`/`midi` subcommands (asserted on real output
content, not just a zero exit status), the `play` → SuperDirt OSC path, all
nine tutorial templates, and the rebrand check for leftover Strudel
dependencies/import paths. The seven behavioural gates assert on content;
none of them can pass vacuously. (`go vet` and the race-enabled suite are
judged on exit status — that is what those tools report.) The two gates
that delegate to `go test -run` grep the `-v` output
for the named test's own `--- PASS:` line, because `go test -run` prints
"no tests to run" and exits **0** when the pattern matches nothing — a
renamed or deleted test would otherwise turn its gate green forever. See the
comments in `scripts/check.sh` for what each gate checks and why, including
the specific bugs it was built to catch (a vacuously passing `eval` gate, a
`make wasm` target that swallows a failed build, and a stale import path in a
generated-parser source file).

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
      `http://localhost:8080` and the page works in a browser: typing
      mini-notation and evaluating it shows haps in the console.

Legacy checklist archived at `docs/archive/strudel-legacy/05-execution-checklist.md`.
