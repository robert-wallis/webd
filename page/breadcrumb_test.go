// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

func Test_Page_Breadcrumbs(t *testing.T) {
	u, _ := url.Parse("http://test/")
	root, err := LoadRoot("test_content", u)
	if err != nil {
		t.Fatal(err)
	}
	m := MapPages(root)

	type test struct {
		path        string
		breadcrumbs []string
	}
	tests := []test{
		{"/", []string{"http://test/"}},
		{"/tree/apple/", []string{"http://test/", "http://test/tree/", "http://test/tree/apple/"}},
	}
	for i := range tests {
		p, ok := m[tests[i].path]
		if !ok {
			t.Errorf("Path not in map %v", tests[i].path)
			continue
		}
		b := p.Breadcrumbs()
		if len(b) != len(tests[i].breadcrumbs) {
			got := []string{}
			for j := range b {
				got = append(got, b[j].URL)
			}
			t.Errorf("path %v expecting %v, got %v", tests[i].path, tests[i].breadcrumbs, got)
			break // don't spam Page structure for every page
		}
		for j := range b {
			expecting := tests[i].breadcrumbs[j]
			got := b[j].URL
			if expecting != got {
				t.Errorf("path %v expecting %v, got %v", tests[i].path, expecting, got)
				return
			}
		}
	}
}
