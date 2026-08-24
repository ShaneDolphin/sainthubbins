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
// A hap's value is normally a control bag (map[string]any). Raw mini-notation
// produces bare values: a bare string becomes s (sound name), and a bare Go
// numeric type (int, int64, float32, float64, uint64) becomes n (note number).
// Note that mini.go stores raw numeric tokens as strings, so "3" arrives as s,
// not n; numeric identity requires either a Go numeric type or a control such
// as core.N or core.Note which produce a control bag. cps and delta always
// travel with the event because SuperDirt sizes its envelope from them.
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
	case int64:
		out = append(out, "n", v)
	case float32:
		out = append(out, "n", v)
	case float64:
		out = append(out, "n", v)
	case uint64:
		// Convert to int64 since the OSC encoder only handles signed integers.
		out = append(out, "n", int64(v))
	default:
		// Unrecognized types (including nil) produce no s/n parameter.
		// This preserves the default synth if one is configured upstream,
		// and catches silent data loss from unexpected type mismatches.
	}

	out = append(out, "cps", cps, "delta", duration)
	return out
}
