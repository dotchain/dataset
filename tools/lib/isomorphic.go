// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib

import (
	"sort"
	"strings"
)

// Normalize takes a bunch of strings and uses a normalization
// procedure to return a new set of strings.
//
// Normalization works by identifying substrings that only ever appear
// together in all the provided input strings.  (Punctuations are
// considered as effectively breaking strings).  It then assigns a
// unique letter for each prefix -- the letter is assigned in order in
// which the prefix appears (those in earlier strings are considered
// earlier prefixes and those in latter strings are considerd latter;
// and within a string, an earlier prefix gets an earlier alphabet)
//
// This process indirectly helps with identifying if two sequences of
// strings are "isomorphic" by comparing if their normalized forms are
// the same.
//
// Note that this process is quite fragile.  For example, it is not
// legal to have a string that has any characters duplicated.  If that
// happens, the behavior is random.
func Normalize(s []string, punctuations string, alphabet []string) []string {
	segments := []segment{}
	for _, str := range getAllCommonSubsequences(s, punctuations) {
		for index := range s {
			offset := strings.Index(s[index], str)
			if offset >= 0 {
				seg := segment{str, "", index, offset}
				segments = append(segments, seg)
				break
			}
		}
	}

	sort.Sort(bySegment(segments))
	for kk := range segments {
		segments[kk].replacement = alphabet[kk]
	}
	result := make([]string, len(s))
	for kk := range s {
		result[kk] = bySegment(segments).replace(s[kk])
	}
	return result
}

type segment struct {
	s, replacement string
	index, offset  int
}

type bySegment []segment

func (b bySegment) Len() int      { return len(b) }
func (b bySegment) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b bySegment) Less(i, j int) bool {
	if b[i].index == b[j].index {
		return b[i].offset < b[j].offset
	}
	return b[i].index < b[j].index
}

func (b bySegment) replace(str string) string {
	if str == "" {
		return ""
	}

	for kk := range b {
		if strings.HasPrefix(str, b[kk].s) {
			return b[kk].replacement + b.replace(str[len(b[kk].s):])
		}
	}
	return str[:1] + b.replace(str[1:])
}

func getAllCommonSubsequences(s []string, punctuations string) []string {
	seg := map[string]bool{}

	var addMany func(parts []string)
	addMany = func(parts []string) {
		if len(parts) == 0 {
			return
		}
		if len(parts) > 1 {
			addMany(parts[1:])
		}

		p := parts[0]
		if p == "" || seg[p] {
			return
		}

		for key := range seg {
			if strings.Contains(key, p) {
				affixes := append(strings.Split(key, p), p)
				delete(seg, key)
				addMany(affixes)
				return
			}

			if strings.Contains(p, key) {
				addMany(strings.Split(p, key))
				return
			}
		}
		seg[p] = true
	}

	for _, str := range s {
		addMany(strings.FieldsFunc(str, func(r rune) bool {
			return strings.IndexRune(punctuations, r) >= 0
		}))
	}

	result := []string{}
	for key := range seg {
		result = append(result, key)
	}
	return result
}
