// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// cmd/saint-hubbins — native CLI for Saint Hubbins.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"codeberg.org/uzu/saint-hubbins/internal/audio"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
	"codeberg.org/uzu/saint-hubbins/web"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: saint-hubbins <eval|serve|render|query> [args]")
		fmt.Println("  eval <code>        — evaluate pattern string")
		fmt.Println("  query              — demo query: Stack(s(\"bd\"), s(\"sd\"))")
		fmt.Println("  serve [addr]       — start live console server (default :8080)")
		fmt.Println("  render <out.wav> <code> — offline render to WAV")
		fmt.Println("  (also available as 'hubbins' — these go to eleven)")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "query":
		demoQuery()
	case "eval":
		if len(os.Args) < 3 {
			fmt.Println("eval <code>")
			os.Exit(1)
		}
		demoEval(os.Args[2])
	case "serve":
		addr := ":8080"
		if len(os.Args) >= 3 {
			addr = os.Args[2]
		}
		serve(addr)
	case "render":
		if len(os.Args) < 4 {
			fmt.Println("render <out.wav> <code>")
			os.Exit(1)
		}
		demoRender(os.Args[2], os.Args[3])
	default:
		fmt.Printf("unknown command %q\n", os.Args[1])
		os.Exit(1)
	}
}

func demoQuery() {
	p := core.Stack(core.S("bd"), core.S("sd"))
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(2))
	b, _ := json.MarshalIndent(hapsToJSON(haps), "", "  ")
	fmt.Println(string(b))
	fmt.Printf("\n%d haps over 2 cycles\n", len(haps))
}

func demoEval(code string) {
	mini.RegisterStringParser()
	var pat core.Pattern
	if p, _, err := core.Evaluate(code, nil); err == nil {
		pat = p
	} else {
		pat = mini.Mini(code)
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	b, _ := json.MarshalIndent(hapsToJSON(haps), "", "  ")
	fmt.Println(string(b))
	fmt.Printf("\n%d haps\n", len(haps))
}

func serve(addr string) {
	srv := web.NewServer(addr)
	if err := srv.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func demoRender(outPath, code string) {
	mini.RegisterStringParser()
	var pat core.Pattern
	if p, _, err := core.Evaluate(code, nil); err == nil {
		pat = p
	} else {
		pat = mini.Mini(code)
	}
	if pat.Query == nil {
		pat = core.Pure(code)
	}
	samples, err := audio.RenderPatternAudio(pat, 4, 48000)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := audio.WriteWAV(outPath, samples, 48000); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%d samples)\n", outPath, len(samples))
}

func hapsToJSON(haps []core.Hap) []map[string]any {
	out := make([]map[string]any, len(haps))
	for i, h := range haps {
		m := map[string]any{}
		if h.Whole != nil {
			m["whole"] = h.Whole.String()
		}
		m["part"] = h.Part.String()
		m["value"] = h.Value
		out[i] = m
	}
	return out
}
