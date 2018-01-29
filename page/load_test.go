// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

// Load the content directory, and make sure the expected data is loaded.
func Test_Page_LoadRoot(t *testing.T) {
	baseUrl, _ := url.Parse("http://test/")
	root, err := LoadRoot("test_content", baseUrl)
	if err != nil {
		t.Fatal(err)
	}
	if root == nil {
		t.Fatal("No content was loaded.")
	}
	if root.Title != "Main Test Page" {
		t.Fatal("Index page was not loaded.")
	}
	if len(root.SubTitle) == 0 {
		t.Fatal("Index page is missing a subtitle.")
	}
	for p := range root.SubPages {
		sub := root.SubPages[p]
		if sub.Title == "Main Test Page" {
			t.Fatal("No subpage of root should be the root subpage, index.yaml should be the root page")
		}
	}
}
