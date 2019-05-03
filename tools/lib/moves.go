// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib

import "github.com/dotchain/dot/changes"

// Moves implements a bunch of useful utiities for working
// with moves
type Moves struct {
	Input string
}

// ForEach generates all possible move operations for the
// given input string and calls the provided callback with the
// op encoded as in this spec:
// http://github.com/dotchain/dataset/CompactJSON.md
func (m *Moves) ForEach(fn func(string)) {
	input := m.Input
	for offset := 0; offset <= len(input); offset++ {
		for end := offset; end <= len(input); end++ {
			for dest := 0; dest <= len(input); dest++ {
				if dest <= offset {
					fn(m.EncodeCompact(offset, end-offset, dest-offset))
				} else if dest >= end {
					fn(m.EncodeCompact(offset, end-offset, dest-end))
				}
			}
		}
	}
}

// EncodeCompact encodes a move into a compact format
func (m *Moves) EncodeCompact(offset, count, distance int) string {
	move := changes.Move{Offset: offset, Count: count, Distance: distance}
	return Compact{}.Encode1(m.Input, move)
}

// ForEachPair generates pairs of operations
func (m *Moves) ForEachPair(fn func(left, right string)) {
	m.ForEach(func(s1 string) {
		m.ForEach(func(s2 string) {
			fn(s1, s2)
		})
	})
}

// ForEachUniquePair generates only unique pairs of operations and
// uses the provided alphabet for the "uniqueness" calculation
func (m *Moves) ForEachUniquePair(alphabet []string, fn func(string, string, string)) {
	seen := map[string]map[string]string{}
	m.ForEachPair(func(l, r string) {
		normalized := Normalize([]string{m.Input, l, r}, specials, alphabet)
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
