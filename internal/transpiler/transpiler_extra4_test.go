// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package transpiler

import "testing"

func TestTranspileMultipleStrings(t *testing.T) {
	code := `s("bd sd") ; n("0 1")`
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	// Should contain two m(' calls
	count := 0
	for i := 0; i+1 < len(res.Output); i++ {
		if res.Output[i] == 'm' && res.Output[i+1] == '(' {
			count++
		}
	}
	if count < 2 {
		t.Fatalf("expected 2 m( calls got %d output %q", count, res.Output)
	}
	if len(res.MiniLocations) < 2 {
		t.Fatalf("locations %v", res.MiniLocations)
	}
}

func TestTranspileNoString(t *testing.T) {
	code := `42 + 1`
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if res.Output != code {
		t.Logf("output %q", res.Output)
	}
	if len(res.MiniLocations) != 0 {
		t.Fatalf("no string should have 0 locations got %v", res.MiniLocations)
	}
}
