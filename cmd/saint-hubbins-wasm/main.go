//go:build js && wasm
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// cmd/saint-hubbins-wasm — WASM entry for live console bridge.

package main

import "syscall/js"

func queryPattern(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf("missing code arg")
	}
	code := args[0].String()
	return js.ValueOf(map[string]any{
		"code":   code,
		"length": len(code),
		"haps":   []any{},
	})
}

func main() {
	js.Global().Set("saintHubbins", map[string]any{
		"version":      js.ValueOf("0.1.0-hubbins"),
		"queryPattern": js.FuncOf(queryPattern),
	})
	select {}
}
