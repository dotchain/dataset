// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"github.com/dotchain/dot"
	"strings"
	"unicode/utf16"
)

// Splices implements a bunch of useful utiities for working
// with splices
type Splices struct {
	Input   string
	Inserts []string
}

// ForEach generates all possible splice operations for the
// given input string and calls the provided callback with the
// splice encoded as in this spec:
// http://github.com/dotchain/dataset/CompactJSON.md
func (s *Splices) ForEach(fn func(string)) {
	input, inserts := s.Input, s.Inserts
	for offset := 0; offset <= len(input); offset++ {
		for end := offset; end <= len(input); end++ {
			before := input[offset:end]
			for _, insert := range inserts {
				fn(s.EncodeCompact(offset, before, insert))
			}
		}
	}
}

// EncodeCompact encodes a splice into a compact format
func (s *Splices) EncodeCompact(offset int, before, after string) string {
	left := s.Input[:offset]
	right := s.Input[offset+len(before):]
	return left + "(" + before + "=" + after + ")" + right
}

// ForEachPair generates pairs of operations
func (s *Splices) ForEachPair(fn func(left, right string)) {
	s.ForEach(func(s1 string) {
		s.ForEach(func(s2 string) {
			fn(s1, s2)
		})
	})
}

// ForEachUniquePair generates only unique pairs of operations and
// uses the provided alphabet for the "uniqueness" calculation
func (s *Splices) ForEachUniquePair(alphabet []string, fn func(string, string, string)) {
	seen := map[string]map[string]string{}
	s.ForEachPair(func(l, r string) {
		normalized := Normalize([]string{s.Input, l, r}, "(=)", alphabet)
		input, left, right := normalized[0], normalized[1], normalized[2]
		if _, ok := seen[left]; !ok {
			seen[left] = map[string]string{}
		}
		if _, ok := seen[left][right]; ok {
			return
		}
		seen[left][right] = input
		fn(input, left, right)
	})
}

func (s *Splices) ToChange(str string) (string, dot.Change) {
	l, r := strings.Index(str, "("), strings.Index(str, ")")
	mid := strings.Split(str[l+1:r], "=")
	before, after := mid[0], mid[1]
	input := str[:l] + before + str[r+1:]
	offset := len(utf16.Encode([]rune(str[:l])))
	splice := &dot.SpliceInfo{Offset: offset, Before: before, After: after}
	return input, dot.Change{Splice: splice}
}

func (s *Splices) EncodeChanges(input string, c []dot.Change) []string {
	result := make([]string, len(c))
	for kk, ch := range c {
		u := utf16.Encode([]rune(input))
		before := s.Stringify(ch.Splice.Before)
		after := s.Stringify(ch.Splice.After)
		left := string(utf16.Decode(u[:ch.Splice.Offset]))
		right := string(utf16.Decode(u[ch.Splice.Offset+len(utf16.Encode([]rune(before))):]))

		result[kk] = left + "(" + before + "=" + after + ")" + right
		input = s.ApplyChange(input, ch)
	}
	return result
}

func (s *Splices) ApplyChange(input string, ch dot.Change) string {
	return s.Stringify(dot.Utils(dot.Transformer{}).Apply(input, []dot.Change{ch}))
}

func (s *Splices) Stringify(x interface{}) string {
	var ret string
	if e, err := json.Marshal(x); err != nil {
		panic(err)
	} else if err := json.Unmarshal(e, &ret); err != nil {
		panic(err)
	}
	return ret
}

func (s *Splices) ApplyChanges(input string, c []dot.Change) string {
	for _, ch := range c {
		input = s.ApplyChange(input, ch)
	}
	return input
}
