// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPianorollSearchParam(t *testing.T) {
	s := NewServer("")
	mux := s.Handler()
	// pianoroll expects POST; verify 400 on GET
	req := httptest.NewRequest(http.MethodGet, "/api/pianoroll", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("pianoroll GET expected 400 got %d", w.Code)
	}
	// POST with JSON
	body := strings.NewReader(`{"code":"bd sd"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/api/pianoroll", body)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != 200 {
		t.Fatalf("pianoroll POST status %d body %q", w2.Code, w2.Body.String())
	}
	if ct := w2.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type %q", ct)
	}
	if !strings.Contains(w2.Body.String(), "haps") {
		t.Fatalf("body missing haps: %q", w2.Body.String())
	}
}

func TestCORSPreflight(t *testing.T) {
	s := NewServer("")
	mux := s.Handler()
	req := httptest.NewRequest(http.MethodOptions, "/api/pianoroll", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("missing CORS header on OPTIONS")
	}
}
