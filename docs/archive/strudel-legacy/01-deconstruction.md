# 01 — Complete Deconstruction of Strudel JS

> Source snapshot: `/Users/shanemorris/Documents/strudel/js` (cloned from https://codeberg.org/uzu/strudel, 2026-08-21, pnpm monorepo).
> Total packages: 34 under `packages/*` + `website/` + `src-tauri/` + `examples/` + `test/` + `bench/` + `tools/` + `docs/`.

## 1. Monorepo Skeleton

| Path | Purpose |
|------|---------|
| `package.json` (root) | `@strudel/monorepo` v0.5.0, scripts `dev/build/test/lint`, private |
| `pnpm-workspace.yaml` | workspaces `packages/*`, `examples/*`, `tools/dbpatch`, `website/` |
| `lerna.json` | lerna v8 |
| `vitest.config.mjs` | vitest + `vite-plugin-bundle-audioworklet`, `vitest.setup.mjs` |
| `eslint.config.mjs` / `.prettierrc` | lint + format gates |
| `jsdoc/jsdoc.config.json` | doc generation → `doc.json` |

## 2. Package Inventory (34 packages)

> LOC counts are `wc -l` over `*.mjs` in each package. Dependency list truncated.

| Package | npm name | LOC | Core files | Key deps | What it does |
|---------|----------|-----|------------|----------|--------------|
| **core** | `@strudel/core` | 11420 | `pattern.mjs` 4191, `controls.mjs` 3319, `signal.mjs` 1099, `repl.mjs` 579, `util.mjs` 508, `euclid.mjs`, `hap.mjs` 178, `fraction.mjs` 147, `timespan.mjs` 117, `cyclist.mjs`, `zyklus.mjs`, `state.mjs` | `fraction.js` | **Heart**. Pattern type, Hap/TimeSpan/State/Fraction, 295 control params, scheduler (zyklus/clockworker), evaluate/evalScope, util, signal ops. All other packages depend on it. |
| **mini** | `@strudel/mini` | ~600 + 2497 gen | `mini.mjs`, `krill-parser.js` (gen), `krill.pegjs` | `@strudel/core` | Mini-notation parser (Tidal mini-language). PEG grammar → AST → Pattern. |
| **transpiler** | `@strudel/transpiler` | ~500 | `transpiler.mjs` 348, `plugin-mini.mjs`, `plugin-sample.mjs`, `plugin-kabelsalat.mjs`, `plugin-widgets.mjs` | `acorn`, `escodegen`, `estree-walker`, `@strudel/core`, `@strudel/mini` | JS transpiler: parses user code with acorn, rewrites mini template literals, sample strings, widgets, emits locations. Registers languages via `registerLanguage`. |
| **webaudio** | `@strudel/webaudio` | ~400 | `webaudio.mjs`, `supradough.mjs`, `scope.mjs`, `spectrum.mjs` | `@strudel/core`, `superdough`, `supradough` | Bridge core → audio: `webaudioOutput`, `renderPatternAudio` (OfflineAudioContext), worklet registration. |
| **superdough** | `superdough` | ~6000 | `superdough.mjs` 1054, `worklets.mjs` 1572, `helpers.mjs` 675, `sampler.mjs`, `synth.mjs` 567, `audioContext.mjs`, `reverb.mjs`, `vowel.mjs`, `wavetable.mjs`, `zzfx*.mjs` | `@kabelsalat/lib`, `@kabelsalat/web`, `nanostores` | Sample-accurate DSP engine: synths, samplers, effects (filter, delay, reverb, distortion, LFO, envelope), wavetables, worklets, node pools. |
| **supradough** | `supradough` | ~1200 | `dough.mjs` 1119, `dough-worklet.mjs`, `dough-export.mjs` | (none) | Alternate dough impl / worklet bundle. |
| **draw** | `@strudel/draw` | ~700 | `pianoroll.mjs` 318, `spiral.mjs`, `draw.mjs`, `animate.mjs`, `color.mjs`, `pitchwheel.mjs` | `@strudel/core` | Visualizers: pianoroll, spiral, pitchwheel, generic draw. |
| **codemirror** | `@strudel/codemirror` | 2439 | `codemirror.mjs` 505, `autocomplete.mjs` 464, `widget.mjs`, `highlight.mjs`, `keybindings.mjs`, `themes/`, `flash.mjs` | `@codemirror/*` (6 pkgs) | Editor integration: CodeMirror 6 setup, autocomplete, widgets (slider), block utilities, flash, themes. |
| **repl** | `@strudel/repl` | ~500 | `repl-component.mjs`, `prebake.mjs` | all core+draw+tonal+etc | REPL component tying editor + evaluation + audio together. |
| **tonal** | `@strudel/tonal` | 1368 | `tonal.mjs` 329, `ireal.mjs` 523, `tonleiter.mjs`, `voicings.mjs` | `@tonaljs/tonal`, `chord-voicings`, `webmidi` | Music theory: scales, chords, voicings, iReal parsing, tonleiter. |
| **edo** | `@strudel/edo` | 628 | `edo.mjs`, `edoscale.mjs`, `pitches.mjs`, `intervals.mjs`, `ratios.mjs` | `@tonaljs/tonal`, `chord-voicings`, `webmidi` | Equal-division-of-octave microtonal support. |
| **xen** | `@strudel/xen` | 293 | `xen.mjs`, `tune.mjs`, `tunejs.js` | `@strudel/core` | Xenharmonic tunings. |
| **midi** | `@strudel/midi` | 923 | `midi.mjs` 683, `input.mjs`, `util.mjs` | `webmidi` | MIDI output/input via WebMIDI. |
| **osc** | `@strudel/osc` | ~250 | `osc.mjs`, `server.js`, `superdirtoutput.js`, `tidal-sniffer.js` | `osc`, `ws` | OSC output to SuperDirt/SuperCollider, server + sniffer. |
| **serial** | `@strudel/serial` | 116 | `serial.mjs` | (core) | Web Serial API. |
| **mqtt** | `@strudel/mqtt` | 129 | `mqtt.mjs` | `paho-mqtt` | MQTT messaging. |
| **motion** | `@strudel/motion` | 389 | `motion.mjs` 386 | (core) | Device motion/orientation → patterns. |
| **gamepad** | `@strudel/gamepad` | 253 | `gamepad.mjs` | (core) | Gamepad API → patterns. |
| **hydra** | `@strudel/hydra` | 51 | `hydra.mjs` | `hydra-synth`, `@strudel/draw` | Hydra visuals bridge. |
| **csound** | `@strudel/csound` | 175 | `index.mjs`, `livecode.orc`, `presets.orc`, `project.csd` | `@csound/browser` | Csound integration. |
| **soundfonts** | `@strudel/soundfonts` | 4055 | `list.mjs` 2028, `gm.mjs` 1787, `sfumato.mjs`, `fontloader.mjs` | `sfumato`, `soundfont2` | SoundFont loading/mapping. |
| **sampler** | `@strudel/sampler` | small | `sample-server.mjs` | `cowsay` | Sample server. |
| **dough** | `@strudel/dough` | 75 | `dough.mjs` | `dough-synth` | Dough synth wrapper. |
| **mondo** | `@strudel/mondo` | 487 + test 979 | `mondo.mjs` | `mondolang`, `@strudel/transpiler` | Mondo language bridge. |
| **mondough** | — | small | `mondough.mjs` | — | Mondo + dough glue. |
| **tidal** | `@strudel/tidal` | 73 | `tidal.mjs` | `hs2js` | Tidal Haskell → JS via hs2js. |
| **hs2js** | `hs2js` | — | `src/` (tree-sitter) | `web-tree-sitter` | Haskell parser (tree-sitter-haskell) → JS. |
| **soundfonts** (dup) | — | — | — | — | — |
| **web** | `@strudel/web` | small | `web.mjs` | core+edo+mini+tonal+transpiler+webaudio | Umbrella re-export. |
| **embed** | `@strudel/embed` | small | `embed.js` | — | Iframe embed. |
| **reference** | `@strudel/reference` | small | `index.mjs` | — | Generated reference data (`undocumented.json`). |
| **desktopbridge** | `@strudel/desktopbridge` | small | `midibridge.mjs`, `oscbridge.mjs`, `loggerbridge.mjs` | `@tauri-apps/api` | Tauri desktop bridges. |
| **vite-plugin-bundle-audioworklet** | — | small | `vite-plugin-bundle-audioworklet.js` | — | Vite plugin for AudioWorklet bundling. |

### Additional top-level areas

| Area | Contents |
|------|----------|
| `website/` | Astro 5 + React 19 + Tailwind 3 + PWA + Supabase + DocSearch + Tauri bindings. `src/repl/` (Repl.jsx, ReplEditor.jsx, panels, piano, audiograph), `src/pages/`, `src/components/`, `src/layouts/`, `public/`. |
| `src-tauri/` | Rust Tauri app (`Cargo.toml`: `tauri` 1.4, `midir` 0.9, `tokio` 1, `rosc` 0.10). Bridges MIDI/OSC/filesystem/clipboard. |
| `examples/` | `minimal-repl`, `buildless`, `codemirror-repl`, `headless-repl`, `superdough`, `tidal-repl`. |
| `test/` | `tunes.test.mjs`, `examples.test.mjs`, `metadata.test.mjs`, `runtime.mjs`, `__snapshots__/`, `testtunes.mjs`. |
| `bench/` | `tunes.bench.mjs` (vitest bench). |
| `tools/dbpatch/` | DB patch tool. |
| `docs/` | `technical-manual/`, `iclc2023-paper/`. |

## 3. Dependency Graph (simplified)

```
core (no strudel deps; only fraction.js)
  ├─ mini, draw, xen, gamepad, serial, motion, csound (partial)
  ├─ tonal, edo  (+ @tonaljs/tonal)
  ├─ transpiler (+ acorn/escodegen/estree-walker) ── mini
  ├─ midi (+ webmidi)
  ├─ webaudio (+ superdough, supradough, draw)
  ├─ superdough (+ @kabelsalat/*, nanostores)  [standalone DSP]
  ├─ soundfonts (+ sfumato)
  └─ repl ── codemirror, draw, hydra, midi, mini, tonal, transpiler, webaudio, edo, ...
       └─ website (Astro/React) ── all
       └─ src-tauri (Rust) ── midir/rosc/tokio
```

## 4. Core Data Model (must be exact in Go)

- **Fraction** (`fraction.mjs` 147 LOC, wraps `fraction.js`): numerator/denominator, exact rational arithmetic, lcm, mod, numeral parsing.
- **TimeSpan** (`timespan.mjs` 117 LOC): `{begin: Fraction, end: Fraction}`, intersection, span arithmetic.
- **Hap** (`hap.mjs` 178 LOC): `{whole: TimeSpan|null, part: TimeSpan, value: any, context: {}, stateful: bool}`, duration/clip logic, whole/part invariant (part ⊆ whole).
- **State** (`state.mjs` 28 LOC): `{span: TimeSpan, controls: {}}` with `setSpan`/`setControls` (immutable copy).
- **Pattern** (`pattern.mjs` 4191 LOC): `class Pattern { query: (State)->Hap[], _steps: Fraction }` + ~200 methods: functor (`withValue`/`fmap`), applicative (`appLeft`/`appRight`), monad (`bind`/`squeezeBind`), time ops (`fast`/`slow`/`early`/`late`/`off`/`ply`), structure (`stack`/`cat`/`fastcat`/`slowcat`/`seq`/`polymeter`), randomness (`rand`/`choose`/`degradeBy`), euclidean, signal integration.
- **Controls** (`controls.mjs` 3319 LOC): 295 `createParam` exports (`note`, `s`, `gain`, `cutoff`, `delay`, …) each returning `{param: value}` bags; alias handling (`isMulti`).
- **Signals** (`signal.mjs` 1099 LOC): continuous patterns (`sine`, `saw`, `tri`, `square`, `rand`, `perlin`).
- **Cyclist/Zyklus** (`cyclist.mjs`, `zyklus.mjs` 54 LOC, `clockworker.js`): clock + scheduler loop, CPS handling.

## 5. What Can Be Stubbed Initially vs Must Be Exact

- **Must be exact from day 1**: Fraction, TimeSpan, Hap, State, Pattern query semantics, mini parser, controls bag semantics.
- **Can be interface-first**: Audio backends (start with `io.Discard` / WAV file), MIDI/OSC (start with logging), draw (start with JSON), website (start with minimal Go `html/template`).

