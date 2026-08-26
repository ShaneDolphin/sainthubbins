# Implementation Status — Saint Hubbins (as of 2026-08-23)

Archived. Historical implementation notes from the Strudel-Go prototype era have been moved to `docs/archive/strudel-legacy/IMPLEMENTATION_STATUS.md`.

Current status is whatever `./scripts/check.sh` reports: ten gates covering `go vet`, the race-enabled test suite, the WASM build, `eval` in both of its input forms (bare mini-notation and JS pattern code), `render`, `midi`, `play`, all nine tutorial templates, and a rebrand check. The eight behavioural gates assert on output content, not just a zero exit status; `go vet` and the race-enabled suite are judged on exit status, as those tools intend. `docs/05-execution-checklist.md` lists what each gate asserts.
See `README.md` and `docs/00-overview.md` for the Saint Hubbins surface.
