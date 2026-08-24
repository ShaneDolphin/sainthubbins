package transpiler

import "testing"

func TestTranspileSingleQuoteMini(t *testing.T) {
	res, err := Transpile(`s('bd sd ~')`, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.MiniLocations) == 0 {
		t.Fatalf("expected mini location for single quote bd sd, got %v", res.MiniLocations)
	}
}
