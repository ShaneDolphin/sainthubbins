// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package sampler

// Sample represents a loaded audio sample
type Sample struct {
	Name string
	URL  string
	Data []float32
}

// Server is sample-server.mjs stub
type Server struct {
	Samples map[string]Sample
}

func NewServer() *Server { return &Server{Samples: map[string]Sample{}} }

func (s *Server) Load(name, url string) error {
	s.Samples[name] = Sample{Name: name, URL: url}
	return nil
}
