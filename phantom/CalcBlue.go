
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

package phantom

import (
	"strings"
	"fmt"
	"os"
)

type Block struct {
	Name string					// Name used for simulation. In reality, block header hash can be used as the 'Name' of a block.
	Score int
	Prev map[string]*Block
	Next map[string]*Block		// Block don't have this info, it comes from the analysis of existing chain
	Blue map[string]bool		// Blue is relative to each Tip
}


func findTips(G map[string]*Block) map[string]*Block {

	tips := make(map[string]*Block)
	for k, v := range G {
		tips[k] = v
	}

	for _, v := range G {
		for _, prev := range v.Prev {
			// if block is referenced by anyone, then it must not a tip
			_, ok := tips[prev.Name]
			if ok {
				delete(tips, prev.Name)
			}
		}
	}

	return tips
}

func pastSet(B *Block, past map[string]*Block){

	for k, v := range B.Prev {
		pastSet(v, past)
		past[k] = v
	}
}

func futureSet(B *Block, future map[string]*Block){

	for k, v := range B.Next {
		future[k] = v
		futureSet(v, future)
	}
}


func countBlue(G map[string]*Block, tip *Block) int{

	var blueBlocks = 0
	for _, v := range G {
		if blue, ok := v.Blue[tip.Name]; ok {
			if blue==true{
				blueBlocks++
			}
		} else if v.Name=="Genesis"{
			blueBlocks++
		}

	}

	return blueBlocks
}

func antiCone(G map[string]*Block, B *Block) map[string]*Block{

	anticone := make(map[string]*Block)

	past := make(map[string]*Block)
	pastSet(B, past)

	future := make(map[string]*Block)
	futureSet(B, future)

	for name, block := range G {
		if _,ok := past[name]; ok {
			continue				// block not in B's past
		}

		if _,ok := future[name]; ok {
			continue				// block not in B's future
		}

		if name==B.Name {
			continue				// block not B
		}

		anticone[name] = block		// then this block belongs to anticone
	}

	return anticone
}

func CalcBlue(G map[string]*Block, k int, topTip *Block){

	defer func() {
		if x := recover(); x != nil {
			// recovering from a panic; x contains whatever was passed to panic()
			fmt.Println("CalcBlue(): tip=", topTip.Name, ". run time panic =", x)

			//panic(x)
			os.Exit(-1)
		}
	}()

	if len(G)==1 {
		if _,ok := G["Genesis"]; ok {
			fmt.Println("CalcBlue(): return from Genesis")
			return
		}else{
			fmt.Println("CalcBlue(): error! len(G)=1 but not Genesis block")
			os.Exit(-1)
		}
	} else if len(G)==0 {
		fmt.Println("CalcBlue(): error! impossible to reach here. len(G)=0")
		os.Exit(-1)
	}

	// step 4,5
	tips := findTips(G)
	maxBlue := -1
	var Bmax *Block = nil
	for _, tip := range tips {
		past := make(map[string]*Block)
		pastSet(tip, past)

		fmt.Println("calcBlue(): info of next recursive call - tip=", tip.Name, " past=", len(past))

		CalcBlue(past, k, tip)

		// step 6
		blueBlocks := countBlue(past, tip)
		if blueBlocks>maxBlue {
			maxBlue = blueBlocks
			Bmax = tip
		} else if blueBlocks==maxBlue {
			// tie-breaking
			if strings.Compare(tip.Name,Bmax.Name)>0 {
				Bmax = tip
			}
		}
	}

	// step 7
	for _, v := range G {
		for name, blue := range v.Blue {
			for _, tip := range tips {
				if blue == true && name != Bmax.Name && name==tip.Name {
					v.Blue[name] = false // clear all other tips blue blocks, only keep the Bmax blue ones
				}
			}
		}
	}
	Bmax.Blue[Bmax.Name] = true		// BLUEk(G) = BLUEk(Bmax) U {Bmax}

	// step 8,9,10
	anticoneBmax := antiCone(G, Bmax)
	for _, B := range anticoneBmax {
		nBlueAnticone := 0
		anticone := antiCone(G, B)
		for _, block := range anticone {
			if block.Blue[Bmax.Name]==true {
				nBlueAnticone++
			}
		}

		if nBlueAnticone<=k {
			B.Blue[Bmax.Name] = true
		}
	}

	// additional step: replace Blue[Bmax] with [topTip]
	for _, B := range G {
		if blue, ok := B.Blue[Bmax.Name]; ok && blue==true {
			B.Blue[Bmax.Name] = false
			B.Blue[topTip.Name] = true
		}
	}

}
