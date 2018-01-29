// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

func Test_Page_Flatten(t *testing.T) {
	// GIVEN a tree of data
	baseUrl, _ := url.Parse("http://test/")
	root, err := LoadRoot("test_content", baseUrl)
	if err != nil {
		t.Fatalf("LoadRoot failed: %v", err)
	}

	// WHEN the content is flattened
	pages := root.Flatten()

	// THEN the count should be more than just the root page
	if len(pages) <= 1 {
		t.Fatalf("The number of loaded pages should be more than 1 was %v", len(pages))
	}

	var apple *Page = nil
	var rock *Page = nil

	for p := range pages {
		page := pages[p]
		if page == root {
			t.Errorf("The root page was in among the content pages, it shouldn't be included. %d", p)
		}
		if page.Title == "Apple Tree Test Page" {
			apple = page
		}
		if page.Title == "Rock Test Page" {
			rock = page
		}
	}
	if apple == nil {
		t.Errorf("Couldn't find Apple page %v", apple)
	}
	if rock == nil {
		t.Errorf("Couldn't find Rock page %v", rock)
	}
}

func Test_Page_Flatten_listHidden(t *testing.T) {
	baseUrl, _ := url.Parse("http://test/")
	root, err := LoadRoot("test_content", baseUrl)
	if err != nil {
		t.Fatalf("LoadRoot failed: %v", err)
	}

	// WHEN the content is flattened
	pages := root.Flatten()
	for i := range pages {
		if pages[i].ListHidden {
			t.Errorf("Page should not be listed: %v", pages[i].URL)
		}
		if pages[i].Dir {
			t.Errorf("Dir should not be listed: %v", pages[i].URL)
		}
	}
}
