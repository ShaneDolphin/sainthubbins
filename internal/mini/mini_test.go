package mini

import (
	"testing"

	"codeberg.org/uzu/saint-hubbins/internal/core"
)

func TestMiniBasic(t *testing.T) {
	p := Mini("bd sd")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Mini bd sd expected 2 got %d", len(haps))
	}
}

func TestMiniEuclid(t *testing.T) {
	p := Mini("bd(3,8)")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Mini bd(3,8) expected 3 got %d", len(haps))
	}
}

func TestMiniRest(t *testing.T) {
	p := Mini("bd ~ sd")
	haps := p.QueryArc(core.FractionFromInt(0), core.FractionFromInt(1))
	if len(haps) < 2 {
		t.Fatalf("Mini bd ~ sd expected >=2 got %d", len(haps))
	}
}
