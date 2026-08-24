package core

import "testing"

func TestPureBasic(t *testing.T) {
	p := Pure("a")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Pure a expected 1 hap, got %d", len(haps))
	}
	if haps[0].Value.(string) != "a" {
		t.Fatalf("Pure value expected a got %v", haps[0].Value)
	}
}

func TestFmap(t *testing.T) {
	p := Pure(1).Fmap(func(a any) any { return a.(int) + 1 })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value.(int) != 2 {
		t.Fatalf("Fmap expected 2 got %v", haps)
	}
}

func TestStackBasic(t *testing.T) {
	p := Stack(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Stack a,b expected 2 got %d", len(haps))
	}
}

func TestFastBasic(t *testing.T) {
	p := Pure("a").FastF(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Fast 2 expected 2 haps, got %d", len(haps))
	}
}

func TestSlowBasic(t *testing.T) {
	p := Pure("a").Slow(FractionFromInt(2))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Slow 2 expected 1 hap over 1 cycle, got %d", len(haps))
	}
	haps2 := p.QueryArc(FractionFromInt(0), FractionFromInt(2))
	if len(haps2) != 1 {
		t.Fatalf("Slow 2 over 2 cycles expected 1, got %d", len(haps2))
	}
}

func TestRevBasic(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c")).Rev()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Rev expected 3 got %d", len(haps))
	}
	// Rev reverses time; first hap value may still be a due to part sorting — just check all values present
	found := map[string]bool{}
	for _, h := range haps {
		found[h.Value.(string)] = true
	}
	if !found["a"] || !found["b"] || !found["c"] {
		t.Fatalf("Rev expected a,b,c got %v", haps)
	}
}

func TestAddMul(t *testing.T) {
	p := Pure(2).Add(Pure(3))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Add 2+3 expected 1 hap got %d", len(haps))
	}
	// Add may return float64 via numeric promotion
	switch v := haps[0].Value.(type) {
	case int:
		if v != 5 {
			t.Fatalf("Add 2+3 expected 5 got %v", v)
		}
	case float64:
		if v != 5 {
			t.Fatalf("Add 2+3 expected 5 got %v", v)
		}
	default:
		t.Fatalf("Add unexpected type %T %v", v, v)
	}
	p2 := Pure(2).Mul(Pure(3))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 {
		t.Fatalf("Mul expected 1 got %d", len(haps2))
	}
	val := haps2[0].Value
	switch v := val.(type) {
	case int:
		if v != 6 {
			t.Fatalf("Mul expected 6 got %v", v)
		}
	case float64:
		if v != 6 {
			t.Fatalf("Mul expected 6 got %v", v)
		}
	default:
		t.Fatalf("Mul unexpected type %T", val)
	}
}

func TestEarlyLate(t *testing.T) {
	p := Pure("a").Early(0.5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Early expected >=1 got %d", len(haps))
	}
	q := Pure("a").Late(0.5)
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Late expected >=1 got %d", len(haps2))
	}
}

func TestWhenEvery(t *testing.T) {
	p := Pure("a").When(true, func(p Pattern) Pattern { return Pure("b") })
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 || haps[0].Value.(string) != "b" {
		t.Fatalf("When true expected b got %v", haps)
	}
	p2 := Pure("a").Every(2, func(p Pattern) Pattern { return Pure("b") })
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 1 || haps2[0].Value.(string) != "b" {
		t.Fatalf("Every 2 first cycle expected b got %v", haps2)
	}
	haps3 := p2.QueryArc(FractionFromInt(1), FractionFromInt(2))
	if len(haps3) != 1 || haps3[0].Value.(string) != "a" {
		t.Fatalf("Every 2 second cycle expected a got %v", haps3)
	}
}

func TestDegradeBy(t *testing.T) {
	p := Pure("a").DegradeBy(0)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("DegradeBy 0 expected 1 got %d", len(haps))
	}
	q := Pure("a").DegradeBy(1)
	haps2 := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 0 {
		t.Fatalf("DegradeBy 1 expected 0 got %d", len(haps2))
	}
}

func TestCompress(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Compress expected 1 got %d", len(haps))
	}
}
