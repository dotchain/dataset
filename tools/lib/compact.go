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

// Compact implements some helper routines for working with the
// compact form as defined in CompactJSON.md
type Compact struct{}

// Decode takes a compact form string and converts it to a dot.Change
func (Compact) Decode(s string) (string, dot.Change) {
	l, r := strings.Index(s, "("), strings.LastIndex(s, ")")
	left, mid, right := s[:l], s[l+1:r], s[r+1:]
	if strings.Contains(mid, "=") {
		e := strings.Index(mid, "=")
		before, after := mid[:e], mid[e+1:]
		offset := len(utf16.Encode([]rune(left)))
		input := left + before + right
		splice := &dot.SpliceInfo{Offset: offset, Before: before, After: after}
		return input, dot.Change{Splice: splice}
	}

	if strings.Contains(left, "=") {
		parts := strings.Split(left, "=")
		input := parts[0] + parts[1] + mid + right
		offset := len(utf16.Encode([]rune(parts[0] + parts[1])))
		distance := -len(utf16.Encode([]rune(parts[1])))
		count := len(utf16.Encode([]rune(mid)))
		move := &dot.MoveInfo{Offset: offset, Count: count, Distance: distance}
		return input, dot.Change{Move: move}
	}

	if strings.Contains(right, "=") {
		parts := strings.Split(right, "=")
		input := left + mid + parts[0] + parts[1]
		offset := len(utf16.Encode([]rune(left)))
		distance := len(utf16.Encode([]rune(parts[0])))
		count := len(utf16.Encode([]rune(mid)))
		move := &dot.MoveInfo{Offset: offset, Count: count, Distance: distance}
		return input, dot.Change{Move: move}
	}

	panic("Unknown formatted string")
}

// Stringify converts a string-like value to string
func (c Compact) Stringify(x interface{}) string {
	if x == nil {
		return ""
	}

	var ret string
	if e, err := json.Marshal(x); err != nil {
		panic(err)
	} else if err := json.Unmarshal(e, &ret); err != nil {
		panic(err)
	}
	return ret
}

// Apply takes a single change and applies to the input string
func (c Compact) Apply(input string, ch dot.Change) string {
	x := dot.Utils(dot.Transformer{})
	return c.Stringify(x.Apply(input, []dot.Change{ch}))
}

// ApplyMany sequentially applies the set of changes
func (c Compact) ApplyMany(input string, changes []dot.Change) string {
	for _, ch := range changes {
		input = c.Apply(input, ch)
	}
	return input
}

// Encode takes an input and a set of changes and converts it into the
// compact form
func (c Compact) Encode(input string, changes []dot.Change) []string {
	result := make([]string, len(changes))
	for kk, ch := range changes {
		result[kk] = c.Encode1(input, ch)
		input = c.Apply(input, ch)
	}
	return result
}

// Encode1 is like Encode but it only takes one change
func (c Compact) Encode1(input string, ch dot.Change) string {
	u := utf16.Encode([]rune(input))
	if ch.Splice != nil {
		left := string(utf16.Decode(u[:ch.Splice.Offset]))
		before := c.Stringify(ch.Splice.Before)
		after := c.Stringify(ch.Splice.After)
		right := string(utf16.Decode(u[ch.Splice.Offset+len(utf16.Encode([]rune(before))):]))
		return left + "(" + before + "=" + after + ")" + right
	}
	if ch.Move != nil {
		mid := string(utf16.Decode(u[ch.Move.Offset : ch.Move.Offset+ch.Move.Count]))
		left := string(utf16.Decode(u[:ch.Move.Offset]))
		right := string(utf16.Decode(u[ch.Move.Offset+ch.Move.Count:]))

		if ch.Move.Distance < 0 {
			l1 := string(utf16.Decode(u[:ch.Move.Offset+ch.Move.Distance]))
			l2 := string(utf16.Decode(u[ch.Move.Offset+ch.Move.Distance : ch.Move.Offset]))
			left = l1 + "=" + l2
		} else {
			end := ch.Move.Offset + ch.Move.Count
			r1 := string(utf16.Decode(u[end : end+ch.Move.Distance]))
			r2 := string(utf16.Decode(u[end+ch.Move.Distance:]))
			right = r1 + "=" + r2
		}

		return left + "(" + mid + ")" + right
	}
	panic("Unknown op")
}
