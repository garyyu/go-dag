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

package utils

import (
	"os"
	"fmt"
	."github.com/garyyu/go-phantom/phantom"
)


func ChainAddBlock(Name string, References []string, chain map[string]*Block) *Block{

	//create this block
	thisBlock := Block{Name, -1, -1,make(map[string]*Block), make(map[string]*Block),make(map[string]bool)}

	//add references
	for _, Reference := range References {
		prev, ok := chain[Reference]
		if ok {
			thisBlock.Prev[Reference] = prev
			prev.Next[Name] = &thisBlock
		}else{
			fmt.Println("chainAddBlock(): error! block reference invalid. block name =", Name, " references=", Reference)
			os.Exit(-1)
		}
	}

	thisBlock.SizeOfPastSet = SizeOfPastSet(&thisBlock)

	//add this block to the chain
	chain[Name] = &thisBlock
	return &thisBlock
}
