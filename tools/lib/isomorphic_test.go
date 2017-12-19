// Copyright (C) 2017 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package lib_test

import (
	"github.com/dotchain/dataset/tools/lib"
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	punctuations := ":[]"
	alphabet := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p"}

	tests := [][]string{
		{"he[lo:ro]", "[he:bo]lo", "a[b:c]", "[a:d]b"},
		{"hea:rt", "he:pam:rt", "ab:c", "a:dbe:c"},
	}

	for _, test := range tests {
		input := test[:len(test)/2]
		expected := test[len(test)/2:]
		actual := lib.Normalize(input, punctuations, alphabet)
		if !reflect.DeepEqual(expected, actual) {
			t.Error("Expected", expected, "but got", actual)
		}
	}
}
