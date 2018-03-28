package main

import (
	"os"
	"fmt"
	"sort"
	."phantom/phantom"
)

var chain map[string]*Block

func chainAddBlock(Name string, References []string, chain map[string]*Block) *Block{

	//create this block
	thisBlock := Block{Name, -1, make(map[string]*Block), make(map[string]*Block),make(map[string]bool)}

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

	//add this block to the chain
	chain[Name] = &thisBlock
	return &thisBlock
}

func chainInitialize() map[string]*Block{

	//initial an empty chain
	chain = make(map[string]*Block)

	//add blocks

	chainAddBlock("Genesis", []string{}, chain)

	chainAddBlock("B", []string{"Genesis"}, chain)
	chainAddBlock("C", []string{"Genesis"}, chain)
	chainAddBlock("D", []string{"Genesis"}, chain)
	chainAddBlock("E", []string{"Genesis"}, chain)

	chainAddBlock("F", []string{"B","C"}, chain)
	chainAddBlock("I", []string{"C","D"}, chain)
	chainAddBlock("H", []string{"E"}, chain)

	chainAddBlock("J", []string{"F","D"}, chain)
	chainAddBlock("L", []string{"F"}, chain)
	chainAddBlock("K", []string{"J","I","E"}, chain)
	chainAddBlock("N", []string{"D","H"}, chain)

	chainAddBlock("M", []string{"L","K"}, chain)
	chainAddBlock("O", []string{"K"}, chain)
	chainAddBlock("P", []string{"K"}, chain)
	chainAddBlock("Q", []string{"N"}, chain)

	chainAddBlock("R", []string{"O","P","N"}, chain)

	chainAddBlock("S", []string{"Q"}, chain)
	chainAddBlock("T", []string{"S"}, chain)
	chainAddBlock("U", []string{"T"}, chain)

	chainAddBlock("Virtual", []string{"M","R", "U"}, chain)

	fmt.Println("chainInitialize(): done. blocks=", len(chain)-1)

	return chain
}



func main() {

	chainInitialize()

	CalcBlue(chain, 3, chain["Virtual"])

	// print the result of blue sets

	names := make([]string, len(chain))
	i := 0
	for name := range chain {
		names[i] = name
		i++
	}
	sort.Strings(names)

	fmt.Print("CalcBlue(): done completely. blue blocks = (G).")
	for _, name := range names {
		if name=="Genesis" || name=="Virtual" {
			continue
		}

		block := chain[name]
		for _,blue := range block.Blue {
			if blue==true {
				fmt.Print(block.Name, ".")
			}
		}
	}
	fmt.Println("(V)")
}
