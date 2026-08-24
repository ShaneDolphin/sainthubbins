// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package desktop

// Bridge is Tauri desktop bridge stub (midir/rosc/tokio)
type Bridge struct {
	MIDIEnabled bool
	OSCEnabled  bool
}

func New() *Bridge { return &Bridge{} }
