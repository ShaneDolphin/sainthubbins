package core

import "testing"

func TestFreeWrappers(t *testing.T) {
	if len(StackFree(Pure("a"), Pure("b")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("StackFree")
	}
	if len(CatFree(Pure("a")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("CatFree")
	}
	if len(FastFree(2, Pure("a")).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("FastFree")
	}
	if len(ChopFree(2, Pure(map[string]any{"s": "bd"})).QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("ChopFree")
	}
	if len(SliceFree(4, 0, "bd").QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("SliceFree")
	}
}
