// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func postEvaluate(t *testing.T, code string) map[string]any {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"code": code})
	req := httptest.NewRequest(http.MethodPost, "/api/evaluate", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	NewServer("").Handler().ServeHTTP(rec, req)

	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("bad JSON: %v (%s)", err, rec.Body.String())
	}
	return out
}

func TestConsoleEvaluatesJS(t *testing.T) {
	out := postEvaluate(t, `s("bd sd")`)
	haps, _ := out["haps"].([]any)
	if len(haps) != 2 {
		t.Fatalf("got %d haps, want 2 — the console should evaluate JS now", len(haps))
	}
	first := haps[0].(map[string]any)
	val, ok := first["value"].(map[string]any)
	if !ok || val["s"] != "bd" {
		t.Errorf("value = %v, want a control bag carrying s:bd", first["value"])
	}
}

func TestConsoleStillEvaluatesMiniNotation(t *testing.T) {
	out := postEvaluate(t, "bd*4")
	haps, _ := out["haps"].([]any)
	if len(haps) != 4 {
		t.Errorf("got %d haps, want 4 — mini-notation must keep working", len(haps))
	}
}

func TestConsoleReportsRealErrors(t *testing.T) {
	out := postEvaluate(t, `s("bd").nope()`)
	if out["error"] == nil || out["error"] == "" {
		t.Errorf("a bad method should be reported, got %v", out)
	}
}

// TestConsoleSyntaxErrorNotLiteralHap is the case that could actually fail:
// unparseable JS that mini's own last-resort word-fallback would otherwise
// happily echo back as a single literal-string hap.
func TestConsoleSyntaxErrorNotLiteralHap(t *testing.T) {
	out := postEvaluate(t, `s("bd" +`)
	if out["error"] == nil || out["error"] == "" {
		t.Fatalf("want a JS error reported, got %v", out)
	}
	haps, _ := out["haps"].([]any)
	for _, h := range haps {
		m, _ := h.(map[string]any)
		if v, ok := m["value"].(string); ok && strings.Contains(v, `s("bd"`) {
			t.Fatalf(`got a hap whose value is the literal source text: %v`, m)
		}
	}
}

func TestPianorollReportsRealErrors(t *testing.T) {
	body, _ := json.Marshal(map[string]string{"code": `s("bd").nope()`})
	req := httptest.NewRequest(http.MethodPost, "/api/pianoroll", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	NewServer("").Handler().ServeHTTP(rec, req)

	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("bad JSON: %v (%s)", err, rec.Body.String())
	}
	if out["error"] == nil || out["error"] == "" {
		t.Errorf("a bad method should be reported by pianoroll too, got %v", out)
	}
}
