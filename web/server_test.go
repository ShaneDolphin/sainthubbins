package web

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestServerHandler(t *testing.T) {
	s := NewServer(":0")
	h := s.Handler()
	if h == nil {
		t.Fatalf("handler nil")
	}
}

func TestPianorollHandler(t *testing.T) {
	s := NewServer(":0")
	h := s.Handler()
	body, _ := json.Marshal(map[string]string{"code": `s("bd sd")`})
	req := httptest.NewRequest("POST", "/api/pianoroll", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("pianoroll status %d", w.Code)
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode %v", err)
	}
	if _, ok := resp["haps"]; !ok {
		t.Fatalf("haps missing")
	}
}

func TestEvaluateHandler(t *testing.T) {
	s := NewServer(":0")
	h := s.Handler()
	body, _ := json.Marshal(map[string]string{"code": `s("bd sd")`})
	req := httptest.NewRequest("POST", "/api/evaluate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("evaluate status %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("CORS missing")
	}
}
