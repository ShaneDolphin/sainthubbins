// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// tools/gen-controls — generates internal/core/controls_gen.go from controls definition (Saint Hubbins controls).

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func main() {
	exeDir, _ := os.Getwd()
	candidates := []string{
		filepath.Join(exeDir, "js", "packages", "core", "controls.mjs"),
		filepath.Join(exeDir, "..", "js", "packages", "core", "controls.mjs"),
		filepath.Join(exeDir, "..", "..", "js", "packages", "core", "controls.mjs"),
	}
	var jsPath string
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			jsPath = c
			break
		}
	}
	if jsPath == "" {
		// No JS source in this repo — controls_gen.go is already generated and checked in.
		// This is intentional for Saint Hubbins: controls are maintained directly in Go.
		fmt.Fprintln(os.Stderr, "controls.mjs not found — skipping generation (checked-in file is authoritative)")
		return
	}
	data, err := os.ReadFile(jsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	text := string(data)
	lines := strings.Split(text, "\n")
	var exports []string
	var buf string
	inExp := false
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "//") {
			continue
		}
		if strings.Contains(line, "export const {") {
			inExp = true
			buf = line
			if strings.Contains(line, ";") {
				exports = append(exports, buf)
				buf = ""
				inExp = false
			}
			continue
		}
		if inExp {
			buf += " " + strings.TrimSpace(line)
			if strings.Contains(line, ";") {
				exports = append(exports, buf)
				buf = ""
				inExp = false
			}
		}
	}
	reExport := regexp.MustCompile(`\{\s*([^}]+?)\s*\}\s*=\s*register(Multi)?Control\((.*)\);`)
	validIdent := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	quoteRe := regexp.MustCompile(`'([^']+)'|"([^"]+)"`)
	leftToPrimary := map[string][]string{}

	type expData struct {
		left    []string
		isMulti bool
		args    string
	}
	var expList []expData
	for _, exp := range exports {
		m := reExport.FindStringSubmatch(exp)
		if m == nil {
			continue
		}
		leftStr := m[1]
		isMulti := m[2] == "Multi"
		args := m[3]
		parts := strings.Split(leftStr, ",")
		var left []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" && validIdent.MatchString(p) {
				left = append(left, p)
			}
		}
		if len(left) == 0 {
			continue
		}
		expList = append(expList, expData{left, isMulti, args})
	}

	parseArgs := func(argsStr string) (string, string) {
		argsStr = strings.TrimSpace(argsStr)
		if strings.HasPrefix(argsStr, "[") {
			depth := 0
			for i, ch := range argsStr {
				if ch == '[' {
					depth++
				} else if ch == ']' {
					depth--
				}
				if depth == 0 && i > 0 {
					return argsStr[:i+1], strings.Trim(strings.TrimPrefix(argsStr[i+1:], ","), " ")
				}
			}
			return argsStr, ""
		} else if strings.HasPrefix(argsStr, "'") || strings.HasPrefix(argsStr, "\"") {
			q := argsStr[0]
			end := strings.IndexByte(argsStr[1:], q)
			if end >= 0 {
				end += 1
				return argsStr[:end+1], strings.Trim(strings.TrimPrefix(argsStr[end+1:], ","), " ")
			}
		}
		parts := strings.SplitN(argsStr, ",", 2)
		if len(parts) == 1 {
			return parts[0], ""
		}
		return parts[0], strings.TrimSpace(parts[1])
	}
	extractQuotes := func(s string) []string {
		matches := quoteRe.FindAllStringSubmatch(s, -1)
		var out []string
		for _, m := range matches {
			if m[1] != "" {
				out = append(out, m[1])
			} else if m[2] != "" {
				out = append(out, m[2])
			}
		}
		return out
	}

	for _, e := range expList {
		if !e.isMulti {
			first, _ := parseArgs(e.args)
			var primary []string
			first = strings.TrimSpace(first)
			if strings.HasPrefix(first, "[") {
				primary = extractQuotes(first)
			} else {
				primary = extractQuotes(first)
			}
			for _, js := range e.left {
				leftToPrimary[js] = primary
			}
		} else {
			argsStr := strings.TrimSpace(e.args)
			var namesArg string
			var rest string
			if strings.HasPrefix(argsStr, "[") {
				depth := 0
				idx := -1
				for i, ch := range argsStr {
					if ch == '[' {
						depth++
					} else if ch == ']' {
						depth--
					}
					if depth == 0 {
						idx = i
						break
					}
				}
				namesArg = argsStr[:idx+1]
				rest = strings.Trim(strings.TrimPrefix(argsStr[idx+1:], ","), " ")
				_ = rest
			} else {
				m := quoteRe.FindString(argsStr)
				namesArg = m
				rest = strings.Trim(strings.TrimPrefix(argsStr[len(m):], ","), " ")
				_ = rest
			}
			var namesList []string
			namesArg = strings.TrimSpace(namesArg)
			if strings.HasPrefix(namesArg, "[") {
				namesList = extractQuotes(namesArg)
			} else {
				namesList = extractQuotes(namesArg)
			}
			for _, js := range e.left {
				reNum := regexp.MustCompile(`^(.*?)(\d+)$`)
				m2 := reNum.FindStringSubmatch(js)
				var primary []string
				if m2 != nil {
					base := m2[1]
					hasBase := false
					for _, other := range e.left {
						if other == base {
							hasBase = true
							break
						}
					}
					if strings.HasSuffix(js, "1") && hasBase {
						primary = namesList
					} else {
						numStr := m2[2]
						if numStr == "1" {
							if hasBase {
								primary = namesList
							} else {
								var tmp []string
								for _, n := range namesList {
									tmp = append(tmp, n+"1")
								}
								primary = tmp
							}
						} else {
							var tmp []string
							for _, n := range namesList {
								tmp = append(tmp, n+numStr)
							}
							primary = tmp
						}
					}
				} else {
					primary = namesList
				}
				leftToPrimary[js] = primary
			}
		}
	}

	jsToGo := func(js string) string {
		if js == "" {
			return ""
		}
		return strings.ToUpper(js[:1]) + js[1:]
	}
	valid := map[string]string{}
	for js := range leftToPrimary {
		if validIdent.MatchString(js) {
			valid[js] = jsToGo(js)
		}
	}
	existing := map[string]bool{
		"S": true, "Sound": true, "N": true, "Note": true, "Gain": true, "Velocity": true, "Vel": true,
		"Cutoff": true, "Lpf": true, "Resonance": true, "Lpq": true, "Delay": true, "Room": true, "Size": true,
		"Pan": true, "Speed": true, "Begin": true, "End": true, "Bank": true, "Orbit": true, "Octave": true,
		"Coarse": true, "CRush": true, "Shape": true, "Distort": true, "Cut": true, "Legato": true, "Sustain": true,
		"Release": true, "Attack": true, "Decay": true, "Vowel": true, "Hpf": true, "Hpq": true, "Bpf": true, "Bpq": true,
		"Freq": true, "Up": true, "Off": true, "Crush": true,
	}
	newVars := map[string][]string{}
	overrideVars := map[string][]string{}
	for js, goName := range valid {
		primary := leftToPrimary[js]
		if len(primary) == 0 {
			continue
		}
		if existing[goName] {
			if len(primary) != 1 || !strings.EqualFold(primary[0], goName) {
				overrideVars[js] = primary
			}
		} else {
			newVars[js] = primary
		}
	}
	type kv struct{ js string; primary []string; goName string }
	var sorted []kv
	for js, primary := range newVars {
		sorted = append(sorted, kv{js, primary, valid[js]})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].goName < sorted[j].goName })

	header := `// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Controls — generated by tools/gen-controls — DO NOT EDIT.

package core
`
	var out string
	out += header + "\nvar (\n"
	for _, kv := range sorted {
		args := ""
		for i, p := range kv.primary {
			if i > 0 {
				args += ", "
			}
			args += fmt.Sprintf("%q", p)
		}
		out += fmt.Sprintf("\t%s = createParam(%s)\n", kv.goName, args)
	}
	out += ")\n"
	if len(overrideVars) > 0 {
		out += "\nfunc init() {\n"
		grouped := map[string][]string{}
		primaryForGroup := map[string][]string{}
		for js, primary := range overrideVars {
			key := strings.Join(primary, "|")
			grouped[key] = append(grouped[key], js)
			primaryForGroup[key] = primary
		}
		var keys []string
		for k := range grouped {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			jss := grouped[k]
			primary := primaryForGroup[k]
			args := ""
			for i, p := range primary {
				if i > 0 {
					args += ", "
				}
				args += fmt.Sprintf("%q", p)
			}
			if len(jss) == 1 {
				out += fmt.Sprintf("\t%s = createParam(%s)\n", valid[jss[0]], args)
			} else {
				sort.Slice(jss, func(a, b int) bool { return valid[jss[a]] < valid[jss[b]] })
				tmp := valid[jss[0]]
				out += fmt.Sprintf("\t%s = createParam(%s)\n", tmp, args)
				for _, js2 := range jss[1:] {
					out += fmt.Sprintf("\t%s = %s\n", valid[js2], tmp)
				}
			}
		}
		out += "}\n"
	}

	goRootCandidates := []string{
		filepath.Join(exeDir, "internal", "core", "controls_gen.go"),
		filepath.Join(exeDir, "..", "internal", "core", "controls_gen.go"),
		filepath.Join(exeDir, "..", "..", "internal", "core", "controls_gen.go"),
	}
	var outPath string
	for _, c := range goRootCandidates {
		dir := filepath.Dir(c)
		if _, err := os.Stat(dir); err == nil {
			outPath = c
			break
		}
	}
	if outPath == "" {
		outPath = "controls_gen.go"
	}
	if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("generated %s (%d new, %d overrides)\n", outPath, len(newVars), len(overrideVars))
}
