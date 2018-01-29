// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import "testing"

func Test_fileBase(t *testing.T) {

	type test struct {
		filename string
		expected string
	}

	tests := []test{
		{"/blah/blah your mom/file.xyz", "file"},
		{"C:\\Test\\correct.ext", "correct"},
		{"/dir/slash.xyz/", "slash"},
		{"", ""},
		{"///", ""},
	}

	var base string
	for i := range tests {
		base = fileBase(tests[i].filename)
		if base != tests[i].expected {
			t.Errorf(`Expecting "%v" from "%v", got "%v"`, tests[i].expected, tests[i].filename, base)
		}
	}
}
