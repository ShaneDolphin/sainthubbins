// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

import (
	"fmt"
	"sync"
)

// GlobalScope is the global evaluation scope — holds registered patterns and helpers.
var GlobalScope = map[string]any{}
var globalScopeMu sync.RWMutex

// UserDefinedKeys tracks user-defined vars for ClearScope.
var UserDefinedKeys = map[string]bool{}
var userDefinedMu sync.RWMutex

// ClearScope removes user-defined entries from GlobalScope.
func ClearScope() {
	userDefinedMu.Lock()
	defer userDefinedMu.Unlock()
	globalScopeMu.Lock()
	defer globalScopeMu.Unlock()
	for k := range UserDefinedKeys {
		delete(GlobalScope, k)
	}
	UserDefinedKeys = map[string]bool{}
}

// EvalScope merges modules into GlobalScope.
func EvalScope(modules ...map[string]any) {
	globalScopeMu.Lock()
	defer globalScopeMu.Unlock()
	userDefinedMu.Lock()
	defer userDefinedMu.Unlock()
	for _, mod := range modules {
		for k, v := range mod {
			GlobalScope[k] = v
		}
	}
}

// Evaluate evaluates code string into Pattern.
func Evaluate(code string, transpiler func(string) (string, error)) (Pattern, map[string]any, error) {
	meta := map[string]any{}
	if transpiler != nil {
		out, err := transpiler(code)
		if err != nil {
			return Silence(), meta, err
		}
		code = out
		meta["transpiled"] = out
	}
	trimmed := code
	globalScopeMu.RLock()
	if pat, ok := GlobalScope[trimmed]; ok {
		if p, ok := pat.(Pattern); ok {
			globalScopeMu.RUnlock()
			return p, meta, nil
		}
		if fn, ok := pat.(func() Pattern); ok {
			globalScopeMu.RUnlock()
			return fn(), meta, nil
		}
		p := Reify(pat)
		globalScopeMu.RUnlock()
		return p, meta, nil
	}
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == ';' {
		noSemi := trimmed[:len(trimmed)-1]
		if pat, ok := GlobalScope[noSemi]; ok {
			if p, ok := pat.(Pattern); ok {
				globalScopeMu.RUnlock()
				return p, meta, nil
			}
			p := Reify(pat)
			globalScopeMu.RUnlock()
			return p, meta, nil
		}
	}
	globalScopeMu.RUnlock()
	if stringParser != nil {
		miniPat := stringParser(trimmed)
		if miniPat.Query != nil {
			haps := miniPat.QueryArc(FractionFromInt(0), FractionFromInt(1))
			if len(haps) > 0 || trimmed == "~" || trimmed == "" {
				return miniPat, meta, nil
			}
			if trimmed != "" {
				return miniPat, meta, nil
			}
		}
	}
	if trimmed != "" {
		return Pure(trimmed), meta, nil
	}
	return Silence(), meta, fmt.Errorf("Evaluate: cannot evaluate %q without JS runtime (goja not yet integrated)", code)
}

// RegisterScope is helper to register a name in scope (for testing).
func RegisterScope(name string, value any) {
	globalScopeMu.Lock()
	defer globalScopeMu.Unlock()
	GlobalScope[name] = value
	userDefinedMu.Lock()
	UserDefinedKeys[name] = true
	userDefinedMu.Unlock()
}
