// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package transpiler

import "testing"

func TestTranspileStripImports(t *testing.T) {
	code := "import { s } from 'hubbins'; s(\"bd sd\")"
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	got := res.Output
	// Should strip import line but keep s("bd sd")
	if len(got) == 0 {
		t.Fatalf("empty")
	}
	// Should not contain "import"
	for i := 0; i+6 <= len(got); i++ {
		if got[i:i+6] == "import" {
			t.Fatalf("still contains import: %q", got)
		}
	}
}

func TestTranspilePlain(t *testing.T) {
	code := "s(\"bd sd\").slow(2)"
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	got := res.Output
	if len(got) == 0 {
		t.Fatalf("empty")
	}
	t.Logf("got %q", got)
}
