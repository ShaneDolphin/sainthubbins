//go:build goja
// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Goja-powered transpiler evaluation (temporary shim, not for production audio).

package transpiler

import (
	"fmt"

	"github.com/dop251/goja"
)

// EvaluateJSGoja evaluates JS code via goja and returns stringified result.
func EvaluateJSGoja(code string) (string, error) {
	vm := goja.New()
	v, err := vm.RunString(code)
	if err != nil {
		return "", fmt.Errorf("goja eval: %w", err)
	}
	return v.String(), nil
}
