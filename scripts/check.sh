#!/usr/bin/env bash
# Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
#
# Every gate from docs/05-execution-checklist.md and docs/03-roadmap.md, in
# one command. Each gate asserts on *content*, not just an exit code — a
# command that runs successfully and produces the wrong (or literal,
# unevaluated) output must still FAIL here. That distinction is the entire
# point of this script: docs/05-execution-checklist.md used to contain
# `saint-hubbins eval 's("bd sd")' returns haps`, which passed even though no
# evaluator for that syntax exists — the command "succeeded" by returning one
# hap whose value was the literal source text. A human ticking a checkbox
# had no way to notice. A script that greps the output for a real hap value
# does.
#
# set -uo pipefail (deliberately not -e): every gate runs and reports,
# rather than stopping at the first failure. You want to see everything
# that's broken in one pass, not fix-rerun-discover-fix-rerun one at a time.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT

fail=0
check() { # check <name> <command...>
  local name="$1"; shift
  if "$@" >"$workdir/log" 2>&1; then
    printf '  PASS  %s\n' "$name"
  else
    printf '  FAIL  %s\n' "$name"
    sed 's/^/        /' "$workdir/log" | head -40
    fail=1
  fi
}

echo "Saint Hubbins — checks"
echo

check "go vet"        go vet ./...
check "go test -race" go test ./... -race -count=1

# make wasm: the Makefile target ends in
#   ... || echo "wasm build skipped (no wasm target)"
# which means a *failed* GOOS=js GOARCH=wasm build still exits 0 (the `||`
# swallows the error and prints a message instead). `make wasm` returning
# success proves nothing by itself — only inspecting the resulting binary
# does. A real WASM module starts with the 4-byte magic number \0asm.
check "wasm builds" bash -c '
  make wasm || exit 1
  f="web/static/saint-hubbins.wasm"
  [ -s "$f" ] || { echo "$f missing or empty"; exit 1; }
  magic=$(dd if="$f" bs=1 count=4 2>/dev/null)
  [ "$magic" = "$(printf "\0asm")" ] || { echo "$f does not start with the wasm magic number (\\0asm)"; exit 1; }
'

# eval: assert the JSON contains real haps carrying "bd"/"sd" values, not a
# single hap whose value is the literal, unevaluated source string. This is
# the gate that used to pass vacuously (see header comment above). Bare
# mini-notation ("bd sd") is what actually evaluates today — see
# docs/03-roadmap.md "Text evaluator" for the separately tracked gap where
# s("bd sd") function-call syntax comes back as one literal-string hap.
check "eval produces real haps" bash -c '
  out=$(go run ./cmd/saint-hubbins eval "bd sd")
  echo "$out" | grep -q "\"value\": \"bd\"" || { echo "no bd hap in output:"; echo "$out"; exit 1; }
  echo "$out" | grep -q "\"value\": \"sd\"" || { echo "no sd hap in output:"; echo "$out"; exit 1; }
  echo "$out" | grep -q "^2 haps$" || { echo "expected exactly 2 haps, got:"; echo "$out"; exit 1; }
'

# render: check the RIFF/WAVE header AND that the sample data is not all
# zero. A silent WAV with a valid header would pass a header-only check
# while the actual pattern engine was producing nothing audible.
check "render writes genuine audio" bash -c '
  out="'"$workdir"'/gate.wav"
  go run ./cmd/saint-hubbins render "$out" "bd sd" >/dev/null
  [ -s "$out" ] || { echo "no wav written"; exit 1; }
  riff=$(dd if="$out" bs=1 skip=0 count=4 2>/dev/null)
  wave=$(dd if="$out" bs=1 skip=8 count=4 2>/dev/null)
  [ "$riff" = "RIFF" ] || { echo "missing RIFF tag"; exit 1; }
  [ "$wave" = "WAVE" ] || { echo "missing WAVE tag"; exit 1; }
  nonzero=$(tail -c +45 "$out" | od -An -tx1 | tr -d " \n0")
  [ -n "$nonzero" ] || { echo "wav header is valid but the audio data is all zero (silence)"; exit 1; }
'

# midi: check the MThd header AND that the CLI itself reports a nonzero note
# count. runMIDI in cmd/saint-hubbins/main.go prints
# "wrote <path> (<cycles> cycles, <N> notes)" — N=0 means a file with a
# valid header but nothing in it (and triggers a separate stderr warning we
# do not depend on here; the note count in the primary report line is the
# actual assertion).
check "midi writes a genuine Standard MIDI File" bash -c '
  out="'"$workdir"'/gate.mid"
  report=$(go run ./cmd/saint-hubbins midi "$out" "bd sd" 1)
  [ -s "$out" ] || { echo "no mid written"; exit 1; }
  hdr=$(dd if="$out" bs=1 skip=0 count=4 2>/dev/null)
  [ "$hdr" = "MThd" ] || { echo "missing MThd tag"; exit 1; }
  echo "$report" | grep -Eq "[1-9][0-9]* notes\)" || { echo "CLI reported zero notes: $report"; exit 1; }
'

# play: cmd/saint-hubbins/play_test.go's TestRunPlaySendsToSuperDirt already
# opens a real UDP listener and asserts the OSC bytes that arrive carry the
# sound name ("bd") — real content over a real socket. Run it rather than
# reimplementing a weaker version here.
check "play reaches a socket (OSC to SuperDirt)" \
  go test ./cmd/saint-hubbins/... -run TestRunPlaySendsToSuperDirt -v -count=1

# templates: examples/examples_test.go's TestTemplatesRenderAudio builds and
# runs all nine tutorial templates and asserts each WAV has a nonzero peak
# (not silent), does not clip, and has the expected frame count — stronger
# than "the binary exited 0 and a file exists". Run it explicitly so a
# failure here is attributed to the templates, not buried in the full
# `go test ./...` gate above.
check "all nine example templates build and run" \
  go test ./examples/... -run TestTemplatesRenderAudio -v -count=1

# Rebrand check: NOT a bare grep for "strudel" or "REPL".
#
# Bare "strudel" flags legitimate content that should stay: comparisons in
# the tutorial ("...from Strudel or TidalCycles, it says so"), a technical
# note in internal/io/hapmidi.go about Strudel-canonical naming, and this
# script/checklist's own description of what it checks. Bare "REPL" is pure
# noise outside the archive — every hit is Replicate/ReplaceAll/replacing,
# not a real REPL reference. Neither proves anything about whether the
# rebrand is complete.
#
# What actually indicates an incomplete rebrand is a *dependency or import
# path* pointing at the old project: an npm package (@strudel/*), the old
# Go import/repo name (strudel-go), or the old module path
# (codeberg.org/uzu/strudel — the module is codeberg.org/uzu/saint-hubbins
# now). Those three return zero hits today; if one comes back, it means
# someone reintroduced a real Strudel dependency, not just a comparison or
# attribution mention. Do not "simplify" this back to a bare word search —
# that was tried and it only produces false positives/negatives.
#
# Exclusions: ATTRIBUTION.md (deliberately documents the upstream project),
# docs/archive/ (frozen legacy history), docs/superpowers/ and .superpowers/
# (planning artifacts about this very rebrand, which necessarily discuss
# these markers as text), internal/mini/krill.peg's header (a copyright
# attribution for a ported grammar file, same category as ATTRIBUTION.md),
# and this script plus the execution checklist it backs — both *document*
# the markers they check for, in comments/prose, so they legitimately
# contain the literal substrings without being a dependency on them.
check "no leftover Strudel dependency/import path" bash -c '
  hits=$(grep -rnI \
      --exclude-dir=.git --exclude-dir=.superpowers --exclude-dir=archive --exclude-dir=superpowers \
      -e "@strudel/" -e "strudel-go" -e "codeberg\.org/uzu/strudel" . \
    | grep -v "^\./ATTRIBUTION.md:" \
    | grep -v "^\./internal/mini/krill.peg:" \
    | grep -v "^\./scripts/check.sh:" \
    | grep -v "^\./docs/05-execution-checklist.md:" \
    || true)
  [ -z "$hits" ] || { echo "$hits"; exit 1; }
'

echo
if [ "$fail" -eq 0 ]; then
  echo "All checks passed."
else
  echo "One or more checks FAILED — see output above."
fi
exit $fail
