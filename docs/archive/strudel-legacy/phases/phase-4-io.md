# Phase 4 — I/O Integrations (MIDI / OSC / Serial / MQTT / Gamepad / Motion / Hydra / CSound)

> Duration: 2-3 weeks. Depends on Phase 1. Can overlap with Phase 3.

## Objective

Port all **output/input bridges** so patterns can drive external hardware/software. Each bridge is an `io` package with a common `Output` interface.

## Common Interface

```go
// internal/io/output.go
type Output interface {
    Name() string
    Trigger(hap core.Hap, deadline float64, duration float64, cps float64) error
    Close() error
}
// Registry mirrors superdough's soundMap
var Outputs = map[string]Output{}
func RegisterOutput(o Output)
```

## 4.1 MIDI (`packages/midi/` 923 LOC: `midi.mjs` 683, `input.mjs`, `util.mjs`)

- **JS**: WebMIDI (`navigator.requestMIDIAccess`).
- **Go native**: `gitlab.com/gomidi/midi/v2` + `gitlab.com/gomidi/rtmididrv` (or `github.com/rakyll/portmidi`).
- **Go WASM**: `syscall/js` bridge to WebMIDI (same as JS).
- **Mock** for tests: records `NoteOn`/`NoteOff`/`CC` messages to slice.
- Files: `internal/io/midi/midi.go`, `input.go`, `util.go`, `midi_mock.go`, `midi_native.go`, `midi_wasm.go`.

## 4.2 OSC (`packages/osc/` ~250 LOC: `osc.mjs`, `server.js`, `superdirtoutput.js`, `tidal-sniffer.js`)

- **JS**: `osc` + `ws` packages, sends to SuperDirt.
- **Go**: `github.com/hypebeast/go-osc` (client + server), `nhooyr.io/websocket` for WS.
- Implement `SuperDirtOutput` (OSC bundle for SuperDirt) + `TidalSniffer` (listens for Tidal OSC).
- Server: `internal/io/osc/server.go` — `go run ./cmd/strudel osc --port 57120`.

## 4.3 Serial (`serial.mjs` 116 LOC)

- **JS**: Web Serial API.
- **Go native**: `github.com/bugst/go-serial`.
- **Go WASM**: `syscall/js` bridge.
- Mock: writes to `bytes.Buffer`.

## 4.4 MQTT (`mqtt.mjs` 129 LOC)

- **JS**: `paho-mqtt`.
- **Go**: `github.com/eclipse/paho.mqtt.golang`.
- Config: broker URL, topic prefix.

## 4.5 Gamepad (`gamepad.mjs` 253 LOC)

- **JS**: `navigator.getGamepads()`.
- **Go native**: `github.com/0xcafed00d/gamepad` or `github.com/simulatedgreg/electron-gamepad` equivalent; poll via goroutine.
- **Go WASM**: `syscall/js` bridge.
- As pattern source: `Gamepad().Button(0)` → `Pattern`.

## 4.6 Motion (`motion.mjs` 386 LOC)

- **JS**: `DeviceMotionEvent` / `DeviceOrientationEvent`.
- **Go native**: not applicable (no device motion on server); provide mock + WASM bridge.
- WASM: `js.Global().Get("DeviceMotionEvent")`.

## 4.7 Hydra (`hydra.mjs` 51 LOC)

- **JS**: `hydra-synth` visuals.
- **Go**: visuals are JS-only; Go `hydra` package is a bridge that sends Hap values to Hydra via WASM `js` call. Native: no-op or log.

## 4.8 CSound (`csound/ 175 LOC`)

- **JS**: `@csound/browser`.
- **Go**: `github.com/future-architect/csound-wasm` or `csound` CGO binding (optional). Start as WASM-only bridge; native can be `//go:build ignore` initially.

## 4.9 Desktop Bridge (`desktopbridge/` — Tauri bridges)

- **JS**: `@tauri-apps/api` bridges for MIDI/OSC/logger.
- **Go**: if using `wails`, these become `wails` bindings directly; no separate `desktopbridge` package needed. Document mapping: `midibridge.mjs` → `internal/io/midi` native, `oscbridge.mjs` → `internal/io/osc` native, `loggerbridge.mjs` → `log/slog`.

## Testing

- Each `internal/io/<name>/` has `*_test.go` using mock output that records triggers; `TestMIDIOutput` asserts correct MIDI bytes for `note("c3 e3")`.
- Integration: `TestOSCBundle` verifies SuperDirt bundle format against JS fixture.

## Acceptance Checklist

- [ ] `Output` interface defined; registry works
- [ ] Each I/O package compiles for both `!js` (native) and `js && wasm` (WASM) with mocks
- [ ] MIDI mock records correct NoteOn/Off for known pattern; `TestMIDI` green
- [ ] OSC client sends valid bundle; `TestOSC` green with local UDP listener
- [ ] Serial/MQTT mocks green
- [ ] Gamepad/Motion WASM bridges compile (no native test required beyond mock)
- [ ] No I/O package imports JS-only deps in native build

