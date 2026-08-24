package transpiler

import "testing"

func TestTranspileBasic(t *testing.T) {
	res, err := Transpile(`s("bd sd")`, Options{})
	if err != nil {
		t.Fatalf("transpile err %v", err)
	}
	if res.Output == "" {
		t.Fatalf("empty output")
	}
	t.Logf("output: %s", res.Output)
}

func TestTranspileMiniLocations(t *testing.T) {
	res, _ := Transpile(`s("bd ~ sd")`, Options{})
	if len(res.MiniLocations) == 0 {
		t.Logf("no mini locations, output %s", res.Output)
	}
}
