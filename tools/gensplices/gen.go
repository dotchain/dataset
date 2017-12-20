// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// main generates json/compact/splices.json.  Please see
// github.com/dotchain/dataset
package main

import (
	"encoding/json"
	"fmt"
	"github.com/dotchain/dataset/tools/lib"
	"github.com/dotchain/dot"
	"log"
	"strings"
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

	input := "abcdefg"
	// Note that the alphabet is deliberately unicode to make sure
	// that the tests work properly with this.
	alphabet := strings.Split("𝐀𝐁𝐂𝐃𝐄𝐅𝐆𝐇𝐈𝐉𝐊𝐋𝐌𝐍𝐎𝐏", "")

	t := dot.Transformer{}
	x := &lib.Splices{Input: input, Inserts: []string{"", "xyz", "XYZ"}}
	first := ""
	compact := lib.Compact{}
	x.ForEachUniquePair(alphabet, func(input, left, right string) {
		inputl, l := compact.Decode(left)
		inputr, r := compact.Decode(right)
		if inputl != inputr || input != inputl {
			log.Fatal("Invalid inputs", inputl, inputr, left, right)
		}
		lx, rx := []dot.Change{l}, []dot.Change{r}
		mergedl, mergedr := t.MergeChanges(lx, rx)
		allLeft := append([]dot.Change{l}, mergedl...)
		allRight := append([]dot.Change{r}, mergedr...)

		encodedl := compact.Encode(input, allLeft)
		encodedr := compact.Encode(input, allRight)

		outputl := compact.ApplyMany(input, allLeft)
		outputr := compact.ApplyMany(input, allRight)
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
