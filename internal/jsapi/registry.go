// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Tables binding JS names to engine operations.

package jsapi

import (
	"fmt"
	"math"

	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
	"codeberg.org/uzu/saint-hubbins/internal/mini"
)

// controls are the constructors that turn a value into a control pattern.
// The JS name is the key; the Go constructor is the value.
var controls = map[string]func(any) core.Pattern{
	"s": core.S, "sound": core.S, "note": core.Note, "n": core.N,
	"gain": core.Gain, "cutoff": core.Cutoff, "lpf": core.Lpf,
	"pan": core.Pan, "room": core.Room, "speed": core.Speed,
	"attack": core.Attack, "release": core.Release, "shape": core.Shape,
}

// toPattern coerces a JS argument into a Pattern: a wrapped pattern passes
// through, a string is mini-notation, a number is a constant.
func toPattern(v any) core.Pattern {
	switch x := v.(type) {
	case *jsPattern:
		return x.pat
	case string:
		return mini.Mini(x)
	case float64, int, int64:
		return core.Pure(x)
	}
	return core.Silence()
}

// register installs every global into the VM.
func register(vm *goja.Runtime) error {
	wrap := func(p core.Pattern) goja.Value { return vm.ToValue(newJSPattern(vm, p)) }

	for name, ctor := range controls {
		ctor := ctor
		if err := vm.Set(name, func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				return wrap(core.Silence())
			}
			arg := call.Argument(0).Export()
			// A string argument is mini-notation, so s("bd sd") is a sequence.
			if str, ok := arg.(string); ok {
				return wrap(ctor(mini.Mini(str)))
			}
			if jp, ok := arg.(*jsPattern); ok {
				return wrap(ctor(jp.pat))
			}
			return wrap(ctor(arg))
		}); err != nil {
			return err
		}
	}
	return nil
}

// unaryOps take no arguments.
var unaryOps = map[string]func(core.Pattern) core.Pattern{
	"rev":        func(p core.Pattern) core.Pattern { return p.Rev() },
	"palindrome": func(p core.Pattern) core.Pattern { return p.Palindrome() },
	"degrade":    func(p core.Pattern) core.Pattern { return p.Degrade() },
	"hush":       func(core.Pattern) core.Pattern { return core.Silence() },
}

// numericOps take a single number.
var numericOps = map[string]func(core.Pattern, float64) core.Pattern{
	"fast":      func(p core.Pattern, v float64) core.Pattern { return p.FastF(core.FractionFromFloat(v)) },
	"slow":      func(p core.Pattern, v float64) core.Pattern { return p.SlowF(core.FractionFromFloat(v)) },
	"ply":       func(p core.Pattern, v float64) core.Pattern { return p.Ply(int(v)) },
	"segment":   func(p core.Pattern, v float64) core.Pattern { return p.Segment(v) },
	"late":      func(p core.Pattern, v float64) core.Pattern { return p.Late(v) },
	"early":     func(p core.Pattern, v float64) core.Pattern { return p.Early(v) },
	"degradeBy": func(p core.Pattern, v float64) core.Pattern { return p.DegradeBy(v) },
	"add":       func(p core.Pattern, v float64) core.Pattern { return p.Add(v) },
}

// numericArgRules constrains a numericOps argument beyond "a finite number",
// for ops where a zero or negative value is not just unusual but actively
// breaks the underlying engine call rather than doing something coherent:
//
//   - fast/slow/segment divide by (or otherwise scale by) the argument
//     (FastF, SlowF, Segment). Zero is a division by zero — SlowF computes
//     1/frac eagerly, so .slow(0) panics synchronously inside the JS call
//     itself ("Fraction.Div: division by zero"), outside any recover, and
//     crashes the process exactly like the ±Infinity case below; .fast(0)
//     panics lazily at Query time instead, caught by the recover in
//     QueryArc (internal/core/pattern.go), which prints "query panic: ..."
//     to stdout and silently returns zero haps — no crash, but the same
//     silent-failure shape this whole plan exists to remove. A negative
//     factor doesn't panic anywhere, but empirically (verified by hand)
//     produces zero haps just as silently, for the same underlying reason:
//     this engine's time-scaling machinery assumes a positive factor.
//     Segment already clamps a non-positive rate to 1 inside
//     core.Pattern.Segment rather than crashing, but that's a *silent*
//     override of what the caller asked for — the same failure shape,
//     just without the crash or the log line.
//   - ply repeats each event by its argument via SqueezeJoin. Zero is
//     coherent on its own terms (repeat something zero times is nothing,
//     and core.Pattern.Ply already resolves it to Silence() deliberately,
//     without panicking or printing anything) so it is left alone. A
//     negative repeat count has no coherent reading and — like fast/slow/
//     segment — empirically just produces zero haps silently.
//
// late, early, degradeBy and add are deliberately absent: a negative
// offset is just an offset in the other direction (verified: .early(-0.25)
// and .late(-0.25) both produce the expected event count, not zero
// haps or a panic); degradeBy(0)/degradeBy(negative) both mean "never
// drop" and return every hap, which is exactly the coherent boundary
// behaviour DegradeBy's own doc comment implies; add's whole point is
// signed arithmetic, so a negative addend is the normal case, not an
// edge case.
var numericArgRules = map[string]func(v float64) error{
	"fast":    requirePositive,
	"slow":    requirePositive,
	"segment": requirePositive,
	"ply":     requireNonNegative,
}

func requirePositive(v float64) error {
	if v <= 0 {
		return fmt.Errorf("must be a positive number, got %v", v)
	}
	return nil
}

func requireNonNegative(v float64) error {
	if v < 0 {
		return fmt.Errorf("must not be negative, got %v", v)
	}
	return nil
}

// normalizeNumber ensures a control value coming from goja is a Go float64
// rather than int64: goja's Export() returns int64 for any JS number with
// no fractional part — `.cutoff(800)` and even `.cutoff(800.0)` both export
// as int64(800), while `.gain(0.5)` exports as float64 — so a control bag's
// numeric values would otherwise vary in Go type depending on whether the
// JS literal happened to be whole, rather than staying float64 consistently.
func normalizeNumber(v any) any {
	switch n := v.(type) {
	case int64:
		return float64(n)
	case int:
		return float64(n)
	case float32:
		return float64(n)
	}
	return v
}

// patternFromJSValue converts a JS value into a Pattern the same way
// unwrap does for a top-level Evaluate result: a wrapped pattern passes
// through, a bare string is mini-notation. It reports false rather than
// erroring for anything else, because its only caller (the every/off
// callback re-entry below) runs at Query time, outside any goja call frame
// — there is no way to turn that into a Go error the original Evaluate
// caller would ever see.
func patternFromJSValue(v goja.Value) (core.Pattern, bool) {
	if v == nil || goja.IsUndefined(v) || goja.IsNull(v) {
		return core.Silence(), false
	}
	switch x := v.Export().(type) {
	case *jsPattern:
		return x.pat, true
	case string:
		return mini.Mini(x), true
	}
	return core.Silence(), false
}

// attachMethods installs every chainable method on a wrapped pattern's JS
// object: the two op tables above, the controls map reused as setters, and
// the handful of methods (euclid, every) that don't fit either table's
// shape.
//
// Every method here validates its own arguments and raises a JS TypeError
// (via panic(vm.NewTypeError(...)), the same mechanism goja's own builtins
// use) rather than silently coercing a bad argument into 0/NaN or returning
// the pattern unchanged. A panic raised while goja is executing JS is
// recovered by the VM and turned into the error Evaluate returns — it is
// not a crash. An unrecognized method name (`.nosuchmethod()`) needs no
// special handling: goja already raises "Object has no member" on its own
// for any name attachMethods never Sets.
//
// Methods are Set on a fresh plain object used as obj's prototype, not on
// obj itself. obj wraps *jsPattern via reflection (vm.ToValue(jp) in
// newJSPattern) so that Export() can recover the underlying Pattern — but
// goja treats that reflected wrapper as a "host object", which rejects
// Object.Set for any new property ("Cannot assign to property ... of a host
// object", verified empirically; the error was there to see, just easy to
// discard along with wrap's return value). SetPrototype, unlike Set, works
// on a host object, and property lookup still walks the prototype chain, so
// `p.fast(2)` resolves through proto to the closures below.
func attachMethods(vm *goja.Runtime, obj *goja.Object, jp *jsPattern) {
	proto := vm.NewObject()
	wrap := func(p core.Pattern) goja.Value { return vm.ToValue(newJSPattern(vm, p)) }

	// requireFiniteNumber extracts and validates the argument at idx,
	// raising a TypeError for:
	//   - a missing argument (fewer than idx+1 arguments given);
	//   - one that doesn't convert to a number (ToFloat is NaN — e.g. a
	//     string like "banana");
	//   - ±Infinity — core.FractionFromFloat itself panics on a non-finite
	//     float ("invalid float +Inf"), synchronously, inside this very JS
	//     call; goja only recovers panics of its own error types, so an
	//     unguarded .fast(Infinity) crashes the whole process (verified),
	//     not just this one Evaluate call. euclid's steps argument has the
	//     same crash shape one step removed: ToInteger() maps +Infinity to
	//     math.MaxInt64, which reaches make([]int, steps) in
	//     internal/core/euclid.go and panics with "makeslice: len out of
	//     range" — a *runtime.Error goja doesn't recognize either, so it
	//     repanics just the same. This is the single choke point every
	//     numeric argument in this file goes through, specifically so a
	//     fix here can't land in numericOps (which review confirmed after
	//     the first round) while leaving euclid/every's hand-written
	//     argument handling uncovered (which review confirmed the second
	//     round — euclid's steps and every's n take numeric arguments but
	//     sit outside the numericOps table, so this check never reached
	//     them until now).
	// A JS boolean argument (.fast(true)) is deliberately NOT rejected
	// here: ToFloat(true) is a well-defined, finite 1.0, matching ordinary
	// JS arithmetic (`true + 1 === 2`) rather than this engine inventing a
	// stricter rule departure just for chained methods.
	requireFiniteNumber := func(name string, call goja.FunctionCall, idx int) float64 {
		if len(call.Arguments) <= idx {
			panic(vm.NewTypeError("%s: requires a numeric argument", name))
		}
		v := call.Argument(idx).ToFloat()
		if math.IsNaN(v) {
			panic(vm.NewTypeError("%s: argument %q is not a number", name, call.Argument(idx).String()))
		}
		if math.IsInf(v, 0) {
			panic(vm.NewTypeError("%s: argument must be finite, got %v", name, v))
		}
		return v
	}

	// requireNumber is requireFiniteNumber for a numericOps entry's single
	// argument (always at index 0), plus that op's own domain rule from
	// numericArgRules, if it has one (see that map's comment).
	requireNumber := func(name string, call goja.FunctionCall) float64 {
		v := requireFiniteNumber(name, call, 0)
		if rule, ok := numericArgRules[name]; ok {
			if err := rule(v); err != nil {
				panic(vm.NewTypeError("%s: %v", name, err))
			}
		}
		return v
	}

	for name, op := range unaryOps {
		name, op := name, op
		_ = proto.Set(name, func(goja.FunctionCall) goja.Value { return wrap(op(jp.pat)) })
	}

	for name, op := range numericOps {
		name, op := name, op
		_ = proto.Set(name, func(call goja.FunctionCall) goja.Value {
			return wrap(op(jp.pat, requireNumber(name, call)))
		})
	}

	// Controls double as setters when called on a pattern: .gain(0.5) merges
	// a gain control into every event. A string argument is mini-notation
	// (matching the top-level constructors in register()); a wrapped
	// pattern argument is unwrapped to a core.Pattern so a modulated control
	// like .gain(sine) reaches createParam's own Pattern branch instead of
	// being embedded as an opaque *jsPattern value.
	for name, ctor := range controls {
		name, ctor := name, ctor
		_ = proto.Set(name, func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) == 0 {
				panic(vm.NewTypeError("%s: requires an argument", name))
			}
			switch arg := call.Argument(0).Export().(type) {
			case string:
				return wrap(jp.pat.Set(ctor(mini.Mini(arg))))
			case *jsPattern:
				return wrap(jp.pat.Set(ctor(arg.pat)))
			default:
				return wrap(jp.pat.Set(ctor(normalizeNumber(arg))))
			}
		})
	}

	// euclid takes two arguments, so it is not in the numeric table.
	//
	// euclidMaxSteps bounds the steps argument: core.Pattern.Euclid ends up
	// at make([]int, steps) in internal/core/euclid.go, sized by steps
	// alone (pulses is only ever compared, never used as a length). A
	// non-finite steps is already rejected by requireFiniteNumber above,
	// but a merely huge *finite* value (.euclid(3, 1_000_000_000), ~8GB at
	// 8 bytes/int) clears that check and allocates for real. Evaluate's
	// 5-second vm.Interrupt timer cannot help here: it only preempts
	// between JS bytecode instructions, and by the time steps reaches this
	// handler the allocation is a native Go call already in flight, not JS
	// the VM can interrupt — validating the argument before the call is
	// the only available defense. 1024 is far beyond any musically
	// meaningful Euclidean rhythm (steps is typically single digits to a
	// few dozen) while staying trivially cheap to allocate even at the
	// boundary.
	const euclidMaxSteps = 1024
	_ = proto.Set("euclid", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("euclid: requires two arguments (pulses, steps)"))
		}
		pulses := requireFiniteNumber("euclid", call, 0)
		steps := requireFiniteNumber("euclid", call, 1)
		if steps < 0 || steps > euclidMaxSteps {
			panic(vm.NewTypeError("euclid: steps must be between 0 and %d, got %v", euclidMaxSteps, steps))
		}
		return wrap(jp.pat.Euclid(int(pulses), int(steps)))
	})

	// every takes a cycle count and a callback, re-entered as a pattern op.
	// core.Pattern.Every already reads the cycle number off the query span
	// and calls SplitQueries() itself (pattern_time.go:193/198), so a wide
	// query is safe here without any extra splitting on our part.
	_ = proto.Set("every", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("every: requires two arguments (n, fn)"))
		}
		nFloat := requireFiniteNumber("every", call, 0)
		if nFloat <= 0 {
			// core.Pattern.Every treats n <= 0 as "return the pattern
			// unchanged" (pattern_time.go) — silent from this binding's
			// point of view. NaN/Infinity/a non-numeric first argument are
			// already rejected above by requireFiniteNumber; this catches
			// the remaining silent case review found: 0, a negative
			// number, or null (which converts to 0 via ToFloat).
			panic(vm.NewTypeError("every: n must be a positive number, got %v", nFloat))
		}
		n := int(nFloat)
		fn, ok := goja.AssertFunction(call.Argument(1))
		if !ok {
			panic(vm.NewTypeError("every: second argument must be a function"))
		}
		return wrap(jp.pat.Every(n, func(p core.Pattern) core.Pattern {
			// This callback is invoked lazily whenever the resulting
			// pattern is queried — which may be long after Evaluate has
			// returned, and outside any goja call frame Evaluate's caller
			// is inside of. Neither a Go error from calling fn nor a
			// non-pattern return value can be turned into the Go error
			// Evaluate returns; panicking here would crash the host process
			// instead (this codebase has no recover() anywhere). Falling
			// back to Silence for just that invocation is the least-silent
			// safe option: a broken per-cycle transform shows up as an
			// audible dropout in the rendered cycle, not as a subtly wrong
			// or duplicated one.
			res, err := fn(goja.Undefined(), vm.ToValue(newJSPattern(vm, p)))
			if err != nil {
				return core.Silence()
			}
			if pat, ok := patternFromJSValue(res); ok {
				return pat
			}
			return core.Silence()
		}))
	})

	_ = obj.SetPrototype(proto)
}
