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
	shio "codeberg.org/uzu/saint-hubbins/internal/io"
	"codeberg.org/uzu/saint-hubbins/internal/jsapi"
	"codeberg.org/uzu/saint-hubbins/internal/osc"
	"codeberg.org/uzu/saint-hubbins/internal/session"
	"codeberg.org/uzu/saint-hubbins/web"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: saint-hubbins <eval|serve|render|play|query|midi> [args]")
		fmt.Println("  eval <code>        — evaluate pattern string")
		fmt.Println("  query              — demo query: Stack(s(\"bd\"), s(\"sd\"))")
		fmt.Println("  serve [addr]       — start live console server (default :8080)")
		fmt.Println("  render <out.wav> <code> — offline render to WAV")
		fmt.Println("  play <code> [host] [port] [secs] — stream to SuperDirt over OSC")
		fmt.Println("  midi <out.mid> <code> [cycles] — render to a Standard MIDI File")
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
			fmt.Println("play <code> [host] [port] [secs]")
			os.Exit(1)
		}
		host, port, secs := "127.0.0.1", 57120, 8.0
		if len(os.Args) >= 4 {
			host = os.Args[3]
		}
		if len(os.Args) >= 5 {
			v, err := strconv.Atoi(os.Args[4])
			if err != nil {
				fmt.Fprintf(os.Stderr, "play: invalid port %q: %v\n", os.Args[4], err)
				os.Exit(1)
			}
			port = v
		}
		if len(os.Args) >= 6 {
			v, err := strconv.ParseFloat(os.Args[5], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "play: invalid seconds %q: %v\n", os.Args[5], err)
				os.Exit(1)
			}
			secs = v
		}
		if err := runPlay(os.Args[2], host, port, secs, os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "midi":
		if len(os.Args) < 4 {
			fmt.Println("midi <out.mid> <code> [cycles]")
			os.Exit(1)
		}
		cycles := 4
		if len(os.Args) >= 5 {
			if v, err := strconv.Atoi(os.Args[4]); err == nil {
				cycles = v
			}
		}
		if err := runMIDI(os.Args[3], os.Args[2], cycles); err != nil {
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

// evalCode evaluates code via jsapi.EvaluateCode (JS first, mini-notation
// fallback) and renders demoEval's pretty-printed JSON output. Kept separate
// from demoEval so the JS-vs-mini wiring can be tested without depending on
// os.Exit.
func evalCode(code string) (string, int, error) {
	pat, err := jsapi.EvaluateCode(code)
	if err != nil {
		return "", 0, err
	}
	haps := pat.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	b, _ := json.MarshalIndent(hapsToJSON(haps), "", "  ")
	return string(b), len(haps), nil
}

func demoEval(code string) {
	out, n, err := evalCode(code)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(out)
	fmt.Printf("\n%d haps\n", n)
}

func serve(addr string) {
	srv := web.NewServer(addr)
	if err := srv.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// renderPattern evaluates code via jsapi.EvaluateCode and writes the
// resulting audio to outPath, returning the sample count. Kept separate
// from demoRender so the JS-vs-mini wiring can be tested without depending
// on os.Exit.
func renderPattern(outPath, code string) (int, error) {
	pat, err := jsapi.EvaluateCode(code)
	if err != nil {
		return 0, err
	}
	samples, err := audio.RenderPatternAudio(pat, 4, 48000)
	if err != nil {
		return 0, err
	}
	if err := audio.WriteWAV(outPath, samples, 48000); err != nil {
		return 0, err
	}
	return len(samples), nil
}

func demoRender(outPath, code string) {
	n, err := renderPattern(outPath, code)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%d samples)\n", outPath, n)
}

// runMIDI evaluates code and writes it out as a Standard MIDI File.
func runMIDI(code, path string, cycles int) error {
	pat, err := jsapi.EvaluateCode(code)
	if err != nil {
		return err
	}
	f := shio.RenderMIDI(pat, cycles, 480)
	if err := f.Write(path); err != nil {
		return err
	}
	notes := f.NoteOnCount()
	fmt.Printf("wrote %s (%d cycles, %d notes)\n", path, cycles, notes)
	if notes == 0 {
		fmt.Fprintf(os.Stderr, "midi: warning: %q produced no notes — the file was written but is silent. "+
			"This usually means the pattern has no resolvable pitches: a bare numeric token like \"0 1 2 3\" "+
			"is a mini-notation string (a sample name), not a note. Use a bare drum name like \"bd\" instead "+
			"(not \"bd:3\", which sets a note number, not a percussive hit), or build the pattern with the "+
			"Go API's core.Note/core.N.\n", code)
	}
	return nil
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
