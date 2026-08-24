// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package transpiler

import "testing"

func TestTranspileImportStripped(t *testing.T) {
	code := "import { s } from \"@hubbins/core\";\ns(\"bd sd\")"
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("transpile err %v", err)
	}
	// Import should be stripped (replaced with whitespace) and not contain 'import'
	if contains(res.Output, "import") {
		t.Fatalf("import not stripped: %s", res.Output)
	}
	if !contains(res.Output, "m('bd sd'") {
		t.Fatalf("mini not transformed: %s", res.Output)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
