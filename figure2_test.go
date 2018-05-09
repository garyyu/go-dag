
// Copyright 2018 The godag Authors
// This file is part of the godag library.
//
// The godag library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The godag library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the godag library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	."github.com/garyyu/go-dag/godag"
	."github.com/garyyu/go-dag/utils"
	"testing"
	"bytes"
)


func chainFig2Initialize() map[string]*Block{

	//initial an empty chain
	chain := make(map[string]*Block)

	//add blocks

	ChainAddBlock("Genesis", []string{}, chain)

	ChainAddBlock("B", []string{"Genesis"}, chain)
	ChainAddBlock("C", []string{"Genesis"}, chain)
	ChainAddBlock("D", []string{"Genesis"}, chain)
	ChainAddBlock("E", []string{"Genesis"}, chain)

	ChainAddBlock("F", []string{"B","C"}, chain)
	ChainAddBlock("G", []string{"C","D"}, chain)
	ChainAddBlock("H", []string{"E"}, chain)

	ChainAddBlock("I", []string{"F","D"}, chain)
	ChainAddBlock("J", []string{"B","G","E"}, chain)
	ChainAddBlock("K", []string{"D","H"}, chain)

	tips := FindTips(chain)
	tipsName := LTPQ(tips, true)	// LTPQ is not relevant here, I just use it to get Tips name.
	ChainAddBlock("Virtual", tipsName, chain)

	return chain
}



// Tests Algorithm 2 Ordering of the DAG, with the example on paper page 3 Fig.2
//
func TestFig2(t *testing.T) {

	var actual bytes.Buffer
	var expected string

	fmt.Println("\n-  BlockDAG Algorithms Simulation - Algorithm 2: Ordering of the DAG.   -")
	fmt.Println("-        The example on Fig.2         -")

	chain := chainFig2Initialize()

	fmt.Println("chainInitialize(): done. blocks=", len(chain)-1)

	orderedList := Order(chain, 3)

	// print the result of blue sets

	ltpq := LTPQ(chain, true)

	fmt.Print("blue set selection done. blue blocks = ")
	nBlueBlocks := 0
	actual.Reset()
	for _, name := range ltpq {
		block := chain[name]
		if IsBlueBlock(block)==true {
			if name=="Genesis" || name=="Virtual" {
				actual.WriteString(fmt.Sprintf("(%s).",name[:1]))
			}else {
				actual.WriteString(name+".")
			}

			nBlueBlocks++
		}
	}
	fmt.Println(actual.String(), "	total blue:", nBlueBlocks)

	expected = "(G).B.C.D.F.G.I.J.(V)."
	if actual.String() != expected {
		t.Errorf("blue selection test not matched. expected=%s", expected)
	}

	// print the result of ordered blocks of this chain

	fmt.Print("ordered chain blocks = ")
	actual.Reset()
	for _, name := range orderedList {
		if name=="Genesis" || name=="Virtual" {
			actual.WriteString(fmt.Sprintf("(%s).",name[:1]))
		} else {
			actual.WriteString(name+".")
		}
	}
	fmt.Println(actual.String())

	expected = "(G).B.C.D.F.E.G.J.I.H.K.(V)."
	if actual.String() != expected {
		t.Errorf("block ordering test not matched. expected=%s", expected)
	}
}


func BenchmarkBlockOrdering(b *testing.B) {

	for i := 0; i < b.N; i++ {
		chain := chainFig2Initialize()
		Order(chain, 3)
	}
}