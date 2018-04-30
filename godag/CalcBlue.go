
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

package godag

import (
	"fmt"
	"os"
	"sort"
)

type Block struct {
	Name string					// Name used for simulation. In reality, block header hash can be used as the 'Name' of a block.
	Score int
	SizeOfPastSet int			// Size of its past set (blocks). It's fixed number once a block is created, and can't change anymore.
	Prev map[string]*Block
	Next map[string]*Block		// Block don't have this info, it comes from the analysis of existing chain
	Blue map[string]bool		// Blue is relative to each Tip
}

/*
 * if a block's Next is not in G, then it's a Tip
 */
func FindTips(G map[string]*Block) map[string]*Block {

	tips := make(map[string]*Block)

	checkNextBlock:
		for k, v := range G {

			if len(v.Next)==0 {
				tips[k] = v
			}else{
				for next := range v.Next {
					if _,ok := G[next]; ok {
						continue checkNextBlock
					}
				}

				tips[v.Name] = v
			}
		}

	return tips
}

func pastSet(B *Block, past map[string]*Block){

	for k, v := range B.Prev {

		if _,ok := past[k]; !ok {
			pastSet(v, past)
		}
		past[k] = v
	}
}

func futureSet(B *Block, future map[string]*Block){

	for k, v := range B.Next {

		if _,ok := future[k]; !ok {
			futureSet(v, future)
		}
		future[k] = v
	}
}


func countBlue(G map[string]*Block, tip *Block) int{

	var blueBlocks = 0
	for _, v := range G {
		if blue, ok := v.Blue[tip.Name]; ok {
			if blue {
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

func IsBlueBlock(B *Block) bool {

	if B==nil {
		return false
	}

	if B.Name=="Genesis" {
		return true
	}

	for _,blue := range B.Blue {
		if blue==true {
			return true
		}
	}

	return false
}

func SizeOfPastSet(B *Block) int{

	past := make(map[string]*Block)
	pastSet(B, past)
	return len(past)
}

/*
 * lexicographical topological priority queue
 */
func LTPQ(G map[string]*Block, ascending bool) []string{

	ltpq := make([]struct {
		Name string
		SizeOfPastSet  int
	}, len(G))

	i:=0
	for _, block := range G {
		ltpq[i].Name, ltpq[i].SizeOfPastSet = block.Name, block.SizeOfPastSet
		i++
	}

	/*
	 * The priority of a block is represented by the size of its past set; in case of ties, the block with lowest hash ID is chosen.
	 */
	sort.Slice(ltpq, func(i, j int) bool {
		return ltpq[i].SizeOfPastSet < ltpq[j].SizeOfPastSet || (ltpq[i].SizeOfPastSet == ltpq[j].SizeOfPastSet && ltpq[i].Name < ltpq[j].Name)
		})

	priorityQue := make([]string, len(ltpq))

	if ascending==true {	// asc: little first
		for i := 0; i < len(ltpq); i++ {
			priorityQue[i] = ltpq[i].Name
		}
	}else{					// desc: big first
		for i,j := len(ltpq)-1,0; i >= 0; i,j = i-1,j+1 {
			priorityQue[j] = ltpq[i].Name
		}
	}

	return priorityQue
}

func CalcBlue(G map[string]*Block, k int, topTip *Block){

	if len(G)==1 {
		if _,ok := G["Genesis"]; ok {
			return
		}else{
			fmt.Println("CalcBlue(): error! len(G)=1 but not Genesis block")
			os.Exit(-1)
		}
	} else if len(G)==0 {
		fmt.Println("CalcBlue(): error! impossible to reach here. len(G)=0")
		os.Exit(-1)
	}

	// step 4
	tips := FindTips(G)
	if len(tips)==0 {
		fmt.Println("calcBlue(): error! impossible! Tips Empty.")
		os.Exit(-1)
	}

	maxBlue := -1
	var Bmax *Block = nil

	// step 4'	 	Starting from the highest scoring tip (to be continued...  for the moment, I use reversed LTPQ.)
	var ltpq = LTPQ(tips, false)

	for _, name := range ltpq {
		tip := tips[name]		// step 4'

		past := make(map[string]*Block)
		pastSet(tip, past)

		//fmt.Println("calcBlue(): info of next recursive call - tip=", tip.Name, " past=", len(past))

		// step 5
		CalcBlue(past, k, tip)

		// step 6
		blueBlocks := countBlue(past, tip)
		if blueBlocks>maxBlue {
			maxBlue = blueBlocks
			Bmax = tip
		} else if blueBlocks==maxBlue {
			// tie-breaking
			/*
			 * Important Note: in some cases, the tie-breaking method will decide the final blue selection result! not always converge to maximum k-cluster SubDAG.
			 *				   research to be continued.
			 */
			if tip.Name < Bmax.Name {
				Bmax = tip
			}
		}
	}

	if Bmax==nil {
		fmt.Println("calcBlue(): error! impossible! Bmax=nil.")
		os.Exit(-1)
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

	// step 8
	anticoneBmax := antiCone(G, Bmax)
	ltpq = LTPQ(anticoneBmax, true)			// in 'some' topological ordering: LTPQ
	for _, name := range ltpq {

		B := anticoneBmax[name]

		// step 9
		nBlueAnticone := 0
		anticone := antiCone(G, B)
		for _, block := range anticone {
			if IsBlueBlock(block)==true {
				nBlueAnticone++
			}
		}

		if nBlueAnticone<=k {
			B.Blue[Bmax.Name] = true	// step 10
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
