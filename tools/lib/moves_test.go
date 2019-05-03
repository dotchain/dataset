// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib_test

import (
	"fmt"
	"strings"

	"github.com/dotchain/dataset/tools/lib"
)

func ExampleMovesForEachUniquePair() {
	m := lib.Moves{Input: "abcdefg"}
	alphabet := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	seen := map[string]bool{}
	m.ForEachUniquePair(alphabet, func(i, l, r string) {
		if seen[i+"-->"+l+r] {
			fmt.Println("Unexpected duplicate")
		}
		seen[i+"-->"+l+r] = true
	})
	fmt.Println("Number of unique pairs =", len(seen))

	// Output: Number of unique pairs = 2919
}

func ExampleMoves_ForEachMoveSplicePair() {
	s := lib.Splices{Input: "abcdefg", Inserts: []string{"", "xyz"}}
	m := lib.Moves{Input: s.Input}
	alphabet := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	seen := map[string]bool{}

	s.ForEach(func(left string) {
		m.ForEach(func(right string) {
			normalized := lib.Normalize([]string{s.Input, left, right}, "(=)", alphabet)
			i, l, r := normalized[0], normalized[1], normalized[2]
			key := i + "|" + l + "|" + r
			seen[key] = true
		})
	})
	fmt.Println("Number of unique pairs =", len(seen))

	// Output: Number of unique pairs = 1006
}
