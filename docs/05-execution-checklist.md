# Execution Checklist — Saint Hubbins

- [ ] `go test ./... -race -count=1` passes
- [ ] `go vet ./...` clean
- [ ] `make wasm` builds `web/static/saint-hubbins.wasm`
- [ ] `go run ./cmd/saint-hubbins serve` boots console at `http://localhost:8080`
- [ ] `saint-hubbins eval 's("bd sd")'` returns haps
- [ ] `saint-hubbins render /tmp/out.wav 's("bd sd")'` writes WAV
- [ ] No `strudel`/`REPL` in user-visible surfaces outside `ATTRIBUTION.md` / `docs/archive/`

Legacy checklist archived at `docs/archive/strudel-legacy/05-execution-checklist.md`.
