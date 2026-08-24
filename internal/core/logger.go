// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.

package core

import "log"

func Logger(args ...any) { log.Println(args...) }
func ErrorLogger(err error, ctx string) { log.Printf("[%s] %v", ctx, err) }
