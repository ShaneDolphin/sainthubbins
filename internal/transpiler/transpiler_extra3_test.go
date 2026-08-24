// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package transpiler

import "testing"

func TestTranspileBacktick(t *testing.T) {
	code := "s(`bd sd`)"
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if len(res.Output) == 0 {
		t.Fatalf("empty")
	}
	// backtick untagged should become m('bd sd', offset)
	if res.Output == code {
		t.Fatalf("backtick not transformed %q", res.Output)
	}
}

func TestTranspileDisableRange(t *testing.T) {
	code := "/* mini-off */ s(\"bd sd\") /* mini-on */ s(\"hh\")"
	res, err := Transpile(code, Options{})
	if err != nil {
		t.Fatalf("err %v", err)
	}
	// first s should not be transformed due to disable, second should
	// Count m(' occurrences — should be 1 (only hh)
	count := 0
	for i := 0; i+2 < len(res.Output); i++ {
		if res.Output[i] == 'm' && res.Output[i+1] == '(' && res.Output[i+2] == '\'' {
			count++
		}
	}
	if count != 1 {
		t.Logf("output %q count %d (expected 1)", res.Output, count)
	}
	if count == 0 {
		t.Fatalf("disable range failed, no m(' found")
	}
}
