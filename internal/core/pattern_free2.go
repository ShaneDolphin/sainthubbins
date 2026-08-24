// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Free wrappers for top-level JS functions that are Pattern methods in Go.
// These allow `stack()`, `cat()` etc to be called as functions, mirroring JS `import { stack } from 'hubbins'`

package core

func StackFree(pats ...Pattern) Pattern { return Stack(pats...) }
func CatFree(pats ...Pattern) Pattern { return Cat(pats...) }
func FastCatFree(pats ...Pattern) Pattern { return FastCat(pats...) }
func SlowCatFree(pats ...Pattern) Pattern { return SlowCat(pats...) }
func ArrangeFree(sections ...any) Pattern { return Arrange(sections...) }
func PolymeterFree(pats ...Pattern) Pattern { return Polymeter(pats...) }
func SequenceFree(pats ...Pattern) Pattern { return Sequence(pats...) }
func TimeCatFree(pats ...any) Pattern { return TimeCat(pats...) }

// Additional top-level wrappers for methods that JS exposes as free functions via `register` curry
func FastFree(factor any, pat Pattern) Pattern { return pat.Fast(Reify(factor)) }
func SlowFree(factor any, pat Pattern) Pattern { return pat.SlowF(toFraction(Reify(factor).FirstCycleValue())) }
func EarlyFree(off any, pat Pattern) Pattern { return pat.Early(off) }
func LateFree(off any, pat Pattern) Pattern { return pat.Late(off) }
func CompressFree(b, e any, pat Pattern) Pattern { return pat.Compress(b, e) }
func ZoomFree(s, e any, pat Pattern) Pattern { return pat.Zoom(s, e) }
func PlyFree(factor any, pat Pattern) Pattern { return pat.Ply(factor) }
func OffFree(timePat any, fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.Off(timePat, fn) }
func WhenFree(test any, fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.When(test, fn) }
func EveryFree(n int, fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.Every(n, fn) }
func InsideFree(factor any, fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.Inside(factor, fn) }
func OutsideFree(factor any, fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.Outside(factor, fn) }
func RevFree(pat Pattern) Pattern { return pat.Rev() }
func PalindromeFree(pat Pattern) Pattern { return pat.Palindrome() }
func JuxFree(fn func(Pattern) Pattern, pat Pattern) Pattern { return pat.Jux(fn) }
func ChopFree(n int, pat Pattern) Pattern { return pat.Chop(n) }
func StriateFree(n int, pat Pattern) Pattern { return pat.Striate(n) }
func SliceFree(npat, ipat, opat any) Pattern { return Slice(npat, ipat, opat) }
func SpliceFree(npat, ipat, opat any) Pattern { return Splice(npat, ipat, opat) }
func FitFree(pat Pattern) Pattern { return pat.Fit() }
func LoopAtFree(factor any, pat Pattern) Pattern { return pat.LoopAt(factor) }
func BiteFree(npat, ipat, pat any) Pattern { return Bite(npat, ipat, pat) }
