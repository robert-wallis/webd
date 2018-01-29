// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

func Test_MapPages(t *testing.T) {
	// GIVEN the contents
	base, _ := url.Parse("http://test/")
	root, err := LoadRoot("../page/test_content", base)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN the map is calculated
	m := MapPages(root)

	// THEN it should contain some key entries
	tests := []string{
		"/",
		"/noindex/",
		"/noindex/hidden/",
		"/noindex/something/",
		"/override/",
		"/rock/",
		"/tree/",
		"/tree/apple/",
	}

	for i := range tests {
		expected := tests[i]
		v, ok := m[expected]
		if ok {
			delete(m, expected)
			if v == nil || v.URL == "" {
				t.Errorf("URL was empty: %v, %v", expected, v)
				continue
			}
			u, _ := url.Parse(v.URL)
			if u == nil || u.Path != expected {
				t.Errorf("Expected path %v in page: %v", expected, v)
			}
		} else {
			t.Errorf("Path not in map %v", expected)
		}
	}
	for k, v := range m {
		t.Errorf("Unexpected extra k/v in map - %v: %v", k, v.URL)
	}

	sadtests := []string{
		"/robert-wallis/webd",
		"/robert-wallis/webd/",
	}

	for i := range sadtests {
		missing := sadtests[i]
		v, ok := m[missing]
		if ok {
			t.Errorf("should not map external links %v: %v", missing, v)
		}
	}
}
