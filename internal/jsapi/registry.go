// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Tables binding JS names to engine operations.

package jsapi

import (
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

	// requireNumber extracts a numeric argument for a numericOps entry,
	// raising a TypeError for a missing argument (ToFloat on the resulting
	// undefined is NaN) or one that doesn't convert to a number (e.g. a
	// string like "banana" also converts to NaN) rather than silently
	// building a fraction out of NaN.
	requireNumber := func(name string, call goja.FunctionCall) float64 {
		if len(call.Arguments) == 0 {
			panic(vm.NewTypeError("%s: requires a numeric argument", name))
		}
		v := call.Argument(0).ToFloat()
		if math.IsNaN(v) {
			panic(vm.NewTypeError("%s: argument %q is not a number", name, call.Argument(0).String()))
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
	_ = proto.Set("euclid", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("euclid: requires two arguments (pulses, steps)"))
		}
		return wrap(jp.pat.Euclid(int(call.Argument(0).ToInteger()), int(call.Argument(1).ToInteger())))
	})

	// every takes a cycle count and a callback, re-entered as a pattern op.
	// core.Pattern.Every already reads the cycle number off the query span
	// and calls SplitQueries() itself (pattern_time.go:193/198), so a wide
	// query is safe here without any extra splitting on our part.
	_ = proto.Set("every", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			panic(vm.NewTypeError("every: requires two arguments (n, fn)"))
		}
		n := int(call.Argument(0).ToInteger())
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
