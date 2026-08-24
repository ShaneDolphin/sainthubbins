# Repository Hygiene Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the roadmap, the completion gates and the project's written conventions describe the software that actually exists, and settle what happens to the unused WASM target.

**Architecture:** Documentation and one CI-style script. No engine code changes. This plan is small and independent — it can run before, after or alongside the others.

**Tech Stack:** Markdown, plus a shell script for the verification gates.

**Spec:** `docs/superpowers/specs/2026-08-24-remaining-work.md`

## Global Constraints

- Every gate written into the checklist must be one a machine can run and a human can read the result of. No gate may pass vacuously.
- Do not delete `docs/archive/` — it holds the Strudel-era history and `ATTRIBUTION.md` depends on that lineage being visible.
- `go vet ./...` and `go test ./... -race -count=1` stay clean.

## File Structure

| File | Responsibility |
|---|---|
| `docs/03-roadmap.md` (rewrite) | Current phases and their real gates. |
| `docs/05-execution-checklist.md` (rewrite) | Runnable gates only. |
| `scripts/check.sh` (create) | Executes every gate; one command, clear pass/fail. |
| `CLAUDE.md` (create) | Conventions a contributor or agent must follow. |
| `docs/tutorial/08-limitations.md` (modify) | WASM entry, per the Task 4 decision. |

---

### Task 1: A roadmap that matches reality

**Files:**
- Rewrite: `docs/03-roadmap.md`

The current file gates phase 5 on the console evaluating `s("bd sd")`, which no evaluator implements. Phases 0–5 are otherwise done and phase 6 is "Polish", which is not a gate.

- [ ] **Step 1: Confirm the current state before writing it down**

```bash
go vet ./...
go test ./... -race -count=1
go run ./cmd/saint-hubbins eval "bd sd"
go run ./cmd/saint-hubbins render /tmp/gate.wav "bd sd" && ls -la /tmp/gate.wav
make wasm && ls -la web/static/saint-hubbins.wasm
```

Record what actually passes. Do not write a gate you have not just run.

- [ ] **Step 2: Rewrite the file**

Replace `docs/03-roadmap.md` with a table of shipped phases marked done, then the remaining tracks, each pointing at its plan in `docs/superpowers/plans/` and stating a gate that can be executed. Suggested remaining tracks:

| Track | Gate | Plan |
|---|---|---|
| Real-time OSC output | `saint-hubbins play "bd*4" 127.0.0.1 57120 4` produces audible SuperDirt output; `go test ./internal/osc/` green | `2026-08-24-realtime-osc-output.md` |
| Pattern engine gaps | `docs/tutorial/08-limitations.md` no longer lists a mini-notation differences table | `2026-08-24-pattern-engine-gaps.md` |
| MIDI output | `saint-hubbins midi /tmp/x.mid "c3 e3 g3"` writes a file a DAW opens | `2026-08-24-midi-output.md` |
| Text evaluator | `saint-hubbins eval 's("bd sd")'` returns two haps with `s` controls | `2026-08-24-text-evaluator.md` |
| Repo hygiene | `./scripts/check.sh` exits 0 | this plan |

Keep the pointer to `docs/archive/strudel-legacy/03-roadmap.md`.

- [ ] **Step 3: Commit**

```bash
git add docs/03-roadmap.md
git commit -m "docs: rewrite the roadmap around gates that can actually be run"
```

---

### Task 2: Runnable gates

**Files:**
- Rewrite: `docs/05-execution-checklist.md`
- Create: `scripts/check.sh`

The checklist contains `saint-hubbins eval 's("bd sd")' returns haps`. It passes today — the command returns one hap whose value is the literal text — which is exactly the failure mode a gate is supposed to catch.

- [ ] **Step 1: Write the script**

```bash
#!/usr/bin/env bash
# Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
# Every gate from docs/05-execution-checklist.md, in one command.
set -uo pipefail

fail=0
check() { # check <name> <command...>
  local name="$1"; shift
  if "$@" >/tmp/sh-check.log 2>&1; then
    printf '  PASS  %s\n' "$name"
  else
    printf '  FAIL  %s\n' "$name"; sed 's/^/        /' /tmp/sh-check.log | head -20
    fail=1
  fi
}

echo "Saint Hubbins — checks"
check "go vet"              go vet ./...
check "go test -race"       go test ./... -race -count=1
check "wasm builds"         make wasm

# Behavioural gates: assert on output, not just on exit status.
check "eval produces haps" bash -c '
  out=$(go run ./cmd/saint-hubbins eval "bd sd")
  echo "$out" | grep -q "\"value\": \"bd\"" || { echo "eval did not return a bd hap: $out"; exit 1; }'

check "render writes audio" bash -c '
  go run ./cmd/saint-hubbins render /tmp/sh-gate.wav "bd sd" >/dev/null
  [ -s /tmp/sh-gate.wav ] || { echo "no wav written"; exit 1; }
  head -c 4 /tmp/sh-gate.wav | grep -q RIFF || { echo "not a RIFF file"; exit 1; }'

check "templates render"   bash -c '
  cd "$(git rev-parse --show-toplevel)" && tmp=$(mktemp -d)
  for d in house chicago-house techno minimal-dubstep maximal-dubstep drum-and-bass electronica trance mytrack; do
    go build -o "$tmp/t" "./examples/$d" || exit 1
    (cd "$tmp" && ./t >/dev/null) || exit 1
  done'

check "no strudel in user surfaces" bash -c '
  hits=$(grep -rniI "strudel" --include="*.go" --include="*.md" . \
    | grep -v "^./ATTRIBUTION.md" | grep -v "^./docs/archive/" \
    | grep -v "^./docs/superpowers/" | grep -v "codeberg.org/uzu/strudel" || true)
  [ -z "$hits" ] || { echo "$hits"; exit 1; }'

exit $fail
```

- [ ] **Step 2: Make it executable and run it**

```bash
chmod +x scripts/check.sh
./scripts/check.sh
```

Expected: every line PASS, exit 0. If "no strudel in user surfaces" fails, read the hits — the tutorial legitimately mentions Strudel when documenting grammar differences, so widen the exclusion to `docs/tutorial/` only if the mentions are attribution or comparison rather than branding.

- [ ] **Step 3: Rewrite the checklist to point at the script**

`docs/05-execution-checklist.md` becomes a short document: run `./scripts/check.sh`, plus the manual gates a script cannot cover (SuperDirt audibility, a DAW opening the MIDI file, the console in a browser).

- [ ] **Step 4: Commit**

```bash
git add scripts/check.sh docs/05-execution-checklist.md
git commit -m "chore: make the execution checklist runnable"
```

---

### Task 3: CLAUDE.md

**Files:**
- Create: `CLAUDE.md`

There is no `CLAUDE.md`. Its absence has already cost: the rule that cycle-dependent combinators must call `SplitQueries` was not written down, and `Every`/`LastOf` shipped without it, silently corrupting every rendered variation.

- [ ] **Step 1: Write the file**

Cover, with a reason for each rule rather than a bare instruction:

- **Module layout.** Engine under `internal/`; anything importing it lives in this module. Example programs go in `examples/`.
- **No new third-party dependencies** without discussion. The selling point is a single Go binary.
- **AGPL-3.0-or-later header** on every new file, matching existing files.
- **Cycle-dependent combinators must call `SplitQueries()`.** Anything reading `state.Span.Begin.Sam()` decides per cycle, so a query spanning several cycles must be split. `pattern_random.go` and `pattern_vlpf_morph.go` are the reference; `Every` and `LastOf` were fixed for omitting it. Include the symptom: correct behaviour per-cycle, wrong behaviour under a wide query, and the renderer only ever issues wide queries.
- **The renderer only queries once, over the whole span.** So a bug that only appears in a wide query reaches the audio and nothing else.
- **Test with both query shapes.** Assert that one N-cycle query agrees with N one-cycle queries.
- **Only `Note`/`N`/`S`, `Gain` and `Cutoff`/`Lpf` reach the offline audio.** The other ~290 controls are data-only; do not "fix" a control that appears to do nothing.
- **`examples/` is teaching material.** A comment that misdescribes its code is a defect, not a nitpick. Verify comments by running the pattern.
- **Documented numbers are asserted.** `docs/tutorial/` quotes exact event counts and peaks; changing behaviour means re-running the templates and updating them.
- **Run `./scripts/check.sh`** before claiming work is complete.

- [ ] **Step 2: Verify the claims in it**

Every technical claim must be one you have run. Confirm the `SplitQueries` list is current:

```bash
grep -rn "state.Span.Begin.Sam()" internal/core/*.go | grep -v _test
```

Each hit must either call `SplitQueries()` or iterate `SpanCycles()` itself. If a new one has appeared, fix it or note it.

- [ ] **Step 3: Commit**

```bash
git add CLAUDE.md
git commit -m "docs: add CLAUDE.md recording engine conventions"
```

---

### Task 4: Decide what happens to WASM

**Files:**
- Modify: `docs/tutorial/08-limitations.md`, `README.md`, and possibly `Makefile`, `web/server.go`, `cmd/saint-hubbins-wasm/`

`make wasm` builds `web/static/saint-hubbins.wasm` and the console footer names it, but the page never loads it — the console calls `/api/evaluate` over HTTP. The target is built, shipped in the docs, and dead.

Pick one:

**A. Wire it up.** The console loads `wasm_exec.js` and the module, calling `saintHubbins.queryPattern` in the browser instead of round-tripping. Removes server latency and makes the console work offline. Cost: real front-end work, and the WASM and HTTP paths must not drift.

**B. Keep building it, stop advertising it.** It stays a supported build target for embedders. Remove the console footer's reference and say plainly in the docs that the console uses HTTP.

**C. Remove it.** Delete `cmd/saint-hubbins-wasm/`, the `wasm` make target and the static files. Least code, closes a door.

Recommendation: **B**. It is honest, costs almost nothing, and keeps the target for anyone embedding the engine in a page. Revisit A if the console grows enough that latency matters.

- [ ] **Step 1: Record the decision**

Write it into `docs/03-roadmap.md` under the remaining tracks, with the reasoning, so it is not re-litigated.

- [ ] **Step 2: Apply option B**

In `web/server.go`, remove the WASM reference from the console footer, since it advertises a bridge the page does not use. Keep the `/static/` route. In `docs/tutorial/08-limitations.md`, rewrite the WASM section from a limitation into a statement of intent: the target is built for embedders, the console deliberately uses HTTP.

- [ ] **Step 3: Verify and commit**

```bash
make wasm && ls -la web/static/saint-hubbins.wasm
go test ./web/ && go run ./cmd/saint-hubbins serve :8099 &
sleep 2 && curl -s http://localhost:8099/ | grep -i wasm || echo "no wasm reference in the page (expected)"
kill %1
git add web/server.go docs/tutorial/08-limitations.md docs/03-roadmap.md
git commit -m "docs: settle the WASM target's role and stop advertising an unused bridge"
```

---

## Verification

```bash
./scripts/check.sh
```

Then read `docs/03-roadmap.md` and `docs/05-execution-checklist.md` as a newcomer would: every gate should be something you can run in the next thirty seconds, and none should pass while the underlying feature is missing.
