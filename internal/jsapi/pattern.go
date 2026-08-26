// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// The JS-facing pattern object.

package jsapi

import (
	"github.com/dop251/goja"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

// jsPattern wraps a core.Pattern so JS can hold and chain it. Methods are
// attached in registry.go rather than declared here, so adding a transform is
// a table entry rather than a new method.
type jsPattern struct {
	pat core.Pattern
}

// newJSPattern builds the JS object for a pattern: the wrapper plus every
// chainable method from the tables in registry.go.
func newJSPattern(vm *goja.Runtime, p core.Pattern) *goja.Object {
	jp := &jsPattern{pat: p}
	obj := vm.ToValue(jp).(*goja.Object)
	attachMethods(vm, obj, jp)
	return obj
}
