// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib_test

import (
	"fmt"
	"github.com/dotchain/dataset/tools/lib"
	"strings"
)

func ExampleSplicesForEachUniquePair() {
	s := lib.Splices{Input: "abcdefg", Inserts: []string{"", "xyz", "XYZ"}}
	alphabet := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	seen := map[string]bool{}
	s.ForEachUniquePair(alphabet, func(i, l, r string) {
		if seen[i+"-->"+l+r] {
			fmt.Println("Unexpected duplicate")
		}
		seen[i+"-->"+l+r] = true
	})
	fmt.Println("Number of unique pairs =", len(seen))

	// Output: Number of unique pairs = 515
}
