// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import "testing"

func TestScaleBridge(t *testing.T) {
	p := Pure("c3").Scale("C:major")
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Scale expected haps")
	}
	if ctx, ok := haps[0].Context["scale"]; !ok || ctx != "C:major" {
		t.Fatalf("scale context %v", haps[0].Context)
	}
}

func TestTransposeBridge(t *testing.T) {
	p := Pure("c4").Transpose(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) == 0 {
		t.Fatalf("Transpose expected haps")
	}
	if haps[0].Value != "D4" {
		t.Fatalf("c4+2 expected D4 got %v", haps[0].Value)
	}
	p2 := Pure("c4").Transpose("3M")
	haps2 := p2.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps2[0].Value != "E4" {
		t.Fatalf("c4+3M expected E4 got %v", haps2[0].Value)
	}
	p3 := Pure(map[string]any{"note": "c4"}).Transpose(2)
	haps3 := p3.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if m, ok := haps3[0].Value.(map[string]any); !ok || m["note"] != "D4" {
		t.Fatalf("map transpose expected D4 got %v", haps3[0].Value)
	}
}

func TestTransAlias(t *testing.T) {
	p := Pure("c4").Trans(1)
	if len(p.QueryArc(FractionFromInt(0), FractionFromInt(1))) == 0 {
		t.Fatalf("Trans alias")
	}
}
