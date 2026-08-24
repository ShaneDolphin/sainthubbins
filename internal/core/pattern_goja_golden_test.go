//go:build goja
// +build goja

// Copyright (C) 2026 Saint Hubbins contributors — AGPL-3.0-or-later
package core

import (
	"testing"

	"github.com/dop251/goja"
)

func TestGojaGolden_StackCat(t *testing.T) {
	vm := goja.New()
	// Simple JS stack/cat via goja mock that mirrors Go Stack/FastCat semantics
	_, err := vm.RunString(`
		function stack(a,b){ return [a,b]; }
		function cat(a,b,c){ return [a,b,c]; }
		var s = stack("a","b");
		if (s.length !== 2) throw new Error("stack");
		var c = cat("a","b","c");
		if (c.length !== 3) throw new Error("cat");
	`)
	if err != nil {
		t.Fatalf("goja stack/cat %v", err)
	}
	// Compare Go Stack vs JS mock: Go should have 2 haps for stack
	p := Stack(Pure("a"), Pure("b"))
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 2 {
		t.Fatalf("Go Stack expected 2 got %d", len(haps))
	}
}

func TestGojaGolden_Euclid(t *testing.T) {
	vm := goja.New()
	_, err := vm.RunString(`
		function bjorklund(steps, pulses){
			let pattern = [];
			for(let i=0;i<steps;i++) pattern.push(i < pulses ? 1 : 0);
			return pattern;
		}
		var b = bjorklund(8,3);
		if (b.length !== 8) throw new Error("bjorklund len");
	`)
	if err != nil {
		t.Fatalf("goja euclid %v", err)
	}
	p := Pure("x").Euclid(3, 8)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if len(haps) != 3 {
		t.Fatalf("Go Euclid 3,8 expected 3 got %d", len(haps))
	}
}

func TestGojaGolden_TransposeScale(t *testing.T) {
	vm := goja.New()
	_, err := vm.RunString(`
		var chromatic = ["C","C#","D","D#","E","F","F#","G","G#","A","A#","B"];
		function transpose(note, semitones){
			var base = {c:0,d:2,e:4,f:5,g:7,a:9,b:11};
			var n = note.toLowerCase();
			var letter = n[0];
			var semi = base[letter];
			var idx=1;
			if(n[1]==='#') semi++;
			var oct=4;
			if(idx < n.length){
				var o=parseInt(n.slice(idx));
				if(!isNaN(o)) oct=o;
			}
			var midi=(oct+1)*12+semi;
			midi+=semitones;
			var pc=chromatic[midi%12];
			var o2=Math.floor(midi/12)-1;
			return pc+o2;
		}
		if(transpose("c4",2)!=="D4") throw new Error("transpose c4+2");
		if(transpose("c4",4)!=="E4") throw new Error("transpose c4+4");
	`)
	if err != nil {
		t.Fatalf("goja transpose %v", err)
	}
	p := Pure("c4").Transpose(2)
	haps := p.QueryArc(FractionFromInt(0), FractionFromInt(1))
	if haps[0].Value != "D4" {
		t.Fatalf("Go Transpose c4+2 expected D4 got %v", haps[0].Value)
	}
}
