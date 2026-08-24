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
