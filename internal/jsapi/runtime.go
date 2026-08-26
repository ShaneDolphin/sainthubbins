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

// evaluateTimeout bounds how long vm.RunString itself may run before it is
// interrupted — an infinite loop in the evaluated script (`while(true){}`)
// cannot hang the calling goroutine past this deadline.
//
// That is a narrower guarantee than it sounds like. Two things this timer
// does NOT cover, both confirmed by hanging a real test process on them:
//
//   - Query time. `s("bd").ply(1e9)` returns from RunString quickly — it
//     builds a lazy Pattern, it doesn't query one — so the interrupt never
//     fires. The actual work happens in QueryArc, which every caller of
//     EvaluateCode runs itself, outside this function and outside the
//     timer entirely. A web handler calling QueryArc on a pattern shaped
//     like that hangs (and can OOM) with nothing here to stop it. This is
//     pre-existing in kind, not new: `bd*1000000000` through plain
//     mini-notation has the identical shape of problem, and bounding query
//     time is a separate piece of work, not part of this fix.
//   - Native builtins. vm.Interrupt only takes effect at goja bytecode
//     instruction boundaries. A call into a native builtin that loops
//     internally in Go — `new Array(1e9).join("x")` — does not return to
//     the bytecode loop until it finishes, so the interrupt can't land
//     until after the hang is already over.
//
// What this timer does guarantee: JS whose hang is in the *evaluated
// script's own bytecode* — a loop, unbounded recursion — is interrupted.
// A package-level constant keeps the default in one place and tunable if
// it proves wrong.
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
		// No "jsapi:" prefix here — the InterruptedError branch below adds it.
		vm.Interrupt(fmt.Sprintf("evaluation exceeded %s", evaluateTimeout))
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

// unwrap converts a JS result into a Pattern via toPatternResult: a wrapped
// pattern passes through, a bare string is treated as mini-notation so
// `"bd sd"` works, and anything else — including a bare number, which
// TestEvaluateRejectsNonPatternResult has always required to be an error
// here rather than becoming a one-hap pattern — is reported rather than
// turned into a literal-valued hap.
func unwrap(vm *goja.Runtime, v goja.Value) (core.Pattern, error) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return core.Silence(), fmt.Errorf("jsapi: expression produced no value")
	}
	if p, ok := toPatternResult(v.Export()); ok {
		return p, nil
	}
	return core.Silence(), fmt.Errorf("jsapi: expression produced %T, want a pattern", v.Export())
}
