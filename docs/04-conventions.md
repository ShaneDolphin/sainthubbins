# Conventions — Saint Hubbins

The naming decisions the rebrand settled, plus the few working rules that
are easy to get wrong here.

- **License**: AGPL-3.0-or-later, copyright `Saint Hubbins contributors`.
- **Attribution**: single `ATTRIBUTION.md` for upstream inspiration; no per-file `Ported from` headers.
- **Module**: `codeberg.org/uzu/saint-hubbins`.
- **CLI**: `saint-hubbins` (alias `hubbins`).
- **Console**: user-facing term is "live console" / "console". Package `internal/session` holds scheduler state.
- **WASM**: artifact `web/static/saint-hubbins.wasm`, global `saintHubbins`.
- **Editor**: `internal/codemirror` theme `hubbinsTheme`.
- **Formatting**: format the files you actually edit — not the tree.
  `gofmt -l .` reports 850 unclean files (20 non-test files under
  `internal/core`, the rest tests), so `make fmt` buries a real change under
  hundreds of cosmetic ones. The target stays for anyone who deliberately
  wants it; just don't reach for it to tidy a diff.
- **Dependencies**: `github.com/dop251/goja` is the only direct dependency.
  A single self-contained Go binary is the point — the OSC and MIDI encoders
  are hand-written rather than imported. Adding one is a deliberate trade.
- **Generated files**: `internal/mini/parser.go` comes from `parser.peg`
  (`pigeon -o internal/mini/parser.go internal/mini/parser.peg`); regenerate
  it whenever the `.peg` changes, or the two drift apart silently.
  `internal/core/controls_gen.go` cannot be regenerated here.
- **Query purity**: a combinator that reads the cycle number off the query
  span must call `.SplitQueries()` or iterate `SpanCycles()` itself. Without
  it a multi-cycle query applies one cycle's decision across the whole span —
  invisible to cycle-by-cycle tests, audible in every rendered track.
- **Checks**: `./scripts/check.sh` (vet, race tests, wasm, CLI content gates).

Legacy conventions that referenced Strudel are archived at `docs/archive/strudel-legacy/04-conventions.md`.
