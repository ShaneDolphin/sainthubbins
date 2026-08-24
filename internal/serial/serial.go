// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package serial

// Serial is Web Serial API stub
type Serial struct {
	Port string
}

func New(port string) *Serial { return &Serial{Port: port} }

func (s *Serial) Write(data []byte) (int, error) { return len(data), nil }
func (s *Serial) Close() error { return nil }
