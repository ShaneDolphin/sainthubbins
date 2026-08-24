// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEvaluateMini(t *testing.T) {
	s := NewServer("")
	mux := s.Handler()
	body := strings.NewReader(`{"code":"s(\"bd sd\")"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("evaluate status %d body %q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "haps") {
		t.Fatalf("missing haps %q", w.Body.String())
	}
	if cors := w.Header().Get("Access-Control-Allow-Origin"); cors != "*" {
		t.Fatalf("cors %q", cors)
	}
}

func TestHealth(t *testing.T) {
	s := NewServer(":0")
	mux := s.Handler()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("health %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "ok" {
		t.Fatalf("health body %q", w.Body.String())
	}
}
