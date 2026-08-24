// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package draw

import "codeberg.org/uzu/saint-hubbins/internal/core"

// SpiralEvents converts haps to spiral coordinates stub
func SpiralEvents(haps []core.Hap) []Event {
	return Pianoroll(haps)
}
