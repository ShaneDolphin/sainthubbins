# 05 — Execution Checklist (print and tick)

## Before Phase 0

- [ ] Read `00-overview.md`, `01-deconstruction.md`, `02-target-architecture.md`, `03-roadmap.md`
- [ ] `cp js/LICENSE Go/LICENSE` and verify AGPL-3.0 header template
- [ ] `go mod init` + `Makefile` + `.golangci.yml` from Phase 0 doc
- [ ] Create `internal/core/testdata/` and script to dump JS fixtures (`node dump-fixtures.mjs`)

## Phase Gates (copy to issue tracker)

- [ ] **M0.1** `internal/core` Fraction/TimeSpan/Hap/State green — `go test ./internal/core` passes
- [ ] **M1.1** Pattern core (fmap/bind/cat/stack) green
- [ ] **M1.2** Controls generated (295) + `TestControls` green
- [ ] **M1.3** Scheduler + evaluate green
- [ ] **M2.1** Mini parser green
- [ ] **M2.2** Transpiler (shim or pure) green — `TestTunes` passes
- [ ] **M3.1** AudioContext + sampler/synth green
- [ ] **M3.2** Offline WAV render green — `strudel render` produces valid WAV
- [ ] **M4.1** MIDI + OSC green
- [ ] **M5.1** `make wasm` + `web/server.go` green — REPL serves at :8080
- [ ] **M5.2** End-to-end demo + tag `v0.1.0-go`

## Final Verification

```bash
go vet ./...
golangci-lint run
go test ./... -race -count=1
GOOS=js GOARCH=wasm go build -o web/static/strudel.wasm ./cmd/strudel-wasm
go run ./cmd/strudel serve &  # verify http://localhost:8080
curl -s -X POST http://localhost:8080/api/evaluate -d '{"code":"s(\"bd sd\")"}' | jq .
go run ./cmd/strudel render --out /tmp/test.wav 's("bd sd")' && file /tmp/test.wav
```

