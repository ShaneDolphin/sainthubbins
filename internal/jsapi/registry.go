// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Tables binding JS names to engine operations.

package jsapi

import (
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

// attachMethods is filled in by Task 2. For now it exists so Task 1 compiles.
func attachMethods(vm *goja.Runtime, obj *goja.Object, jp *jsPattern) {}
