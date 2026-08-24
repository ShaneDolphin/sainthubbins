// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/hap.mjs
package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Hap represents an event active during Part, optionally within Whole.
// If Whole is nil, the event is continuous (Value sampled at midpoint).
// The Part must never extend outside Whole when Whole is non-nil.
type Hap struct {
	Whole    *TimeSpan
	Part     TimeSpan
	Value    any
	Context  map[string]any
	Stateful bool
}

func NewHap(whole *TimeSpan, part TimeSpan, value any, context map[string]any) Hap {
	if context == nil {
		context = map[string]any{}
	}
	return Hap{Whole: whole, Part: part, Value: value, Context: context}
}

// Duration returns the logical duration, respecting value.duration and value.clip if present.
func (h Hap) Duration() Fraction {
	// CheckValue is map with duration/clip
	if m, ok := h.Value.(map[string]any); ok {
		var duration Fraction
		if d, exists := m["duration"]; exists {
			switch v := d.(type) {
			case float64:
				duration = FractionFromFloat(v)
			case int:
				duration = FractionFromInt(int64(v))
			case int64:
				duration = FractionFromInt(v)
			case Fraction:
				duration = v
			case string:
				if parsed, err := ParseFraction(v); err == nil {
					duration = parsed
				} else {
					duration = h.WholeDuration()
				}
			default:
				duration = h.WholeDuration()
			}
		} else {
			duration = h.WholeDuration()
		}
		if clip, exists := m["clip"]; exists {
			var clipFrac Fraction
			switch v := clip.(type) {
			case float64:
				clipFrac = FractionFromFloat(v)
			case int:
				clipFrac = FractionFromInt(int64(v))
			case int64:
				clipFrac = FractionFromInt(v)
			case Fraction:
				clipFrac = v
			default:
				clipFrac = FractionFromInt(1)
			}
			return duration.Mul(clipFrac)
		}
		return duration
	}
	return h.WholeDuration()
}

func (h Hap) WholeDuration() Fraction {
	if h.Whole == nil {
		return h.Part.Duration()
	}
	return h.Whole.End.Sub(h.Whole.Begin)
}

func (h Hap) EndClipped() Fraction {
	if h.Whole == nil {
		return h.Part.End
	}
	return h.Whole.Begin.Add(h.Duration())
}

func (h Hap) IsActive(currentTime Fraction) bool {
	if h.Whole == nil {
		return false
	}
	return !h.Whole.Begin.Gt(currentTime) && !h.EndClipped().Lt(currentTime)
}

func (h Hap) IsInPast(currentTime Fraction) bool {
	return currentTime.Gt(h.EndClipped())
}

func (h Hap) HasOnset() bool {
	return h.Whole != nil && h.Whole.Begin.Equals(h.Part.Begin)
}

func (h Hap) WholeOrPart() TimeSpan {
	if h.Whole != nil {
		return *h.Whole
	}
	return h.Part
}

func (h Hap) WithSpan(fn func(TimeSpan) TimeSpan) Hap {
	var whole *TimeSpan
	if h.Whole != nil {
		w := fn(*h.Whole)
		whole = &w
	}
	part := fn(h.Part)
	return Hap{Whole: whole, Part: part, Value: h.Value, Context: h.Context, Stateful: h.Stateful}
}

func (h Hap) WithValue(fn func(any) any) Hap {
	return Hap{Whole: h.Whole, Part: h.Part, Value: fn(h.Value), Context: h.Context, Stateful: h.Stateful}
}

func (h Hap) SpanEquals(other Hap) bool {
	if h.Whole == nil && other.Whole == nil {
		return true
	}
	if h.Whole == nil || other.Whole == nil {
		return false
	}
	return h.Whole.Equals(*other.Whole)
}

func (h Hap) Equals(other Hap) bool {
	return h.SpanEquals(other) && h.Part.Equals(other.Part) && fmt.Sprintf("%v", h.Value) == fmt.Sprintf("%v", other.Value)
}

func (h Hap) Show(compact bool) string {
	var valueStr string
	if m, ok := h.Value.(map[string]any); ok {
		if compact {
			b, _ := json.Marshal(m)
			s := string(b)
			// JS compact: slice 1,-1, replace quotes and commas
			if len(s) >= 2 {
				s = s[1 : len(s)-1]
			}
			s = strings.ReplaceAll(s, "\"", "")
			s = strings.ReplaceAll(s, ",", " ")
			valueStr = s
		} else {
			b, _ := json.Marshal(m)
			valueStr = string(b)
		}
	} else {
		valueStr = fmt.Sprintf("%v", h.Value)
	}
	var spans string
	if h.Whole == nil {
		spans = "~" + h.Part.Show()
	} else {
		isWhole := h.Whole.Begin.Equals(h.Part.Begin) && h.Whole.End.Equals(h.Part.End)
		if !h.Whole.Begin.Equals(h.Part.Begin) {
			spans = h.Whole.Begin.Show() + " ⇜ "
		}
		if !isWhole {
			spans += "("
		}
		spans += h.Part.Show()
		if !isWhole {
			spans += ")"
		}
		if !h.Whole.End.Equals(h.Part.End) {
			spans += " ⇝ " + h.Whole.End.Show()
		}
	}
	return "[ " + spans + " | " + valueStr + " ]"
}

func (h Hap) ShowWhole(compact bool) string {
	wholeStr := "~"
	if h.Whole != nil {
		wholeStr = h.Whole.Show()
	}
	return fmt.Sprintf("%s: %v", wholeStr, h.Value)
}

func (h Hap) CombineContext(other map[string]any) map[string]any {
	combined := map[string]any{}
	for k, v := range h.Context {
		combined[k] = v
	}
	for k, v := range other {
		combined[k] = v
	}
	// Merge locations slices
	var locs []any
	if a, ok := h.Context["locations"]; ok {
		if sl, ok := a.([]any); ok {
			locs = append(locs, sl...)
		}
	}
	if b, ok := other["locations"]; ok {
		if sl, ok := b.([]any); ok {
			locs = append(locs, sl...)
		}
	}
	if len(locs) > 0 {
		combined["locations"] = locs
	}
	return combined
}

func (h Hap) SetContext(ctx map[string]any) Hap {
	return Hap{Whole: h.Whole, Part: h.Part, Value: h.Value, Context: ctx, Stateful: h.Stateful}
}

func (h Hap) EnsureObjectValue() {
	if _, ok := h.Value.(map[string]any); !ok {
		if _, ok2 := h.Value.(string); ok2 {
			// string values are allowed in early pipeline but EnsureObjectValue enforces map
		}
		if h.Value == nil || fmt.Sprintf("%T", h.Value) == "string" {
			// Allow strings? JS throws: expected hap.value to be object
			// But for Go we allow but log?
		}
		if _, isMap := h.Value.(map[string]any); !isMap {
			panic(fmt.Sprintf("expected hap.value to be an object, but got %q. Hint: append .note() or .s() to the end", fmt.Sprintf("%v", h.Value)))
		}
	}
}

func (h Hap) String() string { return h.Show(false) }
