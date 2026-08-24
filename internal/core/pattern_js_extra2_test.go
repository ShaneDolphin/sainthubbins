// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestJSExtra_SetKeep(t *testing.T) {
	// JS: note("c a f e").s("sine").set(s("triangle")) -> s overrides
	p := Pure(map[string]any{"note": "c", "s": "sine"}).Set(Pure(map[string]any{"s": "triangle"}))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Set expected haps")
	}
	if m, ok := haps[0].Value.(map[string]any); ok {
		if m["s"] != "triangle" {
			t.Fatalf("Set s expected triangle got %v", m["s"])
		}
		if m["note"] != "c" {
			t.Fatalf("Set should keep note c got %v", m["note"])
		}
	}
	// Keep: should keep original s
	p2 := Pure(map[string]any{"s": "sine", "n": 1}).Keep(Pure(map[string]any{"s": "triangle", "gain": 0.5}))
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Keep expected haps")
	}
	if m, ok := haps2[0].Value.(map[string]any); ok {
		if m["s"] != "sine" {
			t.Fatalf("Keep s expected sine got %v", m["s"])
		}
		if m["gain"] != 0.5 {
			t.Fatalf("Keep gain expected 0.5 got %v", m["gain"])
		}
	}
}

func TestJSExtra_SqueezeJoin(t *testing.T) {
	// JS: squeezeJoin via SqueezeBind
	p := Pure(Pure("a")).SqueezeJoin()
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("SqueezeJoin expected haps")
	}
}

func TestJSExtra_PolymeterArrange(t *testing.T) {
	// JS: polymeter, arrange
	p := Polymeter(Pure("a"), Pure("b"), Pure("c"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Polymeter expected haps")
	}
	arr := Arrange(1, Pure("a"), 2, Pure("b"))
	haps2 := arr.QueryArc(FractionFromInt(0), FractionFromInt(3))
	if len(haps2) == 0 {
		t.Fatalf("Arrange expected haps")
	}
}

func TestJSExtra_RandomChoose(t *testing.T) {
	p := Pure(0).Choose([]any{"a", "b", "c"})
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Choose expected haps")
	}
	// Shuffle
	p2 := FastCat(Pure("a"), Pure("b"), Pure("c")).Shuffle(3)
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 3 {
		t.Fatalf("Shuffle expected 3 got %d", len(haps2))
	}
}

func TestJSExtra_FilterWithin(t *testing.T) {
	p := FastCat(Pure("a"), Pure("b"), Pure("c"), Pure("d"))
	within := p.Within(0, 0.5, func(pat Pattern) Pattern { return pat.FastF(FractionFromInt(2)) })
	haps := within.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Within expected haps")
	}
	filtered := p.FilterValues(func(v any) bool { return v == "a" || v == "b" })
	haps2 := filtered.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) != 2 {
		t.Fatalf("FilterValues a b expected 2 got %d", len(haps2))
	}
}

func TestJSExtra_CompressZoom(t *testing.T) {
	p := Pure("a").Compress(0, 0.5)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 1 {
		t.Fatalf("Compress 0 0.5 expected 1 got %d", len(haps))
	}
	if !haps[0].Part.Begin.Equals(FractionFromInt(0)) {
		t.Fatalf("Compress part begin 0")
	}
	z := Pure("a").Zoom(0.5, 1)
	haps2 := z.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps2) == 0 {
		t.Fatalf("Zoom expected haps")
	}
}
