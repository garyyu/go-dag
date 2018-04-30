
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

import "testing"

var testBlock *Block = nil

func BenchmarkBlocks(b *testing.B) {


	for i := 0; i < b.N; i++ {
		testBlock = new(Block)
	}
}

