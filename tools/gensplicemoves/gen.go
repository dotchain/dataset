// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// main generates json/compact/splicemoves.json.  Please see
// github.com/dotchain/dataset
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dotchain/dataset/tools/lib"
	"github.com/dotchain/dot/changes"
)

var preamble = `
{
	"format": "compact",
	"test": [
`
var postamble = `
	]
}
`

func main() {
	fmt.Println(preamble)
	defer fmt.Println(postamble)

	input := "abcdefgh"
	// Note that the alphabet is deliberately unicode to make sure
	// that the tests work properly with this.
	alphabet := strings.Split("𝐀𝐁𝐂𝐃𝐄𝐅𝐆𝐇𝐈𝐉𝐊𝐋𝐌𝐍𝐎𝐏", "")

	splices := &lib.Splices{Input: input, Inserts: []string{"", "xyz"}}

	first := ""
	compact := lib.Compact{}

	splices.ForEachUniqueSpliceMovePair(alphabet, func(input, left, right string) {
		inputl, l := compact.Decode(left)
		inputr, r := compact.Decode(right)
		if inputl != inputr || input != inputl {
			log.Fatal("Invalid inputs", inputl, inputr, left, right)
		}
		mergedl, mergedr := changes.Merge(l, r)
		allLeft := changes.ChangeSet{l, mergedl}
		allRight := changes.ChangeSet{r, mergedr}

		encodedl := compact.Encode(input, allLeft)
		encodedr := compact.Encode(input, allRight)

		outputl := compact.Apply(input, allLeft)
		outputr := compact.Apply(input, allRight)
		if outputl != outputr {
			log.Fatal("merge failure: ", input, "\n", left, " x ", right, "\n", encodedl, " x ", encodedr, "\n", outputl, " x ", outputr)
		}

		output := outputl
		encoded, err := json.Marshal([]interface{}{
			input,
			output,
			[]string{left},
			[]string{right},
			encodedl[1:],
			encodedr[1:],
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\t\t%s", first, string(encoded))
		first = ",\n"
	})
}
