// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins live console server — HTTP + WASM bridge.

package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

var consoleTemplate = template.Must(template.New("console").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Saint Hubbins — Live Console</title>
<style>
body{font-family:monospace;background:#1a1a1a;color:#eee;margin:0;padding:20px}
#editor{width:100%;height:200px;background:#2a2a2a;color:#eee;border:1px solid #444;padding:10px;font-size:14px}
#output{white-space:pre-wrap;background:#222;padding:10px;margin-top:10px;min-height:100px}
button{background:#444;color:#eee;border:none;padding:8px 16px;cursor:pointer;margin:5px}
button:hover{background:#555}
footer{margin-top:24px;font-size:12px;color:#888}
footer em{color:#b5a46a}
</style>
</head>
<body>
<h1>Saint Hubbins — Live Console</h1>
<p>Go pattern engine running natively. Type pattern code like <code>s("bd sd")</code> or mini <code>"bd ~ sd"</code> — these go to eleven.</p>
<textarea id="editor">s("bd sd")</textarea>
<br>
<button onclick="evaluate()">Evaluate</button>
<button onclick="hush()">Hush</button>
<div id="output"></div>
<footer><em>Stonehenge — 18" edition</em> &mdash; Saint Hubbins live console (WASM bridge: <code>saintHubbins.queryPattern</code> &rarr; <code>saint-hubbins.wasm</code>)</footer>
<script>
async function evaluate(){
  const code=document.getElementById('editor').value;
  const res=await fetch('/api/evaluate',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({code:code})});
  const data=await res.json();
  document.getElementById('output').textContent=JSON.stringify(data,null,2);
}
function hush(){document.getElementById('output').textContent='[hush]';}
</script>
</body>
</html>`))

type EvaluateRequest struct {
	Code string `json:"code"`
}

type EvaluateResponse struct {
	Haps []map[string]any `json:"haps"`
	Error string `json:"error,omitempty"`
}

// Server is Go HTTP server serving console + API.
type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	if addr == "" {
		addr = ":8080"
	}
	return &Server{Addr: addr}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRoot)
	mux.HandleFunc("/api/evaluate", s.handleEvaluate)
	mux.HandleFunc("/api/pianoroll", s.handlePianoroll)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){ _, _ = w.Write([]byte("ok")) })
	// Static fallback (WASM)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_ = consoleTemplate.Execute(w, nil)
}

func (s *Server) handleEvaluate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 400)
		return
	}
	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var pat core.Pattern
	mini.RegisterStringParser()
	code := req.Code
	if p, _, err := core.Evaluate(code, nil); err == nil {
		pat = p
	} else {
		pat = mini.Mini(code)
		if pat.Query == nil {
			pat = core.Pure(code)
		}
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	out := make([]map[string]any, len(haps))
	for i, h := range haps {
		m := map[string]any{
			"part": h.Part.String(),
			"value": h.Value,
		}
		if h.Whole != nil {
			m["whole"] = h.Whole.String()
		}
		out[i] = m
	}
	resp := EvaluateResponse{Haps: out}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(resp)
}

func (s *Server) handlePianoroll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", 400)
		return
	}
	var req EvaluateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	mini.RegisterStringParser()
	code := req.Code
	var pat core.Pattern
	if p, _, err := core.Evaluate(code, nil); err == nil {
		pat = p
	} else {
		pat = mini.Mini(code)
		if pat.Query == nil {
			pat = core.Pure(code)
		}
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(2))
	type Resp struct {
		Haps []map[string]any `json:"haps"`
	}
	out := make([]map[string]any, len(haps))
	for i, h := range haps {
		m := map[string]any{
			"part": h.Part.String(),
			"value": h.Value,
		}
		if h.Whole != nil {
			m["whole"] = h.Whole.String()
			m["time"] = h.Whole.Begin.Float64()
			m["duration"] = h.Whole.Duration().Float64()
		}
		out[i] = m
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(Resp{Haps: out})
}

func (s *Server) Start() error {
	fmt.Printf("Saint Hubbins console listening on http://localhost%s\n", s.Addr)
	return http.ListenAndServe(s.Addr, s.Handler())
}
