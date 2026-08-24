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
