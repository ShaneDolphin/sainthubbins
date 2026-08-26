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

// evaluateTimeout bounds how long a single piece of running JS — the
// top-level script in Evaluate, or one every()/off()-style callback
// invocation in registry.go — may run before it is interrupted. Both call
// armInterruptTimeout, so both get their own independent deadline: a
// callback re-entered at query time, arbitrarily long after Evaluate's own
// timer already fired and stopped, is not covered by that timer and needs
// one of its own. Without this, `s("bd*4").every(2, x => { while(true){} })`
// hangs forever and leaks a goroutine on every request that triggers it —
// confirmed by leaving a real test process running past 15s on exactly that
// input, where a bare top-level `while(true){}` already returns by 6s.
//
// That is still a narrower guarantee than it sounds like. Two things
// neither timer covers, both confirmed by hanging a real test process on
// them:
//
//   - Query time that isn't inside a callback. `s("bd").ply(1e9)` returns
//     from RunString quickly — it builds a lazy Pattern, it doesn't query
//     one — so no interrupt ever fires for it. The actual work happens in
//     QueryArc, which every caller of EvaluateCode runs itself, outside
//     Evaluate and outside any timer entirely; the pattern algebra inside
//     QueryArc is plain Go code, not JS, so vm.Interrupt cannot reach it
//     even when the pattern in question was built by an every() callback
//     that already returned. A web handler calling QueryArc on a pattern
//     shaped like that hangs (and can OOM) with nothing here to stop it.
//     This is pre-existing in kind, not new: `bd*1000000000` through plain
//     mini-notation has the identical shape of problem, and bounding query
//     time is a separate piece of work, not part of this fix.
//   - Native builtins. vm.Interrupt only takes effect at goja bytecode
//     instruction boundaries. A call into a native builtin that loops
//     internally in Go — `new Array(1e9).join("x")` — does not return to
//     the bytecode loop until it finishes, so the interrupt can't land
//     until after the hang is already over. This applies equally whether
//     the native builtin is called from the top-level script or from
//     inside a callback.
//
// What these timers do guarantee: JS whose hang is in *evaluated script's
// own bytecode* — a loop, unbounded recursion — is interrupted, whether
// that script is the top-level Evaluate call or a callback invoked later
// at query time. A package-level constant keeps the default in one place
// and tunable if it proves wrong.
const evaluateTimeout = 5 * time.Second

// armInterruptTimeout starts a timer that interrupts vm after
// evaluateTimeout and returns a function that must be deferred immediately.
// That function stops the timer (so it can't fire after the protected call
// already returned) and clears any interrupt that did fire — ClearInterrupt
// is required, not optional, because Interrupt() called while the runtime
// isn't actively running JS is queued rather than dropped: if this same vm
// runs more JS afterward (a later every() callback invocation, in the case
// this exists for) without ClearInterrupt(), that queued interrupt fires
// against completely unrelated code the instant it starts.
//
// reason names what timed out, folded into the message an InterruptedError
// carries when the caller unwraps it (see Evaluate and every's callback in
// registry.go).
func armInterruptTimeout(vm *goja.Runtime, reason string) func() {
	timer := time.AfterFunc(evaluateTimeout, func() {
		// No "jsapi:" prefix here — each caller's InterruptedError handling
		// adds it.
		vm.Interrupt(fmt.Sprintf("%s exceeded %s", reason, evaluateTimeout))
	})
	return func() {
		timer.Stop()
		vm.ClearInterrupt()
	}
}

// Evaluate runs code in a fresh VM and returns the pattern it produced.
//
// A fresh runtime per call keeps the goja VM itself safe under the web
// console's concurrent requests — goja runtimes are not goroutine-safe. The
// mini-notation string-parser hook registered below is a separate, package
// -level piece of shared state (see core.SetStringParser); it is guarded by
// its own mutex, so concurrent Evaluate calls are safe on that count too.
//
// armInterruptTimeout interrupts the VM after evaluateTimeout so JS that
// never returns (`while(true){}`) cannot hang the calling goroutine
// indefinitely. See evaluateTimeout's doc comment for what this does and
// does not cover — in particular, a callback invoked later at query time
// (every()/off()) is outside this call's timer and arms its own.
func Evaluate(code string) (core.Pattern, error) {
	mini.RegisterStringParser()
	vm := goja.New()
	if err := register(vm); err != nil {
		return core.Silence(), err
	}

	defer armInterruptTimeout(vm, "evaluation")()

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
