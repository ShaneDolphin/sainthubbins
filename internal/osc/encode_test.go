// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package osc

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"time"
)

func TestPadStringAlwaysNullTerminatedAndAligned(t *testing.T) {
	cases := map[string]int{
		"":      4, // just the null, padded to 4
		"a":     4, // 'a' + null + 2 pad
		"abc":   4, // 3 + null = 4 exactly
		"abcd":  8, // 4 + null needs a fresh block
		"/dirt": 8,
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

func TestEncodeMessageRejectsOutOfRangeInt64(t *testing.T) {
	// Verified before this fix: EncodeMessage("/x", int64(3000000000)) wrote
	// the wrapped-around value -1294967296 with no error, contradicting this
	// function's own doc comment ("an error rather than a silently dropped
	// value").
	if _, err := EncodeMessage("/x", int64(3000000000)); err == nil {
		t.Fatal("want an error for an int64 outside the int32 range, got nil")
	}
	if _, err := EncodeMessage("/x", int64(math.MinInt32)-1); err == nil {
		t.Fatal("want an error for an int64 just below math.MinInt32, got nil")
	}
}

func TestEncodeMessageAcceptsInt32Boundaries(t *testing.T) {
	if _, err := EncodeMessage("/x", int64(math.MaxInt32)); err != nil {
		t.Errorf("int64(math.MaxInt32) should still encode: %v", err)
	}
	if _, err := EncodeMessage("/x", int64(math.MinInt32)); err != nil {
		t.Errorf("int64(math.MinInt32) should still encode: %v", err)
	}
	if _, err := EncodeMessage("/x", int(math.MaxInt32)); err != nil {
		t.Errorf("int(math.MaxInt32) should still encode: %v", err)
	}
	if _, err := EncodeMessage("/x", int(math.MinInt32)); err != nil {
		t.Errorf("int(math.MinInt32) should still encode: %v", err)
	}
}

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

func TestTimetagFractionalSecond(t *testing.T) {
	// A half second must give a fraction of exactly 0x80000000 — the top bit
	// of the 32-bit fractional field, i.e. precisely half the field's range.
	base := time.Unix(0, 0).UTC().Add(500 * time.Millisecond)
	got := timetag(base)
	frac := uint32(got & 0xffffffff)
	if frac != 0x80000000 {
		t.Errorf("fractional field for 0.5s = %#x, want 0x80000000", frac)
	}
}

func TestEncodeBundleEmptyMessages(t *testing.T) {
	got := EncodeBundle(time.Unix(0, 0).UTC())
	if len(got) != 16 {
		t.Errorf("empty bundle length = %d, want 16 (#bundle + timetag)", len(got))
	}
	if !bytes.HasPrefix(got, []byte("#bundle\x00")) {
		t.Errorf("empty bundle must still start with #bundle, got %q", got[:8])
	}
}
