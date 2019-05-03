// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib

import (
	"encoding/json"
	"strings"

	"github.com/dotchain/dot/changes"
	"github.com/dotchain/dot/changes/types"
)

// Compact implements some helper routines for working with the
// compact form as defined in CompactJSON.md
type Compact struct{}

var specials = "[{(=)}]+:,"

// Decode takes a compact form string and converts it to a changes.Change
func (Compact) Decode(s string) (string, changes.Change) {
	if s == "" {
		return "", nil
	}

	l, r := strings.Index(s, "("), strings.LastIndex(s, ")")
	left, mid, right := s[:l], s[l+1:r], s[r+1:]
	if strings.Contains(mid, "=") {
		e := strings.Index(mid, "=")
		before, after := mid[:e], mid[e+1:]
		offset := types.S16(left).Count()
		input := left + before + right
		splice := changes.Splice{
			Offset: offset,
			Before: types.S16(before),
			After:  types.S16(after),
		}
		return input, splice
	}

	if strings.Contains(left, "=") {
		parts := strings.Split(left, "=")
		input := parts[0] + parts[1] + mid + right
		offset := types.S16(parts[0] + parts[1]).Count()
		distance := -types.S16(parts[1]).Count()
		count := types.S16(mid).Count()
		move := changes.Move{Offset: offset, Count: count, Distance: distance}
		return input, move
	}

	if strings.Contains(right, "=") {
		parts := strings.Split(right, "=")
		input := left + mid + parts[0] + parts[1]
		offset := types.S16(left).Count()
		distance := types.S16(parts[0]).Count()
		count := types.S16(mid).Count()
		move := changes.Move{Offset: offset, Count: count, Distance: distance}
		return input, move
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
func (c Compact) Apply(input string, ch changes.Change) string {
	return string(types.S16(input).Apply(nil, ch).(types.S16))
}

// Encode takes an input and a set of changes and converts it into the
// compact form
func (c Compact) Encode(input string, ch changes.Change) []string {
	result := []string(nil)
	if cs, ok := ch.(changes.ChangeSet); ok {
		for _, cx := range cs {
			result = append(result, c.Encode(input, cx)...)
			input = c.Apply(input, cx)
		}
	} else {
		result = append(result, c.Encode1(input, ch))
	}

	filtered := []string(nil)
	for _, s := range result {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// Encode1 is like Encode but it only takes one change
func (c Compact) Encode1(input string, ch changes.Change) string {
	if ch == nil {
		return ""
	}

	u := types.S16(input)
	switch ch := ch.(type) {
	case changes.Splice:
		left := string(u.Slice(0, ch.Offset).(types.S16))
		before := string(ch.Before.(types.S16))
		after := string(ch.After.(types.S16))
		end := ch.Offset + ch.Before.Count()
		right := string(u.Slice(end, u.Count()-end).(types.S16))
		return left + "(" + before + "=" + after + ")" + right
	case changes.Move:
		mid := string(u.Slice(ch.Offset, ch.Count).(types.S16))
		left := string(u.Slice(0, ch.Offset).(types.S16))
		end := ch.Offset + ch.Count
		len := u.Count()
		right := string(u.Slice(end, len-end).(types.S16))

		if ch.Distance < 0 {
			l1 := string(u.Slice(0, ch.Offset+ch.Distance).(types.S16))
			l2 := string(u.Slice(ch.Offset+ch.Distance, -ch.Distance).(types.S16))
			left = l1 + "=" + l2
		} else {
			r1 := string(u.Slice(end, ch.Distance).(types.S16))
			r2 := string(u.Slice(end+ch.Distance, len-end-ch.Distance).(types.S16))
			right = r1 + "=" + r2
		}

		return left + "(" + mid + ")" + right
	}
	panic(ch)
}
