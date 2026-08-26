// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// A goja runtime with the pattern API bound into it.

package jsapi

import (
	"fmt"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// Evaluate runs code in a fresh VM and returns the pattern it produced.
//
// A fresh runtime per call keeps this safe under the web console's concurrent
// requests — goja runtimes are not goroutine-safe — and costs little next to
// pattern querying.
func Evaluate(code string) (core.Pattern, error) {
	mini.RegisterStringParser()
	vm := goja.New()
	if err := register(vm); err != nil {
		return core.Silence(), err
	}
	v, err := vm.RunString(code)
	if err != nil {
		return core.Silence(), fmt.Errorf("jsapi: %w", err)
	}
	return unwrap(vm, v)
}

// unwrap converts a JS result into a Pattern. A bare string is treated as
// mini-notation so `"bd sd"` works, but anything else is an error rather than
// a literal-valued hap.
func unwrap(vm *goja.Runtime, v goja.Value) (core.Pattern, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return core.Silence(), fmt.Errorf("jsapi: expression produced no value")
	}
	if obj, ok := v.Export().(*jsPattern); ok {
		return obj.pat, nil
	}
	if s, ok := v.Export().(string); ok {
		return mini.Mini(s), nil
	}
	return core.Silence(), fmt.Errorf("jsapi: expression produced %T, want a pattern", v.Export())
}
