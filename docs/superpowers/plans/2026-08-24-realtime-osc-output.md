# Real-Time OSC Output Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Saint Hubbins play patterns in real time by sending OSC to SuperDirt, turning the finished-but-unused scheduler into a live instrument.

**Architecture:** `Cyclist` already queries the pattern each tick, filters to haps with an onset, and computes `targetTime`/`duration`/`deadline` — then calls `OnTrigger`, which every construction site currently sets to `nil`. This plan writes a pure-Go OSC 1.0 encoder and UDP client, converts a `Hap` into SuperDirt's flat key/value argument list, and wires that into `OnTrigger`. A new `play` subcommand starts the clock. No audio synthesis is written: SuperDirt supplies samples, stereo and effects.

**Tech Stack:** Go 1.25 standard library only — `net` for UDP, `encoding/binary` for OSC, `context` for lifecycle. No third-party packages.

**Spec:** `docs/superpowers/specs/2026-08-24-remaining-work.md`

## Global Constraints

- Go 1.25.0, module `codeberg.org/uzu/saint-hubbins`.
- **No new third-party dependencies.** OSC encoding is written by hand.
- Every new file starts with: `// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later`
- Engine packages live under `internal/`; importers must be in this module.
- Tests must be hermetic: bind a loopback UDP socket on port 0 and read from it. Never contact a real SuperDirt.
- `go test ./... -race -count=1` and `go vet ./...` must stay clean.
- Do not change the signature of `osc.New(host string, port int) *Client` or remove `SendSuperDirt`/`Close` — `internal/osc/osc_test.go` and `osc_extra_test.go` exercise them and must keep passing.

## File Structure

| File | Responsibility |
|---|---|
| `internal/osc/encode.go` (create) | OSC 1.0 wire format: padded strings, int32/float32, message and bundle encoding. Pure functions, no I/O. |
| `internal/osc/encode_test.go` (create) | Byte-level tests for the encoder against the OSC 1.0 spec. |
| `internal/osc/osc.go` (modify) | `Client` becomes a real UDP sender. Keeps its existing exported surface. |
| `internal/osc/osc_udp_test.go` (create) | Loopback tests: a real UDP listener receives what the client sends. |
| `internal/osc/superdirt.go` (create) | `Hap` → SuperDirt `/dirt/play` argument list. Knows about control bags, not about sockets. |
| `internal/osc/superdirt_test.go` (create) | Conversion tests for bags, bare values and rests. |
| `internal/session/session.go` (modify) | Gains an output sink so a session can drive OSC. |
| `cmd/saint-hubbins/main.go` (modify) | New `play` subcommand. |

---

### Task 1: OSC value and message encoding

**Files:**
- Create: `internal/osc/encode.go`
- Test: `internal/osc/encode_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces: `func EncodeMessage(addr string, args ...any) ([]byte, error)`, and the unexported helper `padString(s string) []byte`. Argument types accepted: `string` (OSC `s`), `int`/`int32`/`int64` (OSC `i`), `float32`/`float64` (OSC `f`). Any other type returns an error.

Background: an OSC message is an address string, then a type-tag string beginning with `,`, then the arguments. Every string is null-terminated and zero-padded to a multiple of 4 bytes. Integers and floats are 32-bit big-endian.

- [ ] **Step 1: Write the failing test**

```go
// internal/osc/encode_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"bytes"
	"testing"
)

func TestPadStringAlwaysNullTerminatedAndAligned(t *testing.T) {
	cases := map[string]int{
		"":       4,  // just the null, padded to 4
		"a":      4,  // 'a' + null + 2 pad
		"abc":    4,  // 3 + null = 4 exactly
		"abcd":   8,  // 4 + null needs a fresh block
		"/dirt":  8,
	}
	for in, wantLen := range cases {
		got := padString(in)
		if len(got) != wantLen {
			t.Errorf("padString(%q) length %d, want %d", in, len(got), wantLen)
		}
		if len(got)%4 != 0 {
			t.Errorf("padString(%q) not 4-byte aligned: %d", in, len(got))
		}
		if got[len(in)] != 0 {
			t.Errorf("padString(%q) missing null terminator", in)
		}
	}
}

func TestEncodeMessageLayout(t *testing.T) {
	got, err := EncodeMessage("/dirt/play", "s", "bd", 3, float32(0.5))
	if err != nil {
		t.Fatalf("EncodeMessage: %v", err)
	}
	// "/dirt/play" is 10 bytes -> padded to 12.
	if !bytes.HasPrefix(got, []byte("/dirt/play\x00\x00")) {
		t.Fatalf("address not encoded first: %q", got[:12])
	}
	// Type tags: , s s i f  -> 5 bytes -> padded to 8.
	if !bytes.Contains(got, []byte(",ssif\x00\x00\x00")) {
		t.Fatalf("type tag string %q not found in %q", ",ssif", got)
	}
	if len(got)%4 != 0 {
		t.Errorf("message not 4-byte aligned: %d", len(got))
	}
}

func TestEncodeMessageRejectsUnsupportedType(t *testing.T) {
	if _, err := EncodeMessage("/x", struct{}{}); err == nil {
		t.Fatal("want an error for an unsupported argument type, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/osc/ -run "PadString|EncodeMessage" -v`
Expected: FAIL — `undefined: padString`, `undefined: EncodeMessage`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/osc/encode.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// OSC 1.0 wire encoding. Pure functions — no sockets here.

package osc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// padString encodes an OSC-string: the bytes, a null terminator, then zero
// padding to the next multiple of four. A string whose length is already a
// multiple of four still gets a full four bytes of padding, because the null
// terminator is mandatory.
func padString(s string) []byte {
	b := make([]byte, 0, len(s)+4)
	b = append(b, s...)
	b = append(b, 0)
	for len(b)%4 != 0 {
		b = append(b, 0)
	}
	return b
}

// EncodeMessage builds an OSC message: address, type tags, then arguments.
// Supported argument types are string, the signed integer kinds, and the
// float kinds. Anything else is an error rather than a silently dropped value.
func EncodeMessage(addr string, args ...any) ([]byte, error) {
	tags := []byte{','}
	var body bytes.Buffer
	for _, a := range args {
		switch v := a.(type) {
		case string:
			tags = append(tags, 's')
			body.Write(padString(v))
		case int:
			tags = append(tags, 'i')
			_ = binary.Write(&body, binary.BigEndian, int32(v))
		case int32:
			tags = append(tags, 'i')
			_ = binary.Write(&body, binary.BigEndian, v)
		case int64:
			tags = append(tags, 'i')
			_ = binary.Write(&body, binary.BigEndian, int32(v))
		case float32:
			tags = append(tags, 'f')
			_ = binary.Write(&body, binary.BigEndian, math.Float32bits(v))
		case float64:
			tags = append(tags, 'f')
			_ = binary.Write(&body, binary.BigEndian, math.Float32bits(float32(v)))
		default:
			return nil, fmt.Errorf("osc: unsupported argument type %T", a)
		}
	}
	out := make([]byte, 0, 64)
	out = append(out, padString(addr)...)
	out = append(out, padString(string(tags))...)
	out = append(out, body.Bytes()...)
	return out, nil
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/osc/ -run "PadString|EncodeMessage" -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Confirm the existing osc tests still pass**

Run: `go test ./internal/osc/ -v`
Expected: PASS, including `TestOSCMessageBasic` and `TestOSCBundle`

- [ ] **Step 6: Commit**

```bash
git add internal/osc/encode.go internal/osc/encode_test.go
git commit -m "feat(osc): add OSC 1.0 message encoding"
```

---

### Task 2: OSC bundles with timetags

**Files:**
- Modify: `internal/osc/encode.go`
- Test: `internal/osc/encode_test.go`

**Interfaces:**
- Consumes: `EncodeMessage`, `padString` from Task 1.
- Produces: `func EncodeBundle(at time.Time, msgs ...[]byte) []byte` and `func timetag(t time.Time) uint64`.

Background: SuperDirt schedules a bundle for the moment in its timetag, which is how a note lands on the beat despite network and process jitter. `Cyclist` already computes `targetTime` per hap, so bundles are what make its latency budget meaningful. An OSC timetag is 64 bits: 32 bits of seconds since 1900-01-01 and 32 bits of fractional second.

- [ ] **Step 1: Write the failing test**

```go
// append to internal/osc/encode_test.go

func TestTimetagUsesNTPEpoch(t *testing.T) {
	// 1970-01-01 UTC is exactly 2208988800 seconds after the NTP epoch.
	unixEpoch := time.Unix(0, 0).UTC()
	got := timetag(unixEpoch)
	if secs := uint32(got >> 32); secs != 2208988800 {
		t.Errorf("seconds field = %d, want 2208988800", secs)
	}
}

func TestEncodeBundleLayout(t *testing.T) {
	m1, _ := EncodeMessage("/a", "x")
	m2, _ := EncodeMessage("/b", 1)
	got := EncodeBundle(time.Unix(0, 0).UTC(), m1, m2)

	if !bytes.HasPrefix(got, []byte("#bundle\x00")) {
		t.Fatalf("bundle must start with #bundle, got %q", got[:8])
	}
	// 8 (#bundle) + 8 (timetag) + per message 4-byte length prefix + payload.
	want := 8 + 8 + 4 + len(m1) + 4 + len(m2)
	if len(got) != want {
		t.Errorf("bundle length %d, want %d", len(got), want)
	}
	// First element's length prefix must equal len(m1).
	gotLen := binary.BigEndian.Uint32(got[16:20])
	if int(gotLen) != len(m1) {
		t.Errorf("first element length prefix %d, want %d", gotLen, len(m1))
	}
}
```

Add `"encoding/binary"` and `"time"` to the test file's imports.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/osc/ -run "Timetag|EncodeBundle" -v`
Expected: FAIL — `undefined: timetag`, `undefined: EncodeBundle`

- [ ] **Step 3: Write minimal implementation**

```go
// append to internal/osc/encode.go
// (add "time" to the import block)

// ntpEpochOffset is the number of seconds between 1900-01-01 and 1970-01-01.
const ntpEpochOffset = 2208988800

// timetag converts a wall-clock time into a 64-bit OSC timetag: seconds since
// the NTP epoch in the high half, fractional second in the low half.
func timetag(t time.Time) uint64 {
	secs := uint64(t.Unix() + ntpEpochOffset)
	frac := uint64(t.Nanosecond()) << 32 / 1e9
	return secs<<32 | frac
}

// EncodeBundle wraps messages in an OSC bundle scheduled for at. SuperDirt
// plays a bundle at its timetag rather than on arrival, which is what lets a
// note land on the beat despite scheduling jitter.
func EncodeBundle(at time.Time, msgs ...[]byte) []byte {
	out := make([]byte, 0, 64)
	out = append(out, padString("#bundle")...)
	out = binary.BigEndian.AppendUint64(out, timetag(at))
	for _, m := range msgs {
		out = binary.BigEndian.AppendUint32(out, uint32(len(m)))
		out = append(out, m...)
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/osc/ -run "Timetag|EncodeBundle" -v`
Expected: PASS (2 tests)

- [ ] **Step 5: Commit**

```bash
git add internal/osc/encode.go internal/osc/encode_test.go
git commit -m "feat(osc): add bundle encoding with NTP timetags"
```

---

### Task 3: Real UDP client

**Files:**
- Modify: `internal/osc/osc.go`
- Test: `internal/osc/osc_udp_test.go`

**Interfaces:**
- Consumes: `EncodeMessage`, `EncodeBundle` from Tasks 1–2.
- Produces: `func (c *Client) Send(addr string, args ...any) error`, `func (c *Client) SendAt(at time.Time, addr string, args ...any) error`. `New(host string, port int) *Client`, `SendSuperDirt([]interface{}) error` and `Close() error` keep their current signatures.

The existing `Client` is a no-op stub. It must become a real UDP sender **without breaking `osc_test.go`**, which calls `New("", 0)` and expects `SendSuperDirt(nil)` and `Close()` to return nil. So: a client with an empty host stays a no-op sink, and dialing happens lazily so construction never fails.

- [ ] **Step 1: Write the failing test**

```go
// internal/osc/osc_udp_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// listener binds a loopback UDP socket and returns it with its port. Tests
// stay hermetic: nothing leaves the machine and no SuperDirt is required.
func listener(t *testing.T) (*net.UDPConn, int) {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn, conn.LocalAddr().(*net.UDPAddr).Port
}

func TestClientSendReachesTheWire(t *testing.T) {
	conn, port := listener(t)
	c := New("127.0.0.1", port)
	defer c.Close()

	if err := c.Send("/dirt/play", "s", "bd"); err != nil {
		t.Fatalf("Send: %v", err)
	}

	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("nothing received: %v", err)
	}
	if !bytes.HasPrefix(buf[:n], []byte("/dirt/play\x00\x00")) {
		t.Errorf("received %q, want a /dirt/play message", buf[:n])
	}
	if !bytes.Contains(buf[:n], []byte("bd\x00")) {
		t.Errorf("received %q, want it to carry the sound name", buf[:n])
	}
}

func TestClientSendAtSendsABundle(t *testing.T) {
	conn, port := listener(t)
	c := New("127.0.0.1", port)
	defer c.Close()

	if err := c.SendAt(time.Now().Add(time.Second), "/dirt/play", "s", "sd"); err != nil {
		t.Fatalf("SendAt: %v", err)
	}
	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("nothing received: %v", err)
	}
	if !bytes.HasPrefix(buf[:n], []byte("#bundle\x00")) {
		t.Errorf("received %q, want a bundle", buf[:n])
	}
}

// A client with no host is a sink, so tests and offline use never need a peer.
func TestClientWithoutHostIsANoOp(t *testing.T) {
	c := New("", 0)
	if err := c.Send("/dirt/play", "s", "bd"); err != nil {
		t.Errorf("Send on a hostless client should be a no-op, got %v", err)
	}
	if err := c.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/osc/ -run "ClientSend|WithoutHost" -v`
Expected: FAIL — `c.Send undefined`, `c.SendAt undefined`

- [ ] **Step 3: Write minimal implementation**

Replace the whole body of `internal/osc/osc.go` with:

```go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// UDP client for SuperDirt and other OSC receivers.

package osc

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Client sends OSC over UDP. A Client with an empty Host is a sink: every send
// succeeds and goes nowhere, so offline rendering and tests need no peer.
type Client struct {
	Host string
	Port int

	mu   sync.Mutex
	conn net.Conn
	dialErr error
	dialed  bool
}

// New creates a client. It does not dial — UDP has no handshake worth doing
// eagerly, and construction should not fail.
func New(host string, port int) *Client { return &Client{Host: host, Port: port} }

// ensure dials once, on first use.
func (c *Client) ensure() (net.Conn, error) {
	if c.Host == "" {
		return nil, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.dialed {
		return c.conn, c.dialErr
	}
	c.dialed = true
	c.conn, c.dialErr = net.Dial("udp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	return c.conn, c.dialErr
}

func (c *Client) write(b []byte) error {
	conn, err := c.ensure()
	if err != nil {
		return err
	}
	if conn == nil {
		return nil // hostless sink
	}
	_, err = conn.Write(b)
	return err
}

// Send transmits one OSC message immediately.
func (c *Client) Send(addr string, args ...any) error {
	msg, err := EncodeMessage(addr, args...)
	if err != nil {
		return err
	}
	return c.write(msg)
}

// SendAt transmits one OSC message inside a bundle timestamped for at, so the
// receiver plays it at that moment rather than on arrival.
func (c *Client) SendAt(at time.Time, addr string, args ...any) error {
	msg, err := EncodeMessage(addr, args...)
	if err != nil {
		return err
	}
	return c.write(EncodeBundle(at, msg))
}

// SendSuperDirt is retained for compatibility with existing callers.
func (c *Client) SendSuperDirt(haps []interface{}) error {
	if len(haps) == 0 {
		return nil
	}
	return c.Send(DirtAddress, haps...)
}

// Close releases the socket if one was opened.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	c.dialed = false
	return err
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/osc/ -run "ClientSend|WithoutHost" -v`
Expected: FAIL — `undefined: DirtAddress`. That constant arrives in Task 4; add it now as a one-line placeholder so this task stands alone:

```go
// DirtAddress is the OSC address SuperDirt listens on for note events.
const DirtAddress = "/dirt/play"
```

Put it at the top of `internal/osc/osc.go`, below the imports. Re-run; expected: PASS (3 tests).

- [ ] **Step 5: Confirm the pre-existing osc tests still pass**

Run: `go test ./internal/osc/ -v`
Expected: PASS — `TestOSCMessageBasic` and `TestOSCBundle` must still be green. If `TestOSCMessageBasic` now fails, it is because it dials `127.0.0.1:57120`; a UDP write to a port with no listener still succeeds locally, so investigate rather than weakening the test.

- [ ] **Step 6: Commit**

```bash
git add internal/osc/osc.go internal/osc/osc_udp_test.go
git commit -m "feat(osc): replace the no-op stub with a real UDP client"
```

---

### Task 4: Hap to SuperDirt arguments

**Files:**
- Create: `internal/osc/superdirt.go`
- Test: `internal/osc/superdirt_test.go`

**Interfaces:**
- Consumes: `core.Hap` from `internal/core`.
- Produces: `func DirtArgs(h core.Hap, cps, duration float64) []any`. Returns a flat alternating key/value list — SuperDirt's `/dirt/play` takes `"s", "bd", "n", 1, ...`.

A hap's value is either a control bag (`map[string]any`, the normal case) or a bare value from raw mini-notation. A bare string becomes `s`; a bare number becomes `n`. Every message carries `cps` and `delta` (the event's duration in seconds), which SuperDirt needs to size its envelope.

- [ ] **Step 1: Write the failing test**

```go
// internal/osc/superdirt_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// pairs turns the flat argument list into a map so tests can assert on keys
// without depending on ordering.
func pairs(t *testing.T, args []any) map[string]any {
	t.Helper()
	if len(args)%2 != 0 {
		t.Fatalf("argument list must be key/value pairs, got %d items: %v", len(args), args)
	}
	m := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			t.Fatalf("key at %d is %T, want string", i, args[i])
		}
		m[k] = args[i+1]
	}
	return m
}

func hapOf(v any) core.Hap {
	span := core.NewTimeSpan(core.FractionFromInt(0), core.NewFraction(1, 4))
	return core.Hap{Whole: &span, Part: span, Value: v}
}

func TestDirtArgsFromControlBag(t *testing.T) {
	h := hapOf(map[string]any{"s": "bd", "gain": 0.8, "n": 2})
	m := pairs(t, DirtArgs(h, 0.5, 0.5))

	if m["s"] != "bd" {
		t.Errorf("s = %v, want bd", m["s"])
	}
	if m["gain"] != 0.8 {
		t.Errorf("gain = %v, want 0.8", m["gain"])
	}
	if m["n"] != 2 {
		t.Errorf("n = %v, want 2", m["n"])
	}
	if _, ok := m["cps"]; !ok {
		t.Error("every message must carry cps")
	}
	if m["delta"] != 0.5 {
		t.Errorf("delta = %v, want the event duration in seconds", m["delta"])
	}
}

func TestDirtArgsFromBareValues(t *testing.T) {
	// Raw mini-notation produces bare values rather than control bags.
	m := pairs(t, DirtArgs(hapOf("bd"), 0.5, 0.25))
	if m["s"] != "bd" {
		t.Errorf("a bare string should become the sound name, got %v", m["s"])
	}

	m = pairs(t, DirtArgs(hapOf(3), 0.5, 0.25))
	if m["n"] != 3 {
		t.Errorf("a bare number should become n, got %v", m["n"])
	}
}

func TestDirtArgsSkipsInternalKeys(t *testing.T) {
	h := hapOf(map[string]any{"s": "bd", "_cps": 0.5})
	m := pairs(t, DirtArgs(h, 0.5, 0.25))
	if _, ok := m["_cps"]; ok {
		t.Error("underscore-prefixed keys are engine internals and must not be sent")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/osc/ -run DirtArgs -v`
Expected: FAIL — `undefined: DirtArgs`

- [ ] **Step 3: Write minimal implementation**

```go
// internal/osc/superdirt.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Translating pattern events into SuperDirt's /dirt/play arguments.

package osc

import (
	"sort"
	"strings"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// DirtArgs flattens a hap into SuperDirt's alternating key/value argument list.
//
// A hap's value is normally a control bag; raw mini-notation produces a bare
// string or number instead, which maps to s and n respectively. cps and delta
// always travel with the event because SuperDirt sizes its envelope from them.
func DirtArgs(h core.Hap, cps, duration float64) []any {
	out := make([]any, 0, 16)

	switch v := h.Value.(type) {
	case map[string]any:
		// Sorted so the wire format is deterministic and tests are stable.
		keys := make([]string, 0, len(v))
		for k := range v {
			if strings.HasPrefix(k, "_") {
				continue // engine internals such as _cps
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			out = append(out, k, v[k])
		}
	case string:
		out = append(out, "s", v)
	case int:
		out = append(out, "n", v)
	case float64:
		out = append(out, "n", v)
	}

	out = append(out, "cps", cps, "delta", duration)
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/osc/ -run DirtArgs -v`
Expected: PASS (3 tests)

- [ ] **Step 5: Remove the placeholder constant duplication**

`DirtAddress` was added to `osc.go` in Task 3. Leave it there; do not redeclare it in `superdirt.go`.

Run: `go build ./... && go test ./internal/osc/`
Expected: build clean, all osc tests pass

- [ ] **Step 6: Commit**

```bash
git add internal/osc/superdirt.go internal/osc/superdirt_test.go
git commit -m "feat(osc): convert haps to SuperDirt arguments"
```

---

### Task 5: Wire the scheduler to an output sink

**Files:**
- Modify: `internal/session/session.go`
- Test: `internal/session/session_output_test.go` (create)

**Interfaces:**
- Consumes: `core.NewCyclist`, `core.Hap`, `osc.DirtArgs`, `osc.Client` from earlier tasks.
- Produces: `type Sink interface { Play(h core.Hap, at time.Time, cps, duration float64) error }`, `func (r *Session) SetSink(s Sink)`, `func (r *Session) Start(ctx context.Context) error`, `func (r *Session) Stop()`, and `type OSCSink struct{ Client *osc.Client }` implementing `Sink`.

`NewSession` currently builds its `Cyclist` with `OnTrigger = nil`, so nothing happens per hap. A `Sink` interface keeps `session` independent of OSC and makes the trigger path testable without a socket.

Note `Cyclist.OnTrigger` has signature `func(hap core.Hap, deadline, duration, cps, targetTime float64)`, where `targetTime` is an absolute wall-clock time in seconds since the Unix epoch.

- [ ] **Step 1: Write the failing test**

```go
// internal/session/session_output_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package session

import (
	"sync"
	"testing"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// recorder is a Sink that remembers what it was asked to play.
type recorder struct {
	mu   sync.Mutex
	haps []core.Hap
}

func (r *recorder) Play(h core.Hap, at time.Time, cps, duration float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.haps = append(r.haps, h)
	return nil
}

func (r *recorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.haps)
}

func TestSessionSendsHapsToItsSink(t *testing.T) {
	s := NewSession()
	rec := &recorder{}
	s.SetSink(rec)

	if _, err := s.Evaluate("bd*4"); err != nil {
		t.Fatalf("Evaluate: %v", err)
	}

	// Drive the trigger path directly rather than waiting on wall-clock ticks,
	// so the test is fast and deterministic.
	h := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0]
	s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))

	if rec.count() != 1 {
		t.Fatalf("sink received %d haps, want 1", rec.count())
	}
}

func TestSessionWithoutSinkDoesNotPanic(t *testing.T) {
	s := NewSession()
	if _, err := s.Evaluate("bd"); err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	h := s.Pattern.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))[0]
	s.Cyclist.OnTrigger(h, 0, 0.5, 0.5, float64(time.Now().Unix()))
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/session/ -run Sink -v`
Expected: FAIL — `s.SetSink undefined`

- [ ] **Step 3: Write minimal implementation**

Replace `internal/session/session.go` with:

```go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live session — evaluation, pattern, scheduler and output.

package session

import (
	"context"
	"sync"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
	"codeberg.org/uzu/saint-hubbins/internal/osc"
)

// Sink receives scheduled events. Keeping this an interface means the session
// knows nothing about OSC, and the trigger path can be tested without a socket.
type Sink interface {
	Play(h core.Hap, at time.Time, cps, duration float64) error
}

// OSCSink plays events through SuperDirt.
type OSCSink struct{ Client *osc.Client }

func (s *OSCSink) Play(h core.Hap, at time.Time, cps, duration float64) error {
	return s.Client.SendAt(at, osc.DirtAddress, osc.DirtArgs(h, cps, duration)...)
}

// Session ties evaluation, pattern, scheduler and output together.
type Session struct {
	mu      sync.RWMutex
	Pattern core.Pattern
	Cyclist *core.Cyclist
	sink    Sink
}

// NewSession creates a new live session (these go to eleven).
func NewSession() *Session {
	mini.RegisterStringParser()
	s := &Session{
		Cyclist: core.NewCyclist(0.1, nil, nil),
		Pattern: core.Silence(),
	}
	// Cyclist computes targetTime as absolute seconds since the Unix epoch.
	s.Cyclist.OnTrigger = func(h core.Hap, deadline, duration, cps, targetTime float64) {
		s.mu.RLock()
		sink := s.sink
		s.mu.RUnlock()
		if sink == nil {
			return
		}
		at := time.Unix(0, int64(targetTime*1e9))
		_ = sink.Play(h, at, cps, duration)
	}
	return s
}

// SetSink installs the output. Passing nil silences the session.
func (r *Session) SetSink(s Sink) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sink = s
}

// Start runs the scheduler until ctx is cancelled.
func (r *Session) Start(ctx context.Context) error { return r.Cyclist.Start(ctx) }

// Stop halts the scheduler.
func (r *Session) Stop() { r.Cyclist.Stop() }

func (r *Session) Evaluate(code string) (core.Pattern, error) {
	pat, _, err := core.Evaluate(code, nil)
	if err != nil {
		pat = mini.Mini(code)
		if pat.Query == nil {
			pat = core.Pure(code)
		}
		err = nil
	}
	r.mu.Lock()
	r.Pattern = pat
	r.mu.Unlock()
	r.Cyclist.SetPattern(pat)
	return pat, err
}

func (r *Session) Hush() {
	r.mu.Lock()
	r.Pattern = core.Silence()
	r.mu.Unlock()
	r.Cyclist.SetPattern(core.Silence())
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/session/ -v`
Expected: PASS, including any pre-existing session tests

- [ ] **Step 5: Commit**

```bash
git add internal/session/session.go internal/session/session_output_test.go
git commit -m "feat(session): route scheduled haps to a pluggable output sink"
```

---

### Task 6: The play subcommand

**Files:**
- Modify: `cmd/saint-hubbins/main.go`
- Test: `cmd/saint-hubbins/play_test.go` (create)

**Interfaces:**
- Consumes: `session.NewSession`, `session.OSCSink`, `osc.New` from Task 5.
- Produces: `func runPlay(code string, host string, port int, seconds float64, out io.Writer) error` — extracted so it is testable without `os.Exit`.

- [ ] **Step 1: Write the failing test**

```go
// cmd/saint-hubbins/play_test.go
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

func TestRunPlaySendsToSuperDirt(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer conn.Close()
	port := conn.LocalAddr().(*net.UDPAddr).Port

	received := make(chan []byte, 8)
	go func() {
		buf := make([]byte, 2048)
		for {
			_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				close(received)
				return
			}
			cp := make([]byte, n)
			copy(cp, buf[:n])
			received <- cp
		}
	}()

	var out bytes.Buffer
	if err := runPlay("bd*4", "127.0.0.1", port, 1.0, &out); err != nil {
		t.Fatalf("runPlay: %v", err)
	}

	select {
	case msg, ok := <-received:
		if !ok {
			t.Fatal("listener closed without receiving anything")
		}
		if !bytes.Contains(msg, []byte("bd")) {
			t.Errorf("received %q, want it to carry the sound name", msg)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no OSC arrived within the timeout")
	}

	if !strings.Contains(out.String(), "127.0.0.1") {
		t.Errorf("runPlay should report where it is sending, got %q", out.String())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./cmd/saint-hubbins/ -run RunPlay -v`
Expected: FAIL — `undefined: runPlay`

- [ ] **Step 3: Write minimal implementation**

Add to `cmd/saint-hubbins/main.go`:

```go
// runPlay evaluates code and streams it to SuperDirt over OSC for seconds.
// It is separate from the CLI dispatch so it can be tested without os.Exit.
func runPlay(code, host string, port int, seconds float64, out io.Writer) error {
	client := osc.New(host, port)
	defer client.Close()

	s := session.NewSession()
	s.SetSink(&session.OSCSink{Client: client})
	if _, err := s.Evaluate(code); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(seconds*float64(time.Second)))
	defer cancel()

	fmt.Fprintf(out, "playing %q to %s:%d for %.1fs — these go to eleven\n",
		code, host, port, seconds)
	if err := s.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	s.Stop()
	return nil
}
```

Add the imports `context`, `io`, `time`, `codeberg.org/uzu/saint-hubbins/internal/osc` and `codeberg.org/uzu/saint-hubbins/internal/session`.

Then add the dispatch case in `main()`, beside the existing `eval`/`serve`/`render` cases:

```go
	case "play":
		if len(os.Args) < 3 {
			fmt.Println("play <code> [host] [port] [seconds]")
			os.Exit(1)
		}
		host, port, secs := "127.0.0.1", 57120, 8.0
		if len(os.Args) >= 4 {
			host = os.Args[3]
		}
		if len(os.Args) >= 5 {
			if v, err := strconv.Atoi(os.Args[4]); err == nil {
				port = v
			}
		}
		if len(os.Args) >= 6 {
			if v, err := strconv.ParseFloat(os.Args[5], 64); err == nil {
				secs = v
			}
		}
		if err := runPlay(os.Args[2], host, port, secs, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
```

Add `strconv` to the imports, and add this line to the usage block printed when no subcommand is given:

```go
		fmt.Println("  play <code> [host] [port] [secs] — stream to SuperDirt over OSC")
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./cmd/saint-hubbins/ -run RunPlay -v`
Expected: PASS

- [ ] **Step 5: Update the documented usage transcript**

`docs/tutorial/01-first-sounds.md` shows the binary's usage output byte-for-byte and a test-free reader will notice the drift. Run `go run ./cmd/saint-hubbins`, copy the exact output, and replace the transcript in that file. Also add `play` to the CLI Reference table in `README.md`.

Run: `go run ./cmd/saint-hubbins`
Expected: usage now lists five subcommands

- [ ] **Step 6: Full verification and commit**

```bash
go vet ./...
go test ./... -race -count=1
git add cmd/saint-hubbins/main.go cmd/saint-hubbins/play_test.go docs/tutorial/01-first-sounds.md README.md
git commit -m "feat(cli): add play subcommand streaming patterns to SuperDirt"
```

---

## Manual verification

Automated tests never contact a real SuperDirt. Once the tasks are done, confirm end to end:

1. Start SuperCollider with SuperDirt listening on 57120.
2. `go run ./cmd/saint-hubbins play "bd*4, hh*8" 127.0.0.1 57120 8`
3. You should hear a four-on-the-floor kick with eight hats, in time, with real samples.
4. `go run ./cmd/saint-hubbins play "bd(3,8)" 127.0.0.1 57120 8` should sound off-grid.

If nothing plays, check in this order: SuperDirt is listening (`netstat -an | grep 57120`), the pattern produces haps (`go run ./cmd/saint-hubbins eval "bd*4"`), and packets leave (`tcpdump -i lo0 udp port 57120`).

## Follow-on work this unblocks

Once events reach SuperDirt, the sine renderer stops being the only output, which retires most of `docs/tutorial/08-limitations.md`: samples, stereo, and the ~290 controls the offline renderer ignores are all handled by SuperDirt. Update that chapter and the README Features list to distinguish live output from the offline preview.
