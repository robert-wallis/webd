// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

func Test_Site_mapRedirects(t *testing.T) {
	// GIVEN a loaded set of content Pages
	u, _ := url.Parse("http://test/")
	root, err := LoadRoot("../page/test_content", u)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN the redirects are mapped
	m := MapRedirects(Walk(root), u)

	// THEN they should be the expected redirects
	type redirectTest struct {
		src string
		dst string
	}
	tests := []redirectTest{
		{"/apple.html", "/tree/apple/"},
		{"/apple/", "/tree/apple/"},
	}

	for i := range tests {
		test := tests[i]
		dst, ok := m[test.src]
		delete(m, test.src)
		if !ok {
			t.Errorf("Source %v not in redirect map", test.src)
			continue
		}
		if dst != test.dst {
			t.Errorf("Source %v expected to redirect to %v but was %v", test.src, test.dst, dst)
		}
	}

	for k, v := range m {
		t.Errorf("Unexpected redirect %v to %v in map.", k, v)
	}
}
