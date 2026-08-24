# Phase 3 — Audio Engine (superdough / supradough / webaudio / soundfonts / dough)

> Duration: 4-6 weeks. Depends on Phase 1 (Pattern). Can overlap with Phase 2. This is the second-largest phase by LOC.

## Objective

Port the **sample-accurate DSP engine** so patterns can be rendered to audio — offline to WAV (testable) and real-time via native audio or WASM bridge.

## 3.0 Architecture Recap (from 02)

```
internal/audio/
  context.go          # AudioContext interface
  context_native.go   # oto/malgo backend  //go:build !js
  context_wasm.go     # syscall/js bridge  //go:build js && wasm
  superdough.go       # main trigger: superdough(hapValue, time, duration, cps)
  sampler.go          # getSampleBufferSource, loadBuffer
  synth.go            # synths (zzfx, wavetable, etc.)
  effects.go          # filter, delay, reverb, distortion, LFO, envelope, vowel
  helpers.go          # createFilter, gainNode, getCompressor, etc.
  wavetable.go        # wavetable.mjs
  worklets.go         # worklet registration + dspworklet.mjs
  nodepool.go         # nodePools.mjs
  webaudio.go         # webaudioOutput, renderPatternAudio, OfflineAudioContext path
  vorticon/*          # extra DSP if needed
internal/soundfonts/
  loader.go           # fontloader.mjs + sfumato.mjs
  gm.go               # gm.mjs (generated, 1787 LOC)
  list.go             # list.mjs (2028 LOC)
```

## 3.1 AudioContext Abstraction

```go
// context.go
type AudioContext interface {
    CurrentTime() float64
    SampleRate() int
    Destination() AudioNode
    CreateGain() GainNode
    CreateBiquadFilter() BiquadFilterNode
    CreateDelay(maxDelay float64) DelayNode
    CreateBufferSource() BufferSourceNode
    CreateOscillator() OscillatorNode
    CreateWaveShaper() WaveShaperNode
    CreateConvolver() ConvolverNode
    CreateAnalyser() AnalyserNode
    DecodeAudioData(data []byte) (AudioBuffer, error)
}
type AudioNode interface { Connect(AudioNode); Disconnect() }
```

- **Native** (`context_native.go`): implement via `github.com/ebitengine/oto/v3` for output + manual DSP for effects. Start with "null" context that renders to `[]float32` buffer (offline) — sufficient for tests without hardware.
- **WASM** (`context_wasm.go`): `js.Global().Get("AudioContext").New()` via `syscall/js`.

## 3.2 Sampler (`sampler.mjs` 392 LOC)

- `LoadBuffer(url) (AudioBuffer, error)` — fetch + decode (native: `http.Get` + WAV/MP3 decode via `github.com/go-audio/wav` / `ebitengine/oto`).
- `GetSampleBufferSource(sampleName) (BufferSource, error)` — lookup in `soundMap` (nanostores → `sync.Map` in Go).
- Caching: `sync.Map` + LRU if memory-constrained.

## 3.3 Synth (`synth.mjs` 567 LOC + `zzfx*.mjs` + `wavetable.mjs` 353 LOC + `noise.mjs` + `vowel.mjs`)

- Port each synth function to Go DSP operating on `[]float32` at sample rate.
- `zzfx` — port `zzfx_fork.mjs` + `zzfx.mjs` (procedural synth).
- `wavetable` — `resetSeenKeys`, wavetable generation.
- Reference: `@kabelsalat/lib` / `@kabelsalat/web` — if Go DSP diverges, vendor the JS DSP constants and match.

## 3.4 Effects (`helpers.mjs` 675 LOC, `reverb.mjs`, `feedbackdelay.mjs`, `vowel.mjs`, `modulators.mjs`, `dspworklet.mjs`)

- `CreateFilter(type, freq, q)` → biquad (use `github.com/youpy/go-wav` or implement RBJ cookbook).
- `EffectSend`, `GainNode`, `GetCompressor`, `GetDistortion`, `GetLFO`, `GetWorklet`.
- Reverb: `reverb.mjs` + `reverbGen.mjs` — convolution reverb; start with simple Schroeder, add convolution later.
- Modulators: `connectLFO`, `connectEnvelope`, `connectBusModulator`.

## 3.5 Superdough Main (`superdough.mjs` 1054 LOC)

```go
func Superdough(value map[string]any, t, hapDuration, cps float64, wholeBegin *Fraction) error
func SetMaxPolyphony(n int)
func RegisterSound(key string, trigger func(map[string]any) error, data map[string]any)
```

- `value` is the Hap bag (`s`, `note`, `gain`, `cutoff`, `delay`, `room`, etc.).
- Dispatch: sample vs synth vs external (MIDI/OSC handled via `internal/io`).
- Polyphony: FIFO voice stealing when `maxPolyphony` exceeded.

## 3.6 WebAudio Bridge (`webaudio.mjs`)

```go
func WebaudioOutput(hap core.Hap, deadline float64, hapDuration, cps, t float64) error
func RenderPatternAudio(pattern core.Pattern, cps, begin, end float64, sampleRate, maxPolyphony int, multiChannel bool) ([]byte, error) // returns WAV bytes
```

- `RenderPatternAudio` for offline: create `OfflineAudioContext` (native: buffer render; WASM: `OfflineAudioContext` JS).
- CLI: `go run ./cmd/strudel render --out out.wav "s(\"bd sd\")"`

## 3.7 SoundFonts (`soundfonts/` 4055 LOC)

- `gm.mjs` (1787 LOC) + `list.mjs` (2028 LOC) are **data** — generate Go equivalents via `tools/gen-gm`.
- `fontloader.mjs` + `sfumato.mjs` — SoundFont2 parsing; use `github.com/sin Schott/sfumato` Go port or call `sfumato` via `goja` shim initially.

## 3.8 Dough / Supradough (`dough.mjs` 75 LOC, `supradough/dough.mjs` 1119 LOC)

- Dough is a thin synth wrapper; supradough is the worklet-based dough. Port after superdough is green — share `AudioContext` abstraction.

## Testing

- **Offline render test**: `TestRenderPatternAudio` — render 2 cycles of `s("bd sd")` at 44100 Hz, compare Go WAV vs JS WAV (generated fixture) within 1e-4 RMS.
- **Unit**: each effect has `Test<Effect>` that processes a sine buffer and checks output not silent and within expected gain.
- **No hardware required**: all tests use offline/buffer path.

## Acceptance Checklist

- [ ] `AudioContext` interface defined; native null-buffer impl + WASM bridge compile (build tags)
- [ ] Sampler loads and caches sample; `TestSampler` green without network (use embedded fixture WAV)
- [ ] Synth produces audible buffer; `TestSynth` checks RMS > 0
- [ ] `Superdough` triggers without panic for `{s: "bd", gain: 0.5}` bag
- [ ] `RenderPatternAudio` produces valid WAV file; `go run ./cmd/strudel render` works
- [ ] WAV output matches JS offline render within tolerance (or documented delta)
- [ ] `soundfonts/gm.go` generated and `TestSoundfontList` passes

