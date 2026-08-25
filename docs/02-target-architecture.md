# Target Architecture — Saint Hubbins (Go)

## Module

```
module codeberg.org/uzu/saint-hubbins
go 1.25
```

## Directory Layout

```
Go/
  go.mod
  cmd/
    saint-hubbins/        # native CLI: saint-hubbins eval/serve/render/play/query/midi
    saint-hubbins-wasm/   # WASM entry: exports saintHubbins.queryPattern
  internal/
    core/             # Fraction, TimeSpan, Hap, State, Pattern, controls, scheduler, evaluate
    mini/             # mini notation parser (pigeon-PEG)
    transpiler/       # string transform + goja bridge
    audio/            # webaudio/offline renderer
    draw/             # pianoroll, spiral, pitch wheel
    session/          # live session (evaluation + scheduler)
    tonal/ xen/ edo/ codemirror/ io/ osc/ serial/ mqtt/ gamepad/ motion/ hydra/ csound/ ...
  web/
    server.go         # http.ServeMux, console handlers, static, API, inline console template
    static/saint-hubbins.wasm + wasm_exec.js
  tools/gen-controls  # generates controls_gen.go (checked-in)
```

## Build Targets

- `./scripts/check.sh` — every automated gate in one command; run this before
  claiming work is done
- `go test ./... -race -count=1`
- `go vet ./...`
- `GOOS=js GOARCH=wasm go build -o web/static/saint-hubbins.wasm ./cmd/saint-hubbins-wasm`
- `go run ./cmd/saint-hubbins serve` → console at `http://localhost:8080`

Legacy Strudel JS layout is archived at `docs/archive/strudel-legacy/`.
