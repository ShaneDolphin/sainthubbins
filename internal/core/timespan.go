// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
// Saint Hubbins — live pattern engine.
// Original: packages/core/timespan.mjs
package core

import "fmt"

// TimeSpan represents a span of time [Begin, End) — half-open interval.
// Mirrors JS TimeSpan class.
type TimeSpan struct {
	Begin Fraction
	End   Fraction
}

func NewTimeSpan(begin, end Fraction) TimeSpan {
	return TimeSpan{Begin: begin, End: end}
}

func (s TimeSpan) Duration() Fraction {
	return s.End.Sub(s.Begin)
}

func (s TimeSpan) Equals(other TimeSpan) bool {
	return s.Begin.Equals(other.Begin) && s.End.Equals(other.End)
}

func (s TimeSpan) String() string {
	return s.Begin.Show() + " → " + s.End.Show()
}
func (s TimeSpan) Show() string { return s.String() }

// SpanCycles splits the span at cycle boundaries, like JS spanCycles getter.
func (s TimeSpan) SpanCycles() []TimeSpan {
	begin := s.Begin
	end := s.End
	endSam := end.Sam()

	if begin.Equals(end) {
		return []TimeSpan{NewTimeSpan(begin, end)}
	}

	var spans []TimeSpan
	for end.Gt(begin) {
		if begin.Sam().Equals(endSam) {
			spans = append(spans, NewTimeSpan(begin, s.End))
			break
		}
		nextBegin := begin.NextSam()
		spans = append(spans, NewTimeSpan(begin, nextBegin))
		begin = nextBegin
	}
	return spans
}

// CycleArc shifts timespan to one of equal duration that starts within cycle 0.
func (s TimeSpan) CycleArc() TimeSpan {
	b := s.Begin.CyclePos()
	e := b.Add(s.Duration())
	return NewTimeSpan(b, e)
}

// WithTime applies fn to both begin and end.
func (s TimeSpan) WithTime(fn func(Fraction) Fraction) TimeSpan {
	return NewTimeSpan(fn(s.Begin), fn(s.End))
}

// WithEnd applies fn only to end.
func (s TimeSpan) WithEnd(fn func(Fraction) Fraction) TimeSpan {
	return NewTimeSpan(s.Begin, fn(s.End))
}

// WithCycle applies fn relative to the cycle of the span's begin.
func (s TimeSpan) WithCycle(fn func(Fraction) Fraction) TimeSpan {
	sam := s.Begin.Sam()
	b := sam.Add(fn(s.Begin.Sub(sam)))
	e := sam.Add(fn(s.End.Sub(sam)))
	return NewTimeSpan(b, e)
}

// Intersection returns the intersection of two TimeSpans, or nil if they don't intersect.
// Mirrors JS intersection() exactly, including zero-width edge cases.
func (s TimeSpan) Intersection(other TimeSpan) *TimeSpan {
	intersectBegin := s.Begin.Max(other.Begin)
	intersectEnd := s.End.Min(other.End)

	if intersectBegin.Gt(intersectEnd) {
		return nil
	}
	if intersectBegin.Equals(intersectEnd) {
		// Zero-width intersection doesn't count if it's at the end of a non-zero-width span
		if intersectBegin.Equals(s.End) && s.Begin.Lt(s.End) {
			return nil
		}
		if intersectBegin.Equals(other.End) && other.Begin.Lt(other.End) {
			return nil
		}
	}
	result := NewTimeSpan(intersectBegin, intersectEnd)
	return &result
}

// IntersectionE like intersection but panics if no intersection.
func (s TimeSpan) IntersectionE(other TimeSpan) TimeSpan {
	result := s.Intersection(other)
	if result == nil {
		panic(fmt.Sprintf("TimeSpans do not intersect: %s and %s", s.String(), other.String()))
	}
	return *result
}

func (s TimeSpan) Midpoint() Fraction {
	return s.Begin.Add(s.Duration().Div(FractionFromInt(2)))
}
