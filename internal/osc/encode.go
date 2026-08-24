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
