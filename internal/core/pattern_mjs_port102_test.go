package core

import "testing"

func TestMJS_BinaryOpAddMul2(t *testing.T) {
	a := Pure(2)
	added := a.Add(Pure(3))
	haps := added.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("Add") }
	if func()bool{v:=haps[0].Value; switch x:=v.(type){case int: return x!=5; case float64: return x!=5; default: return true}}() { t.Fatalf("5 got %v", haps[0].Value) }
	mul := a.Mul(Pure(4))
	haps = mul.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps)==0 { t.Fatalf("Mul") }
}

func TestMJS_StructureWithSignal3(t *testing.T) {
	p := Pure("bd sd")
	s := p.Struct(Pure(true).FastF(FractionFromInt(2)))
	if s.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Struct") }
}

func TestMJS_EveryWhenOff2(t *testing.T) {
	p := Pure("a").Every(2, func(q Pattern) Pattern { return q.FastF(FractionFromInt(2)) })
	if p.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("Every") }
	w := p.When(func(b bool) bool { return true }, func(q Pattern) Pattern { return q.Rev() })
	if w.QueryArc(FractionFromInt(0), FractionFromInt(1))==nil { t.Fatalf("When") }
}
