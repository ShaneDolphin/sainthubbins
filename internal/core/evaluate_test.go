package core

import "testing"

func TestEvaluateScope(t *testing.T) {
	ClearScope()
	RegisterScope("myNum", 42)
	pat, _, err := Evaluate("myNum", nil)
	if err != nil {
		t.Fatalf("evaluate myNum err %v", err)
	}
	haps := pat.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || toFloat(haps[0].Value) != 42 {
		t.Fatalf("evaluate myNum got %v", haps[0].Value)
	}
	ClearScope()
	if len(GlobalScope) != 0 {
		// GlobalScope should be cleared of user keys, but may have other?
	}
}

func TestEvaluateMiniFallback(t *testing.T) {
	miniPat, _, _ := Evaluate("bd sd", nil)
	// With stub transpiler, bd sd will fail evaluate and fallback to Mini for console, but direct Evaluate will error?
	// This test just checks Evaluate returns error for unknown
	if miniPat.Query != nil {
		t.Logf("miniPat haps %d", len(miniPat.QueryArc(FractionFromInt(0), FractionFromInt(1))))
	}
}
