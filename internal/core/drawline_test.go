package core

import "testing"

func TestDrawLine(t *testing.T) {
	pat := FastCat(Pure("a"), Pure("b"), Pure("c"))
	line := DrawLine(pat, 20)
	if len(line) == 0 {
		t.Fatalf("drawline empty")
	}
	// Should contain cycle separators
	hasPipe := false
	for _, c := range line {
		if c == '|' {
			hasPipe = true
			break
		}
	}
	if !hasPipe {
		t.Fatalf("drawline no pipe: %q", line)
	}
	t.Logf("drawline: %q", line)
}

func TestPerlin(t *testing.T) {
	p := Perlin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("perlin no haps")
	}
	v := toFloat(haps[0].Value)
	if v < -1.5 || v > 1.5 {
		t.Fatalf("perlin out of range %v", v)
	}
	p2 := PerlinWith(1.0)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("perlinWith no haps")
	}
}
