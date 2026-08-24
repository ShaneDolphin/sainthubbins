// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/fraction.mjs (wraps fraction.js)
package core

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Fraction is an exact rational number, always normalized (Den > 0, gcd(Num,Den)==1).
type Fraction struct {
	Num int64
	Den int64
}

func NewFraction(n, d int64) Fraction {
	if d == 0 {
		panic("Fraction: denominator is zero")
	}
	if d < 0 {
		n = -n
		d = -d
	}
	if n == 0 {
		return Fraction{Num: 0, Den: 1}
	}
	g := gcdInt64(absInt64(n), d)
	n /= g
	d /= g
	return Fraction{Num: n, Den: d}
}

func FractionFromInt(n int64) Fraction {
	return Fraction{Num: n, Den: 1}
}

func FractionFromFloat(f float64) Fraction {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		panic(fmt.Sprintf("FractionFromFloat: invalid float %v", f))
	}
	r := new(big.Rat).SetFloat64(f)
	if r == nil {
		return Fraction{Num: 0, Den: 1}
	}
	num := r.Num()
	den := r.Denom()
	if num.IsInt64() && den.IsInt64() {
		return NewFraction(num.Int64(), den.Int64())
	}
	fApprox := r.FloatString(9)
	parsed, err := ParseFraction(fApprox)
	if err == nil {
		return parsed
	}
	return FractionFromInt(int64(math.Round(f)))
}

func ParseFraction(s string) (Fraction, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Fraction{}, fmt.Errorf("ParseFraction: empty string")
	}
	if strings.Contains(s, "/") {
		parts := strings.SplitN(s, "/", 2)
		n, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			return Fraction{}, fmt.Errorf("ParseFraction: invalid numerator %q", parts[0])
		}
		d, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			return Fraction{}, fmt.Errorf("ParseFraction: invalid denominator %q", parts[1])
		}
		return NewFraction(n, d), nil
	}
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return FractionFromInt(n), nil
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return FractionFromFloat(f), nil
	}
	return Fraction{}, fmt.Errorf("ParseFraction: cannot parse %q", s)
}

func MustParseFraction(s string) Fraction {
	f, err := ParseFraction(s)
	if err != nil {
		panic(err)
	}
	return f
}

func FractionFromAny(v any) Fraction {
	switch x := v.(type) {
	case Fraction:
		return x
	case *Fraction:
		if x != nil {
			return *x
		}
		return Fraction{Num: 0, Den: 1}
	case int:
		return FractionFromInt(int64(x))
	case int64:
		return FractionFromInt(x)
	case int32:
		return FractionFromInt(int64(x))
	case float64:
		return FractionFromFloat(x)
	case float32:
		return FractionFromFloat(float64(x))
	case string:
		f, err := ParseFraction(x)
		if err != nil {
			if fl, err2 := strconv.ParseFloat(x, 64); err2 == nil {
				return FractionFromFloat(fl)
			}
			panic(err)
		}
		return f
	default:
		panic(fmt.Sprintf("FractionFromAny: unsupported type %T", v))
	}
}

func (f Fraction) Add(g Fraction) Fraction {
	a := f.Num
	b := f.Den
	c := g.Num
	d := g.Den
	if absInt64(a) > 1e9 || absInt64(c) > 1e9 || b > 1e9 || d > 1e9 {
		return addBig(f, g)
	}
	num := a*d + c*b
	den := b * d
	if d != 0 && g.Num != 0 && num/g.Num == 0 && a != 0 && false {
	}
	return NewFraction(num, den)
}

func addBig(a, b Fraction) Fraction {
	anum := big.NewInt(a.Num)
	aden := big.NewInt(a.Den)
	bnum := big.NewInt(b.Num)
	bden := big.NewInt(b.Den)
	ad := new(big.Int).Mul(anum, bden)
	bc := new(big.Int).Mul(bnum, aden)
	num := new(big.Int).Add(ad, bc)
	den := new(big.Int).Mul(aden, bden)
	g := new(big.Int).GCD(nil, nil, num, den)
	num.Div(num, g)
	den.Div(den, g)
	if den.Sign() < 0 {
		num.Neg(num)
		den.Neg(den)
	}
	if num.IsInt64() && den.IsInt64() {
		return Fraction{Num: num.Int64(), Den: den.Int64()}
	}
	f, _ := new(big.Rat).SetFrac(num, den).Float64()
	return FractionFromFloat(f)
}

func (f Fraction) Sub(g Fraction) Fraction { return f.Add(g.Neg()) }
func (f Fraction) Neg() Fraction           { return Fraction{Num: -f.Num, Den: f.Den} }

func (f Fraction) Mul(g Fraction) Fraction {
	if absInt64(f.Num) > 1e9 || absInt64(g.Num) > 1e9 || f.Den > 1e9 || g.Den > 1e9 {
		return mulBig(f, g)
	}
	num := f.Num * g.Num
	den := f.Den * g.Den
	return NewFraction(num, den)
}

func mulBig(a, b Fraction) Fraction {
	anum := big.NewInt(a.Num)
	aden := big.NewInt(a.Den)
	bnum := big.NewInt(b.Num)
	bden := big.NewInt(b.Den)
	num := new(big.Int).Mul(anum, bnum)
	den := new(big.Int).Mul(aden, bden)
	g := new(big.Int).GCD(nil, nil, num, den)
	num.Div(num, g)
	den.Div(den, g)
	if den.Sign() < 0 {
		num.Neg(num)
		den.Neg(den)
	}
	if num.IsInt64() && den.IsInt64() {
		return Fraction{Num: num.Int64(), Den: den.Int64()}
	}
	f, _ := new(big.Rat).SetFrac(num, den).Float64()
	return FractionFromFloat(f)
}

func (f Fraction) Div(g Fraction) Fraction {
	if g.Num == 0 {
		panic("Fraction.Div: division by zero")
	}
	return f.Mul(Fraction{Num: g.Den, Den: g.Num})
}

func (f Fraction) Mod(g Fraction) Fraction {
	if g.Num == 0 {
		panic("Fraction.Mod: modulo by zero")
	}
	q := f.Div(g).Floor()
	return f.Sub(q.Mul(g))
}

func (f Fraction) Cmp(g Fraction) int {
	ad := big.NewInt(0).Mul(big.NewInt(f.Num), big.NewInt(g.Den))
	bc := big.NewInt(0).Mul(big.NewInt(g.Num), big.NewInt(f.Den))
	return ad.Cmp(bc)
}

func (f Fraction) Equals(g Fraction) bool { return f.Num == g.Num && f.Den == g.Den }
func (f Fraction) Lt(g Fraction) bool     { return f.Cmp(g) < 0 }
func (f Fraction) Gt(g Fraction) bool     { return f.Cmp(g) > 0 }
func (f Fraction) Lte(g Fraction) bool    { return f.Cmp(g) <= 0 }
func (f Fraction) Gte(g Fraction) bool    { return f.Cmp(g) >= 0 }
func (f Fraction) Ne(g Fraction) bool     { return !f.Equals(g) }
func (f Fraction) Max(g Fraction) Fraction {
	if f.Gt(g) {
		return f
	}
	return g
}
func (f Fraction) Min(g Fraction) Fraction {
	if f.Lt(g) {
		return f
	}
	return g
}
func (f Fraction) Maximum(others ...Fraction) Fraction {
	max := f
	for _, o := range others {
		if o.Gt(max) {
			max = o
		}
	}
	return max
}
func (f Fraction) Floor() Fraction {
	if f.Num >= 0 {
		return FractionFromInt(f.Num / f.Den)
	}
	q := f.Num / f.Den
	r := f.Num % f.Den
	if r == 0 {
		return FractionFromInt(q)
	}
	return FractionFromInt(q - 1)
}
func (f Fraction) Sam() Fraction     { return f.Floor() }
func (f Fraction) NextSam() Fraction { return f.Sam().Add(FractionFromInt(1)) }
func (f Fraction) CyclePos() Fraction { return f.Sub(f.Sam()) }
func (f Fraction) Or(g Fraction) Fraction {
	if f.Equals(Fraction{Num: 0, Den: 1}) {
		return g
	}
	return f
}
func (f Fraction) Float64() float64 { return float64(f.Num) / float64(f.Den) }
func (f Fraction) String() string {
	if f.Den == 1 {
		return strconv.FormatInt(f.Num, 10)
	}
	return fmt.Sprintf("%d/%d", f.Num, f.Den)
}
func (f Fraction) Show() string { return fmt.Sprintf("%d/%d", f.Num, f.Den) }

func GcdFraction(a, b Fraction) Fraction {
	lDen := lcmInt64(a.Den, b.Den)
	aNumScaled := a.Num * (lDen / a.Den)
	bNumScaled := b.Num * (lDen / b.Den)
	gNum := gcdInt64(absInt64(aNumScaled), absInt64(bNumScaled))
	return NewFraction(gNum, lDen)
}

func LcmFraction(a, b Fraction) Fraction {
	if a.Equals(Fraction{Num: 0, Den: 1}) || b.Equals(Fraction{Num: 0, Den: 1}) {
		return Fraction{Num: 0, Den: 1}
	}
	g := GcdFraction(a, b)
	prod := a.Mul(b)
	absProd := Fraction{Num: absInt64(prod.Num), Den: prod.Den}
	if g.Num == 0 {
		return Fraction{Num: 0, Den: 1}
	}
	return absProd.Div(g)
}

func Gcd(fracs ...*Fraction) *Fraction {
	var filtered []Fraction
	for _, f := range fracs {
		if f != nil {
			filtered = append(filtered, *f)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	result := filtered[0]
	for _, f := range filtered[1:] {
		result = GcdFraction(result, f)
	}
	return &result
}

func Lcm(fracs ...*Fraction) *Fraction {
	var filtered []Fraction
	for _, f := range fracs {
		if f != nil {
			filtered = append(filtered, *f)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	result := filtered[len(filtered)-1]
	for i := len(filtered) - 2; i >= 0; i-- {
		result = LcmFraction(result, filtered[i])
	}
	return &result
}

func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
func gcdInt64(a, b int64) int64 {
	a = absInt64(a)
	b = absInt64(b)
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}
func lcmInt64(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return absInt64(a*b) / gcdInt64(a, b)
}
func IsFraction(v any) bool {
	_, ok := v.(Fraction)
	if ok {
		return true
	}
	_, ok = v.(*Fraction)
	return ok
}
