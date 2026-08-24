package core

import "testing"

func BenchmarkStack(b *testing.B) {
	p := Stack(Pure("bd"), Pure("sd"), Pure("hh"))
	for i := 0; i < b.N; i++ {
		p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	}
}

func BenchmarkFastCat(b *testing.B) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	for i := 0; i < b.N; i++ {
		p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	}
}

func BenchmarkEuclid(b *testing.B) {
	p := Pure("x").Euclid(3, 8)
	for i := 0; i < b.N; i++ {
		p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	}
}
