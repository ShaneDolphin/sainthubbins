package core

import "testing"

func TestTimePure(t *testing.T) {
	p := Pure("a")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Time Pure expected 1")
	}
}

func TestTimeSlowFast(t *testing.T) {
	p := Pure("a").Slow(FractionFromInt(2))
	q := p.FastF(FractionFromInt(2))
	haps := q.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Slow/Fast roundtrip expected >=1")
	}
}

func TestTimeCompress(t *testing.T) {
	p := Pure("a").Compress(FractionFromFloat(0.25), FractionFromFloat(0.75))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Compress expected 1")
	}
}

func TestTimeRev(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b")).Rev()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Rev 2 expected 2 got %d", len(haps))
	}
}
