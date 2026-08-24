// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Live console animation helpers — Stonehenge edition.

package draw

import (
	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// AnimateParams mirrors JS createParams for draw animation (x,y,w,h,angle,r,fill,smear)
type AnimateParams struct {
	X, Y, W, H, Angle, R, Fill, Smear any
}

// Rescale rescales value in hap context via pattern (stub)
func Rescale(f any, pat core.Pattern) core.Pattern {
	return pat.Fmap(func(v any) any {
		// If v is map with numeric value, scale via f if f is [min,max] etc.
		return v
	})
}

// MoveXY moves x/y
func MoveXY(dx, dy any, pat core.Pattern) core.Pattern {
	return pat.Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			// store dx/dy as context for draw
			m2["x"] = dx
			m2["y"] = dy
			return m2
		}
		return map[string]any{"value": v, "x": dx, "y": dy}
	})
}

// ZoomIn zooms by factor
func ZoomIn(f any, pat core.Pattern) core.Pattern {
	return pat.Fmap(func(v any) any {
		if m, ok := v.(map[string]any); ok {
			m2 := map[string]any{}
			for k, vv := range m {
				m2[k] = vv
			}
			m2["zoom"] = f
			return m2
		}
		return map[string]any{"value": v, "zoom": f}
	})
}

// Framer mirrors JS Framer class (simplified)
type Framer struct {
	LastTime float64
}

func NewFramer() *Framer { return &Framer{} }

func (f *Framer) Tick(time float64) []Event { return nil }

// Drawer mirrors JS Drawer (simplified)
type Drawer struct {
	Ctx string
}

func NewDrawer(ctx string) *Drawer { return &Drawer{Ctx: ctx} }

func (d *Drawer) Draw(events []Event) {}

func GetDrawContext(id string) *Drawer { return NewDrawer(id) }

func CleanupDraw(clearScreen bool, id ...string) {}

func CleanupDrawContext(sessionID string) {}
