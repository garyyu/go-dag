
// Copyright 2018 The go-phantom Authors
// This file is part of the go-phantom library.
//
// The go-phantom library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-phantom library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-phantom library. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	."github.com/garyyu/go-phantom/phantom"
	."github.com/garyyu/go-phantom/utils"
	"testing"
)


func chainFig3Initialize() map[string]*Block{

	//initial an empty chain
	chain := make(map[string]*Block)

	//add blocks

	ChainAddBlock("Genesis", []string{}, chain)

	ChainAddBlock("B", []string{"Genesis"}, chain)
	ChainAddBlock("C", []string{"Genesis"}, chain)
	ChainAddBlock("D", []string{"Genesis"}, chain)
	ChainAddBlock("E", []string{"Genesis"}, chain)

	ChainAddBlock("F", []string{"B","C"}, chain)
	ChainAddBlock("H", []string{"C","D","E"}, chain)
	ChainAddBlock("I", []string{"E"}, chain)

	ChainAddBlock("J", []string{"F","H"}, chain)
	ChainAddBlock("K", []string{"B","H","I"}, chain)
	ChainAddBlock("L", []string{"D","I"}, chain)
	ChainAddBlock("M", []string{"F","K"}, chain)

	tips := FindTips(chain)
	tipsName := LTPQ(tips, true)	// LTPQ is not relevant here, I just use it to get Tips name.
	ChainAddBlock("Virtual", tipsName, chain)

	fmt.Println("chainInitialize(): done. blocks=", len(chain)-1)

	return chain
}


// Tests Algorithm 1 Selection of a blue set, with the example on paper page 7 Fig.3
//
func TestFig3(t *testing.T) {

	fmt.Println("\n- Phantom Paper Simulation - Algorithm 1: Selection of a blue set. -")
	fmt.Println("-                   The example on page 7 Fig.3.                   -\n")

	chain := chainFig3Initialize()

	CalcBlue(chain, 3, chain["Virtual"])

	// print the result of blue sets

	ltpq := LTPQ(chain, true)

	fmt.Print("blue set selection done. blue blocks = ")
	nBlueBlocks := 0
	for _, name := range ltpq {
		block := chain[name]
		if IsBlueBlock(block)==true {
			if name=="Genesis" || name=="Virtual" {
				fmt.Print("(",name[:1],").")
			}else {
				fmt.Print(block.Name, ".")
			}

			nBlueBlocks++
		}
	}
	fmt.Println("	total blue:", nBlueBlocks)
	t.Error("I'm in a bad mood.")
}

