// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package webaudio

type AudioContext struct{ SampleRate int }
func NewAudioContext() *AudioContext { return &AudioContext{SampleRate: 44100} }
func (a *AudioContext) Close() {}
