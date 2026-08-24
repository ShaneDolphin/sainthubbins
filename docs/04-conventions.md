# Conventions — Saint Hubbins

- **License**: AGPL-3.0-or-later, copyright `Saint Hubbins contributors`.
- **Attribution**: single `ATTRIBUTION.md` for upstream inspiration; no per-file `Ported from` headers.
- **Module**: `codeberg.org/uzu/saint-hubbins`.
- **CLI**: `saint-hubbins` (alias `hubbins`).
- **Console**: user-facing term is "live console" / "console". Package `internal/session` holds scheduler state.
- **WASM**: artifact `web/static/saint-hubbins.wasm`, global `saintHubbins`.
- **Editor**: `internal/codemirror` theme `hubbinsTheme`.
- **Formatting**: `gofmt -w .`, `go vet ./...`.

Legacy conventions that referenced Strudel are archived at `docs/archive/strudel-legacy/04-conventions.md`.
