// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// A goja runtime with the pattern API bound into it.

package jsapi

import (
	"errors"
	"fmt"
	"time"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// evaluateTimeout bounds how long a single Evaluate call may run JS before
// it is interrupted. Runaway code — an infinite loop, a pattern that never
// stops recursing — would otherwise hang the calling goroutine forever;
// under Task 4 that goroutine is a web request handler, so a hang there
// becomes a downed console rather than just a stuck test. A package-level
// constant keeps the default in one place and tunable if it proves wrong.
const evaluateTimeout = 5 * time.Second

// Evaluate runs code in a fresh VM and returns the pattern it produced.
//
// A fresh runtime per call keeps the goja VM itself safe under the web
// console's concurrent requests — goja runtimes are not goroutine-safe. The
// mini-notation string-parser hook registered below is a separate, package
// -level piece of shared state (see core.SetStringParser); it is guarded by
// its own mutex, so concurrent Evaluate calls are safe on that count too.
//
// A timer interrupts the VM after evaluateTimeout so JS that never returns
// (`while(true){}`) cannot hang the calling goroutine indefinitely.
func Evaluate(code string) (core.Pattern, error) {
	mini.RegisterStringParser()
	vm := goja.New()
	if err := register(vm); err != nil {
		return core.Silence(), err
	}

	timer := time.AfterFunc(evaluateTimeout, func() {
		vm.Interrupt(fmt.Sprintf("jsapi: evaluation exceeded %s", evaluateTimeout))
	})
	defer timer.Stop()

	v, err := vm.RunString(code)
	if err != nil {
		var interrupted *goja.InterruptedError
		if errors.As(err, &interrupted) {
			return core.Silence(), fmt.Errorf("jsapi: %v", interrupted.Value())
		}
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
