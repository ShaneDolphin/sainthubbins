# Roadmap — Saint Hubbins

## Phases (Saint Hubbins)

| Phase | Focus | Key Gate |
|---|---|---|
| 0 | Foundations (Fraction/TimeSpan/Hap/State) | `go test` fraction/timespan |
| 1 | Core engine (Pattern, controls, scheduler) | `go vet` + core tests |
| 2 | Mini + Transpiler | `mini.Mini("bd sd")` parsers green |
| 3 | Audio (offline WAV, filter) | `saint-hubbins render` produces WAV |
| 4 | I/O (MIDI/OSC/Serial/MQTT stubs) | mock interfaces construct without hardware |
| 5 | Live console + WASM | `saint-hubbins serve` → console evaluates `s("bd sd")` |
| 6 | Polish + Stonehenge seasoning | review |

Historical roadmap that referenced Strudel JS is archived at `docs/archive/strudel-legacy/03-roadmap.md`.
