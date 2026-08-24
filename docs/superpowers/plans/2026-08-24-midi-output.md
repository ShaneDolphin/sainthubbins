# MIDI Output Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Turn the MIDI mock into something useful — export a pattern as a Standard MIDI File any DAW can open, and make live MIDI a one-file addition once a device backend is chosen.

**Architecture:** `internal/io` already defines `MIDIInterface` (`SendNoteOn`/`SendNoteOff`/`SendCC`/`Close`) with `MockMIDI` behind it. This plan adds two real implementations of the same interface: `FileMIDI`, which accumulates events and writes a Standard MIDI File, and `SinkMIDI`, which adapts `MIDIInterface` to the `session.Sink` from the real-time OSC plan so a live pattern can drive MIDI as soon as a device backend exists.

**Tech Stack:** Go 1.25 standard library only. SMF is a simple binary container — variable-length quantities, a header chunk and a track chunk — and writing it by hand avoids a dependency.

**Spec:** `docs/superpowers/specs/2026-08-24-remaining-work.md`

**Depends on:** Task 5 of `2026-08-24-realtime-osc-output.md` for the `session.Sink` interface. Tasks 1–3 here do not need it; Task 4 does.

## Global Constraints

- Go 1.25.0, module `codeberg.org/uzu/saint-hubbins`. No new dependencies.
- AGPL-3.0-or-later header on every new file.
- Do not change `MIDIInterface`, `MockMIDI` or `MIDIFromHaps` signatures — `internal/io`'s existing tests exercise them.
- Tests must be hermetic: no MIDI hardware, no virtual ports.
- `go test ./... -race -count=1` and `go vet ./...` must stay clean.

## Decision required before Task 5

Real-time output to a MIDI device cannot be done from the Go standard library. The realistic options, in order of cost:

1. **Do nothing more.** Ship SMF export (Tasks 1–3). Users import the file into a DAW. No dependency, no platform code. This plan assumes this is where you stop.
2. **Add `gitlab.com/gomidi/midi/v2`** with its `rtmididrv` backend. Pure-Go API, but the driver needs cgo and platform MIDI headers. Breaks "single binary, no runtime deps" on some platforms.
3. **Write a CoreMIDI/ALSA/WinMM binding by hand** behind a build tag. No module dependency, but a meaningful amount of per-platform cgo to maintain.

Task 5 is written for option 2 and is **deliberately left unstarted**. Do not begin it without an explicit decision, because it changes the project's dependency story.

## File Structure

| File | Responsibility |
|---|---|
| `internal/io/smf.go` (create) | Standard MIDI File encoding: VLQ, header and track chunks. |
| `internal/io/smf_test.go` (create) | Byte-level tests against the SMF spec. |
| `internal/io/filemidi.go` (create) | `FileMIDI` — a `MIDIInterface` that records timed events and writes a `.mid`. |
| `internal/io/filemidi_test.go` (create) | Round-trip tests. |
| `internal/io/hapmidi.go` (create) | `Hap` → note number, velocity and channel. |
| `internal/io/hapmidi_test.go` (create) | Conversion tests. |
| `cmd/saint-hubbins/main.go` (modify) | `midi` subcommand: render a pattern to a `.mid`. |

---

### Task 1: Standard MIDI File encoding

**Files:**
- Create: `internal/io/smf.go`
- Test: `internal/io/smf_test.go`

**Interfaces:**
- Produces: `func writeVLQ(v uint32) []byte`, `func EncodeSMF(ticksPerQuarter int, events []TimedEvent) []byte`, and `type TimedEvent struct { Tick uint32; Data []byte }`.

Background: an SMF is `MThd` (header: format, track count, division) followed by `MTrk` chunks. Inside a track, each event is preceded by a delta time encoded as a variable-length quantity — seven bits per byte, high bit set on all but the last.

- [ ] **Step 1: Write the failing test**

```go
// internal/io/smf_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"bytes"
	"testing"
)

func TestWriteVLQ(t *testing.T) {
	cases := map[uint32][]byte{
		0:       {0x00},
		127:     {0x7F},
		128:     {0x81, 0x00},
		8192:    {0xC0, 0x00},
		1048576: {0xC0, 0x80, 0x00},
	}
	for in, want := range cases {
		if got := writeVLQ(in); !bytes.Equal(got, want) {
			t.Errorf("writeVLQ(%d) = % X, want % X", in, got, want)
		}
	}
}

func TestEncodeSMFStructure(t *testing.T) {
	out := EncodeSMF(480, []TimedEvent{
		{Tick: 0, Data: []byte{0x90, 60, 100}},   // note on
		{Tick: 480, Data: []byte{0x80, 60, 0}},   // note off a quarter later
	})
	if !bytes.HasPrefix(out, []byte("MThd")) {
		t.Fatalf("missing MThd header: % X", out[:4])
	}
	if !bytes.Contains(out, []byte("MTrk")) {
		t.Fatal("missing MTrk chunk")
	}
	// End-of-track meta event is mandatory.
	if !bytes.Contains(out, []byte{0xFF, 0x2F, 0x00}) {
		t.Error("missing end-of-track meta event")
	}
}

func TestEncodeSMFUsesDeltaTimes(t *testing.T) {
	// Two events at the same tick must produce a zero delta for the second.
	out := EncodeSMF(480, []TimedEvent{
		{Tick: 100, Data: []byte{0x90, 60, 100}},
		{Tick: 100, Data: []byte{0x90, 64, 100}},
	})
	// The second event's delta is 0x00 immediately before its status byte.
	if !bytes.Contains(out, []byte{0x00, 0x90, 64, 100}) {
		t.Errorf("second simultaneous event should have a zero delta: % X", out)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/io/ -run "VLQ|SMF" -v`
Expected: FAIL — `undefined: writeVLQ`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/io/smf.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Standard MIDI File encoding — no dependencies, just the container format.

package io

import "encoding/binary"

// TimedEvent is a raw MIDI message at an absolute tick.
type TimedEvent struct {
	Tick uint32
	Data []byte
}

// writeVLQ encodes a variable-length quantity: seven bits per byte, with the
// high bit set on every byte except the last.
func writeVLQ(v uint32) []byte {
	if v == 0 {
		return []byte{0}
	}
	var stack []byte
	for v > 0 {
		stack = append(stack, byte(v&0x7F))
		v >>= 7
	}
	out := make([]byte, 0, len(stack))
	for i := len(stack) - 1; i >= 0; i-- {
		b := stack[i]
		if i != 0 {
			b |= 0x80
		}
		out = append(out, b)
	}
	return out
}

// EncodeSMF writes a format-0 single-track file. Events must be sorted by tick;
// deltas are derived from the gaps between them.
func EncodeSMF(ticksPerQuarter int, events []TimedEvent) []byte {
	var track []byte
	var last uint32
	for _, e := range events {
		delta := e.Tick - last
		last = e.Tick
		track = append(track, writeVLQ(delta)...)
		track = append(track, e.Data...)
	}
	// End of track.
	track = append(track, 0x00, 0xFF, 0x2F, 0x00)

	out := make([]byte, 0, len(track)+22)
	out = append(out, "MThd"...)
	out = binary.BigEndian.AppendUint32(out, 6)
	out = binary.BigEndian.AppendUint16(out, 0)                        // format 0
	out = binary.BigEndian.AppendUint16(out, 1)                        // one track
	out = binary.BigEndian.AppendUint16(out, uint16(ticksPerQuarter))  // division
	out = append(out, "MTrk"...)
	out = binary.BigEndian.AppendUint32(out, uint32(len(track)))
	out = append(out, track...)
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/io/ -run "VLQ|SMF" -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Commit**

```bash
git add internal/io/smf.go internal/io/smf_test.go
git commit -m "feat(io): add Standard MIDI File encoding"
```

---

### Task 2: FileMIDI

**Files:**
- Create: `internal/io/filemidi.go`
- Test: `internal/io/filemidi_test.go`

**Interfaces:**
- Consumes: `EncodeSMF`, `TimedEvent`, and the existing `MIDIInterface`.
- Produces: `func NewFileMIDI(ticksPerQuarter int) *FileMIDI`, `func (f *FileMIDI) At(tick uint32)`, `func (f *FileMIDI) Write(path string) error`. `FileMIDI` satisfies `MIDIInterface`.

`MIDIInterface` has no notion of time, so `FileMIDI` carries a cursor: `At` sets the tick that subsequent messages are stamped with.

- [ ] **Step 1: Write the failing test**

```go
// internal/io/filemidi_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileMIDISatisfiesInterface(t *testing.T) {
	var _ MIDIInterface = NewFileMIDI(480)
}

func TestFileMIDIRecordsAndWrites(t *testing.T) {
	f := NewFileMIDI(480)
	f.At(0)
	if err := f.SendNoteOn(0, 60, 100); err != nil {
		t.Fatalf("SendNoteOn: %v", err)
	}
	f.At(480)
	if err := f.SendNoteOff(0, 60); err != nil {
		t.Fatalf("SendNoteOff: %v", err)
	}

	path := filepath.Join(t.TempDir(), "out.mid")
	if err := f.Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.HasPrefix(b, []byte("MThd")) {
		t.Errorf("not a MIDI file: % X", b[:4])
	}
	if !bytes.Contains(b, []byte{0x90, 60, 100}) {
		t.Error("note-on not present in the file")
	}
}

func TestFileMIDIEmptyStillWritesValidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.mid")
	if err := NewFileMIDI(480).Write(path); err != nil {
		t.Fatalf("Write: %v", err)
	}
	b, _ := os.ReadFile(path)
	if !bytes.Contains(b, []byte{0xFF, 0x2F, 0x00}) {
		t.Error("an empty file still needs an end-of-track event")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/io/ -run FileMIDI -v`
Expected: FAIL — `undefined: NewFileMIDI`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/io/filemidi.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// A MIDIInterface that writes a Standard MIDI File instead of talking to a device.

package io

import "os"

// FileMIDI records MIDI messages against a tick cursor and writes them out as
// a Standard MIDI File. MIDIInterface carries no timing, so callers move the
// cursor with At before sending each event.
type FileMIDI struct {
	TicksPerQuarter int
	cursor          uint32
	events          []TimedEvent
}

func NewFileMIDI(ticksPerQuarter int) *FileMIDI {
	if ticksPerQuarter <= 0 {
		ticksPerQuarter = 480
	}
	return &FileMIDI{TicksPerQuarter: ticksPerQuarter}
}

// At moves the cursor. Subsequent messages are stamped at this tick.
func (f *FileMIDI) At(tick uint32) { f.cursor = tick }

func (f *FileMIDI) record(data ...byte) error {
	f.events = append(f.events, TimedEvent{Tick: f.cursor, Data: data})
	return nil
}

func (f *FileMIDI) SendNoteOn(channel, note, velocity int) error {
	return f.record(byte(0x90|channel&0x0F), byte(note&0x7F), byte(velocity&0x7F))
}

func (f *FileMIDI) SendNoteOff(channel, note int) error {
	return f.record(byte(0x80|channel&0x0F), byte(note&0x7F), 0)
}

func (f *FileMIDI) SendCC(channel, cc, val int) error {
	return f.record(byte(0xB0|channel&0x0F), byte(cc&0x7F), byte(val&0x7F))
}

func (f *FileMIDI) Close() error { return nil }

// Write encodes everything recorded so far to path.
func (f *FileMIDI) Write(path string) error {
	return os.WriteFile(path, EncodeSMF(f.TicksPerQuarter, f.events), 0o644)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/io/ -run FileMIDI -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Commit**

```bash
git add internal/io/filemidi.go internal/io/filemidi_test.go
git commit -m "feat(io): add FileMIDI writing Standard MIDI Files"
```

---

### Task 3: Hap to MIDI, and the `midi` subcommand

**Files:**
- Create: `internal/io/hapmidi.go`, `internal/io/hapmidi_test.go`
- Modify: `cmd/saint-hubbins/main.go`

**Interfaces:**
- Consumes: `core.Hap`, `FileMIDI`.
- Produces: `func HapToNote(h core.Hap) (note, velocity, channel int, ok bool)` and `func RenderMIDI(pat core.Pattern, cycles int, ticksPerQuarter int) *FileMIDI`, plus `func runMIDI(code, path string, cycles int) error` in `main`.

Note resolution mirrors the audio renderer's precedence: `n` beats `note`, which beats a drum name from `s`. A hap with no pitch information returns `ok == false` and is skipped rather than defaulting to middle C.

- [ ] **Step 1: Write the failing test**

```go
// internal/io/hapmidi_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package io

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func hapWith(v any) core.Hap {
	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	return core.Hap{Whole: &span, Part: span, Value: v}
}

func TestHapToNoteFromNoteName(t *testing.T) {
	note, _, _, ok := HapToNote(hapWith(map[string]any{"note": "c4"}))
	if !ok {
		t.Fatal("c4 should resolve to a note")
	}
	if note != 60 {
		t.Errorf("c4 = %d, want 60", note)
	}
}

func TestHapToNoteNumericWins(t *testing.T) {
	note, _, _, ok := HapToNote(hapWith(map[string]any{"n": 72, "note": "c4"}))
	if !ok || note != 72 {
		t.Errorf("n should take precedence: got %d, ok=%v", note, ok)
	}
}

func TestHapToNoteVelocityFromGain(t *testing.T) {
	_, vel, _, _ := HapToNote(hapWith(map[string]any{"note": 60, "gain": 0.5}))
	if vel != 63 {
		t.Errorf("velocity = %d, want 63 (gain 0.5 of 127)", vel)
	}
}

func TestHapToNoteSkipsPitchlessEvents(t *testing.T) {
	if _, _, _, ok := HapToNote(hapWith(map[string]any{"gain": 0.5})); ok {
		t.Error("an event with no pitch should be skipped, not defaulted")
	}
}

func TestRenderMIDIProducesPairedNoteEvents(t *testing.T) {
	pat := core.Note(core.FastCat(core.Pure(60), core.Pure(64)))
	f := RenderMIDI(pat, 1, 480)
	if len(f.events) != 4 {
		t.Fatalf("got %d events, want 4 (two notes on and off)", len(f.events))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/io/ -run "HapToNote|RenderMIDI" -v`
Expected: FAIL — `undefined: HapToNote`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/io/hapmidi.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Turning pattern events into MIDI notes.

package io

import (
	"math"
	"strconv"
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// drumNotes maps the drum names the engine knows to General MIDI percussion.
var drumNotes = map[string]int{
	"bd": 36, "sd": 38, "hh": 42, "oh": 46, "ch": 42, "cp": 39,
}

// HapToNote resolves a hap to a MIDI note. Precedence matches the audio
// renderer: n, then note, then a drum name. An event carrying no pitch returns
// ok == false so callers skip it rather than emitting a spurious middle C.
func HapToNote(h core.Hap) (note, velocity, channel int, ok bool) {
	velocity, channel = 100, 0
	m, isBag := h.Value.(map[string]any)
	if !isBag {
		if s, isStr := h.Value.(string); isStr {
			if n, found := drumNotes[s]; found {
				return n, velocity, 9, true
			}
			if n, found := noteNameToMIDI(s); found {
				return n, velocity, channel, true
			}
		}
		return 0, 0, 0, false
	}

	if g, found := m["gain"]; found {
		velocity = int(math.Round(toF(g) * 127))
		if velocity < 1 {
			velocity = 1
		}
		if velocity > 127 {
			velocity = 127
		}
	}
	if c, found := m["channel"]; found {
		channel = int(toF(c))
	}
	if v, found := m["n"]; found {
		return int(toF(v)), velocity, channel, true
	}
	if v, found := m["note"]; found {
		if s, isStr := v.(string); isStr {
			if n, parsed := noteNameToMIDI(s); parsed {
				return n, velocity, channel, true
			}
			return 0, 0, 0, false
		}
		return int(toF(v)), velocity, channel, true
	}
	if v, found := m["s"]; found {
		if s, isStr := v.(string); isStr {
			if n, found := drumNotes[s]; found {
				return n, velocity, 9, true
			}
		}
	}
	return 0, 0, 0, false
}

func toF(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		f, _ := strconv.ParseFloat(x, 64)
		return f
	}
	return 0
}

// noteNameToMIDI parses "c4", "f#3", "eb2". Middle C is c4 = 60.
func noteNameToMIDI(s string) (int, bool) {
	if len(s) < 2 {
		return 0, false
	}
	base := map[byte]int{'c': 0, 'd': 2, 'e': 4, 'f': 5, 'g': 7, 'a': 9, 'b': 11}
	n := strings.ToLower(s)
	semi, found := base[n[0]]
	if !found {
		return 0, false
	}
	i := 1
	if len(n) > 2 && (n[1] == '#' || n[1] == 'b') {
		if n[1] == '#' {
			semi++
		} else {
			semi--
		}
		i = 2
	}
	oct, err := strconv.Atoi(n[i:])
	if err != nil {
		return 0, false
	}
	return (oct+1)*12 + semi, true
}

// RenderMIDI queries cycles of pat and records paired note-on/note-off events.
// One cycle is one bar of four quarters, matching the offline audio renderer.
func RenderMIDI(pat core.Pattern, cycles, ticksPerQuarter int) *FileMIDI {
	f := NewFileMIDI(ticksPerQuarter)
	ticksPerCycle := float64(ticksPerQuarter * 4)
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(int64(cycles)))
	for _, h := range haps {
		if h.Whole == nil || !h.HasOnset() {
			continue
		}
		note, vel, ch, ok := HapToNote(h)
		if !ok {
			continue
		}
		f.At(uint32(h.Whole.Begin.Float64() * ticksPerCycle))
		_ = f.SendNoteOn(ch, note, vel)
		f.At(uint32(h.Whole.End.Float64() * ticksPerCycle))
		_ = f.SendNoteOff(ch, note)
	}
	return f
}
```

Note: `RenderMIDI` appends note-offs as it goes, so events are not sorted by tick. Fix that before encoding — add a sort at the top of `FileMIDI.Write`:

```go
func (f *FileMIDI) Write(path string) error {
	sort.SliceStable(f.events, func(i, j int) bool { return f.events[i].Tick < f.events[j].Tick })
	return os.WriteFile(path, EncodeSMF(f.TicksPerQuarter, f.events), 0o644)
}
```

Add `"sort"` to `filemidi.go`'s imports.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/io/ -v`
Expected: PASS, including the pre-existing `MockMIDI` tests

- [ ] **Step 5: Add the `midi` subcommand**

In `cmd/saint-hubbins/main.go`:

```go
// runMIDI evaluates code and writes it out as a Standard MIDI File.
func runMIDI(code, path string, cycles int) error {
	mini.RegisterStringParser()
	pat, _, err := core.Evaluate(code, nil)
	if err != nil {
		pat = mini.Mini(code)
	}
	f := io.RenderMIDI(pat, cycles, 480)
	if err := f.Write(path); err != nil {
		return err
	}
	fmt.Printf("wrote %s (%d cycles)\n", path, cycles)
	return nil
}
```

Import `codeberg.org/uzu/saint-hubbins/internal/io` — note this collides with the standard `io` package if that is also imported in this file; alias it as `shio "codeberg.org/uzu/saint-hubbins/internal/io"` and use `shio.RenderMIDI`.

Add the dispatch case and a usage line:

```go
	case "midi":
		if len(os.Args) < 4 {
			fmt.Println("midi <out.mid> <code> [cycles]")
			os.Exit(1)
		}
		cycles := 4
		if len(os.Args) >= 5 {
			if v, err := strconv.Atoi(os.Args[4]); err == nil {
				cycles = v
			}
		}
		if err := runMIDI(os.Args[3], os.Args[2], cycles); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
```

- [ ] **Step 6: Verify by hand and commit**

```bash
go run ./cmd/saint-hubbins midi /tmp/test.mid "c3 e3 g3 c4" 4
```

Open `/tmp/test.mid` in a DAW or `file /tmp/test.mid` (expect "Standard MIDI data"). Then update the documented usage transcript in `docs/tutorial/01-first-sounds.md` and the CLI Reference table in `README.md`.

```bash
go vet ./... && go test ./... -race -count=1
git add internal/io/ cmd/saint-hubbins/main.go docs/tutorial/01-first-sounds.md README.md
git commit -m "feat(io): render patterns to Standard MIDI Files"
```

---

### Task 4: Live MIDI through the session sink

**Files:**
- Create: `internal/session/midisink.go`, `internal/session/midisink_test.go`

**Interfaces:**
- Consumes: `session.Sink` (from the OSC plan, Task 5), `io.MIDIInterface`, `io.HapToNote`.
- Produces: `type MIDISink struct { Out shio.MIDIInterface }` implementing `Sink`. `internal/io` is imported as `shio` throughout, because the package name collides with the standard library's `io`.

This makes any `MIDIInterface` playable live. With `MockMIDI` it is testable; with a device backend it makes sound. Note-off is scheduled with a timer rather than blocking the scheduler callback.

- [ ] **Step 1: Write the failing test**

```go
// internal/session/midisink_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"testing"
	"time"

	shio "codeberg.org/uzu/saint-hubbins/internal/io"
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestMIDISinkEmitsNoteOnThenOff(t *testing.T) {
	mock := &shio.MockMIDI{}
	sink := &MIDISink{Out: mock}

	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	h := core.Hap{Whole: &span, Part: span, Value: map[string]any{"note": 60}}

	if err := sink.Play(h, time.Now(), 0.5, 0.05); err != nil {
		t.Fatalf("Play: %v", err)
	}
	if len(mock.Messages) == 0 {
		t.Fatal("no note-on was sent")
	}
	// The note-off is scheduled; give it room to land.
	time.Sleep(200 * time.Millisecond)
	if len(mock.Messages) < 2 {
		t.Errorf("expected a note-off to follow, got %v", mock.Messages)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/session/ -run MIDISink -v`
Expected: FAIL — `undefined: MIDISink`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/session/midisink.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Playing scheduled events through a MIDI interface.

package session

import (
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	shio "codeberg.org/uzu/saint-hubbins/internal/io"
)

// MIDISink plays haps through any MIDIInterface. The note-off is scheduled on
// a timer so the scheduler callback never blocks for the length of a note.
type MIDISink struct {
	Out shio.MIDIInterface
}

func (m *MIDISink) Play(h core.Hap, at time.Time, cps, duration float64) error {
	note, vel, ch, ok := shio.HapToNote(h)
	if !ok {
		return nil
	}
	if err := m.Out.SendNoteOn(ch, note, vel); err != nil {
		return err
	}
	hold := time.Duration(duration * float64(time.Second))
	if hold <= 0 {
		hold = 100 * time.Millisecond
	}
	time.AfterFunc(hold, func() { _ = m.Out.SendNoteOff(ch, note) })
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/session/ -race -v`
Expected: PASS. If the race detector flags `MockMIDI.Messages`, that is real — the timer writes from another goroutine. Guard `MockMIDI` with a mutex in `internal/io/midi.go`, keeping its exported fields and methods unchanged.

- [ ] **Step 5: Commit**

```bash
git add internal/session/midisink.go internal/session/midisink_test.go internal/io/midi.go
git commit -m "feat(session): play scheduled haps through a MIDI interface"
```

---

### Task 5: Device backend — DO NOT START WITHOUT A DECISION

See "Decision required before Task 5" above. This task adds a third-party dependency and platform code, changing the project's "single binary, no runtime deps" story. Bring it to the maintainer before writing any of it.

Sketch, for the discussion only: add `gitlab.com/gomidi/midi/v2` behind a `//go:build midi` tag, implement `MIDIInterface` over an output port, expose it as `--midi-device` on the `play` subcommand. Everything else in this plan already works without it.
