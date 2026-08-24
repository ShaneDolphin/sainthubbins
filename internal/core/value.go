// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/value.mjs (64 LOC) — Value functor and unionWithObj.

package core

import (
	"fmt"
	"log"
)

// UnionWithObj merges two control bags, applying func to common keys.
func UnionWithObj(a, b map[string]any, fn func(any, any) any) map[string]any {
	if b != nil {
		if _, hasValue := b["value"]; hasValue && len(b) == 1 {
			log.Printf("[warn]: Can't do arithmetic on control pattern.")
			return a
		}
	}
	// common keys
	common := []string{}
	for k := range a {
		if _, ok := b[k]; ok {
			common = append(common, k)
		}
	}
	out := map[string]any{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	for _, k := range common {
		out[k] = fn(a[k], b[k])
	}
	return out
}

// Value is a Maybe-like functor for combined values
type Value struct {
	Val any
}

func ValOf(x any) Value { return Value{Val: x} }

func (v Value) IsNothing() bool { return v.Val == nil }

func (v Value) Map(f func(any) any) Value {
	if v.IsNothing() {
		return v
	}
	return ValOf(f(v.Val))
}

func (v Value) Ap(other any) Value {
	switch o := other.(type) {
	case Value:
		return o.Map(func(x any) any {
			if fn, ok := v.Val.(func(any) any); ok {
				return fn(x)
			}
			return x
		})
	default:
		// treat other as Value
		ov := ValOf(other)
		return ov.Map(func(x any) any {
			if fn, ok := v.Val.(func(any) any); ok {
				return fn(x)
			}
			return x
		})
	}
}

func (v Value) UnionWith(other any, fn func(any, any) any) Value {
	ov := ValOf(other)
	if ov.IsNothing() || v.IsNothing() {
		return v
	}
	av, aok := v.Val.(map[string]any)
	bv, bok := ov.Val.(map[string]any)
	if !aok || !bok {
		panic(fmt.Sprintf("unionWith: expected objects got %T %T", v.Val, ov.Val))
	}
	return ValOf(UnionWithObj(av, bv, fn))
}

func Mul(a, b any) any {
	af := toFloat(a)
	bf := toFloat(b)
	return af * bf
}
