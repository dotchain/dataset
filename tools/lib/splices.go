// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib

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
	return left + "[" + before + ":" + after + "]" + right
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
		normalized := Normalize([]string{s.Input, l, r}, ":[]", alphabet)
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