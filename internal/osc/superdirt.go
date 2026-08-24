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
// numeric type becomes n (note number). Note that mini.go stores raw numeric
// tokens as strings, so "3" arrives as s, not n; numeric identity requires
// either a Go numeric type or a control such as core.N or core.Note which
// produce a control bag. The bare int/float64 cases below are reachable only
// from the Go API (e.g. a pattern built from core.N directly) — mini-notation
// never hands DirtArgs a bare Go numeric, only strings, so those cases have
// no mini-notation-driven test coverage and that is expected, not a gap.
// cps and delta always travel with the event because SuperDirt sizes its
// envelope from them, so a bag that also carries a "cps" key (Cyclist reads
// one out of the bag to retune) is deliberately skipped in the loop below
// rather than sent twice.
func DirtArgs(h core.Hap, cps, duration float64) []any {
	out := make([]any, 0, 16)

	switch v := h.Value.(type) {
	case map[string]any:
		// Sorted so the wire format is deterministic and tests are stable.
		keys := make([]string, 0, len(v))
		for k := range v {
			if strings.HasPrefix(k, "_") || k == "cps" || k == "delta" {
				continue // engine internals (_cps) or values added below
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if val, ok := normalizeArg(v[k]); ok {
				out = append(out, k, val)
			}
		}
	case string:
		out = append(out, "s", v)
	default:
		if val, ok := normalizeArg(v); ok {
			out = append(out, "n", val)
		}
	}

	out = append(out, "cps", cps, "delta", duration)
	return out
}

// normalizeArg converts a value coming from a hap (a bare value, or one
// pulled out of a control bag) into a type EncodeMessage can put on the
// wire. Bags can legitimately carry core.Fraction (duration, clip) or
// anything createParam stored verbatim, and those must not reach the
// encoder unconverted — see internal/core/hap.go's `case Fraction:` handling
// for duration/clip, and the bug this fixes: an unencodable bag value used
// to make EncodeMessage error and the whole event vanish.
func normalizeArg(v any) (any, bool) {
	switch x := v.(type) {
	case core.Fraction:
		return x.Float64(), true
	case string, int, int64, float32, float64:
		return x, true
	case bool:
		if x {
			return int64(1), true
		}
		return int64(0), true
	case uint:
		return int64(x), true
	case uint32:
		return int64(x), true
	case uint64:
		return int64(x), true
	case int8:
		return int64(x), true
	case int16:
		return int64(x), true
	case int32:
		return int64(x), true
	default:
		// Anything else (including nil) cannot be represented on the wire.
		// Skip it rather than handing EncodeMessage a value it will reject —
		// that would previously error out the whole message and drop every
		// other key/value pair riding along with it.
		return nil, false
	}
}
