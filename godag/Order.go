
// Copyright 2018 The godag Authors
// This file is part of the godag library.
//
// The g-dag library is free software: you can redistribute it and/or modify
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


func Intersection(A map[string]*Block, B map[string]*Block) map[string]*Block{

	intersection := make(map[string]*Block)

	for _, blockA := range A {
		for _, blockB := range B {
			if blockA.Name==blockB.Name {
				intersection[blockA.Name] = blockA
			}
		}
	}

	return intersection
}


func Order(chain map[string]*Block, k int) []string {

	// step 2
	topo_queue := make([]*Block, 0)
	map_queue := make(map[string]*Block)	// store the blocks already in queue

	// step 3
	ordered_list := make([]string, 0)
	map_list := make(map[string]int)		// store the order of the blocks
	var order = 0

	// step 4
	CalcBlue(chain, k, chain["Virtual"])

	// step 5
	topo_queue = append(topo_queue, chain["Genesis"])	// FIFO Queue Push
	map_queue["Genesis"] = chain["Genesis"]

	// step 6
	var B *Block = nil
	for len(topo_queue)>0 {

		B = topo_queue[0]				// step 7.	FIFO Queue Pop
		topo_queue = topo_queue[1:]		// discard top element of FIFO Queue

		// step 8
		ordered_list = append(ordered_list, B.Name)
		map_list[B.Name] = order
		order ++

		// step 9
		childrenB := make(map[string]*Block)
		futureSet(B, childrenB)

		// step 9'		in 'some' topological ordering: LTPQ
		ltpq := LTPQ(childrenB, true)

		for _, name := range ltpq {
			C := childrenB[name]		// step 9"

			if IsBlueBlock(C)==true {
				// step 10
				pastC := make(map[string]*Block)
				pastSet(C, pastC)

				anticoneB := antiCone(chain, B)
				intersection := Intersection(pastC, anticoneB)

				// step 10'		in 'some' topological ordering: LTPQ
				ltpq2 := LTPQ(intersection, true)

				for _, name := range ltpq2 {
					D := intersection[name]		// step 10"

					if _,okk := map_list[D.Name]; !okk{
						// step 11
						if _,ok3 := map_queue[D.Name]; !ok3 {	// the queue should avoid duplicating elements
							topo_queue = append(topo_queue, D)	// FIFO Queue Push
							map_queue[D.Name] = D
						}
					}
				}

				// step 12
				if _,ok4 := map_queue[C.Name]; !ok4 {	// the queue should avoid duplicating elements
					topo_queue = append(topo_queue, C)	// FIFO Queue Push
					map_queue[C.Name] = C
				}
			}
		}
	}

	return ordered_list
}