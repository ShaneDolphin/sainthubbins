// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// cmd/saint-hubbins — native CLI for Saint Hubbins.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"codeberg.org/uzu/saint-hubbins/internal/audio"
	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
	"codeberg.org/uzu/saint-hubbins/internal/osc"
	"codeberg.org/uzu/saint-hubbins/internal/session"
	"codeberg.org/uzu/saint-hubbins/web"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: saint-hubbins <eval|serve|render|play|query> [args]")
		fmt.Println("  eval <code>        — evaluate pattern string")
		fmt.Println("  query              — demo query: Stack(s(\"bd\"), s(\"sd\"))")
		fmt.Println("  serve [addr]       — start live console server (default :8080)")
		fmt.Println("  render <out.wav> <code> — offline render to WAV")
		fmt.Println("  play <code> [host] [port] [secs] — stream to SuperDirt over OSC")
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
	case "play":
		if len(os.Args) < 3 {
			fmt.Println("play <code> [host] [port] [seconds]")
			os.Exit(1)
		}
		host, port, secs := "127.0.0.1", 57120, 8.0
		if len(os.Args) >= 4 {
			host = os.Args[3]
		}
		if len(os.Args) >= 5 {
			if v, err := strconv.Atoi(os.Args[4]); err == nil {
				port = v
			}
		}
		if len(os.Args) >= 6 {
			if v, err := strconv.ParseFloat(os.Args[5], 64); err == nil {
				secs = v
			}
		}
		if err := runPlay(os.Args[2], host, port, secs, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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

// runPlay evaluates code and streams it to SuperDirt over OSC for seconds.
// It is separate from the CLI dispatch so it can be tested without os.Exit.
func runPlay(code, host string, port int, seconds float64, out io.Writer) error {
	client := osc.New(host, port)
	defer client.Close()

	// Dial eagerly, off the scheduler's tick goroutine, so a bad host or an
	// unresolvable name is reported here rather than as a silent, dropped
	// first event once the clock starts.
	if err := client.Connect(); err != nil {
		return fmt.Errorf("play: could not reach %s:%d: %w", host, port, err)
	}

	s := session.NewSession()
	s.SetSink(&session.OSCSink{Client: client})
	s.OnError = func(err error) {
		fmt.Fprintf(out, "play: sink error: %v\n", err)
	}
	if _, err := s.Evaluate(code); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(seconds*float64(time.Second)))
	defer cancel()

	fmt.Fprintf(out, "playing %q to %s:%d for %.1fs — these go to eleven\n",
		code, host, port, seconds)
	if err := s.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	s.Stop()
	return nil
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
