# Phase 5 — Live Console, Draw & WASM (Saint Hubbins)

Archived. Original Strudel-era plan moved to `docs/archive/strudel-legacy/phases/phase-5-frontend-repl.md`.

Current console lives at `web/server.go` (`Server.Handler()`), with the page itself an inline Go template (`consoleTemplate`) in that same file. The WASM entry point is `cmd/saint-hubbins-wasm` (`saintHubbins.queryPattern`); the console does not load it — see `docs/03-roadmap.md`.
