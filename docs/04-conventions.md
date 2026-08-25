# Conventions — Saint Hubbins

`CLAUDE.md` at the repository root is authoritative for how to work in this
codebase — formatting, dependencies, generated files, the rules that protect
query purity. This file records the *naming* decisions the rebrand settled,
and nothing else; restating engine rules here only gives them somewhere to
drift out of sync.

- **License**: AGPL-3.0-or-later, copyright `Saint Hubbins contributors`.
- **Attribution**: single `ATTRIBUTION.md` for upstream inspiration; no per-file `Ported from` headers.
- **Module**: `codeberg.org/uzu/saint-hubbins`.
- **CLI**: `saint-hubbins` (alias `hubbins`).
- **Console**: user-facing term is "live console" / "console". Package `internal/session` holds scheduler state.
- **WASM**: artifact `web/static/saint-hubbins.wasm`, global `saintHubbins`.
- **Editor**: `internal/codemirror` theme `hubbinsTheme`.
- **Formatting**: see `CLAUDE.md` — format the files you touch, not the tree.
- **Checks**: `./scripts/check.sh` (vet, race tests, wasm, CLI content gates).

Legacy conventions that referenced Strudel are archived at `docs/archive/strudel-legacy/04-conventions.md`.
