// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later

package jsapi

import "testing"

func TestStackLayers(t *testing.T) {
	if got := countHaps(t, `stack(s("bd*4"), s("hh*8"))`); got != 12 {
		t.Errorf("got %d haps, want 12", got)
	}
}

func TestCatAlternates(t *testing.T) {
	if got := countHaps(t, `cat(s("bd"), s("sd"))`); got != 1 {
		t.Errorf("got %d haps in one cycle, want 1", got)
	}
}

func TestSilenceIsEmpty(t *testing.T) {
	if got := countHaps(t, `silence()`); got != 0 {
		t.Errorf("got %d haps, want 0", got)
	}
}

func TestMiniHelper(t *testing.T) {
	if got := countHaps(t, `mini("bd sd hh")`); got != 3 {
		t.Errorf("got %d haps, want 3", got)
	}
}
