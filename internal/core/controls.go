// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/controls.mjs (295 params) — generated + hand-written core.

package core

// createParam creates a function that returns a Pattern of a control bag.
// Mirrors JS createParam (handles multi-name and object bags).
func createParam(names ...string) func(any) Pattern {
	isMulti := len(names) > 1
	name := names[0]
	return func(value any) Pattern {
		// Normalize value into a bag map[string]any
		buildBag := func(v any) map[string]any {
			// If value is already a map with .value key, merge
			if m, ok := v.(map[string]any); ok {
				if inner, hasValue := m["value"]; hasValue {
					bag := map[string]any{}
					for k, val := range m {
						if k != "value" {
							bag[k] = val
						}
					}
					// inner may be array for multi
					if isMulti {
						if arr, ok := inner.([]any); ok {
							for i, item := range arr {
								if i < len(names) {
									bag[names[i]] = item
								}
							}
							return bag
						}
					}
					bag[name] = inner
					return bag
				}
				// If multi and value is array
				if isMulti {
					if arr, ok := v.([]any); ok {
						bag := map[string]any{}
						for i, item := range arr {
							if i < len(names) {
								bag[names[i]] = item
							}
						}
						// copy other keys if any? but v is array, no
						return bag
					}
				}
				// Return as-is if already bag? but we need to set this param
				// For simplicity, if map without value key, treat as bag already containing controls
				// and add this param if not multi? But JS logic: withVal does bag merging
			}
			if isMulti {
				if arr, ok := v.([]any); ok {
					bag := map[string]any{}
					for i, item := range arr {
						if i < len(names) {
							bag[names[i]] = item
						}
					}
					return bag
				}
			}
			return map[string]any{name: v}
		}
		// Reify value if it's a Pattern: pattern of bags
		if p, ok := value.(Pattern); ok {
			return p.Fmap(func(v any) any { return buildBag(v) })
		}
		if p, ok := value.(*Pattern); ok && p != nil {
			return p.Fmap(func(v any) any { return buildBag(v) })
		}
		return Pure(buildBag(value))
	}
}

// Core controls — representative subset. Full 295 generated in controls_gen.go.
// These cover the most common pattern controls; the generator will fill the rest.

var (
	S       = createParam("s")
	Sound   = S
	N       = createParam("n")
	Note    = createParam("note")
	Gain    = createParam("gain")
	Velocity = createParam("velocity")
	Vel     = Velocity
	Cutoff  = createParam("cutoff")
	Lpf     = Cutoff
	Resonance = createParam("resonance")
	Lpq     = Resonance
	Delay   = createParam("delay")
	Room    = createParam("room")
	Size    = createParam("size")
	Pan     = createParam("pan")
	Speed   = createParam("speed")
	Begin   = createParam("begin")
	End     = createParam("end")
	Bank    = createParam("bank")
	Orbit   = createParam("orbit")
	Octave  = createParam("octave")
	Coarse  = createParam("coarse")
	CRush  = createParam("crush")
	Shape  = createParam("shape")
	Distort = createParam("distort")
	Cut     = createParam("cut")
	Legato  = createParam("legato")
	Sustain = createParam("sustain")
	Release = createParam("release")
	Attack  = createParam("attack")
	Decay   = createParam("decay")
	Vowel   = createParam("vowel")
	Hpf     = createParam("hpf")
	Hpq     = createParam("hpq")
	Bpf     = createParam("bpf")
	Bpq     = createParam("bpq")
	Freq    = createParam("freq")
	Up      = createParam("up")
	Off     = createParam("off")
)

// Alias handling for multi-name params (e.g., n is alias for note via createParam(["n","note"]))
// N already handles alias via createParam("n")? But true multi would be createParam("n","note")
// We provide helper for that:
func init() {
	// Multi-name example: note and n share same bag key handling
	// Override N to be multi
	N = createParam("n", "note")
	Note = N
}
