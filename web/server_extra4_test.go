// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluateEmpty(t *testing.T) {
	s := NewServer("")
	mux := s.Handler()
	body := strings.NewReader(`{"code":""}`)
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("empty evaluate status %d body %q", w.Code, w.Body.String())
	}
}

func TestPianorollEmpty(t *testing.T) {
	s := NewServer("")
	mux := s.Handler()
	body := strings.NewReader(`{"code":"~"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/pianoroll", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("pianoroll ~ status %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "haps") {
		t.Fatalf("missing haps %q", w.Body.String())
	}
}
