// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"net/url"
	"testing"
)

func Test_Page_loadPage(t *testing.T) {
	p := &Page{URL: "http://test/"}

	// make sure the main index page loads
	err := p.loadSubPage("test_content/index.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Main Test Page" {
		t.Fatalf(`Couldn't load root page, title was wrong "%v"`, p.Title)
	}
}

func Test_Page_loadPage_url(t *testing.T) {
	p := &Page{URL: "http://test/apps/"}

	// GIVEN a folder
	if err := p.loadSubPage("test_content/tree/index.yaml"); err != nil {
		t.Fatal(err)
	}
	p.URL = "http://test/tree/"

	// WHEN the a page loaded in that folder
	if err := p.loadSubPage("test_content/tree/apple.yaml"); err != nil {
		t.Fatal(err)
	}

	if len(p.SubPages) != 1 {
		t.Fatal("The sub-page was not added.")
	}

	// THEN it's URL should be correct
	sp := p.SubPages[0]
	if sp.URL != "http://test/tree/apple/" {
		t.Errorf("Page URL was not correct: %v", sp.URL)
	}
}

func Test_Page_loadPage_Layout(t *testing.T) {
	// GIVEN the test content
	u, err := url.Parse("http://test/")

	// WHEN the content is loaded
	root, err := LoadRoot("test_content", u)
	if err != nil {
		t.Error(err)
	}
	m := MapPages(root)

	// THEN the layouts in each page should be overridden properly
	type test struct {
		path   string
		layout string
	}
	tests := []test{
		{"/", "dir.html"},
		{"/rock/", "page.html"},
		{"/tree/", "dir.html"},
		{"/tree/apple/", "page.html"},
		{"/noindex/", "dir.html"},
		{"/noindex/something/", "page.html"},
		{"/override/", "overridden.html"},
	}
	for i := range tests {
		p, ok := m[tests[i].path]
		if !ok {
			t.Errorf("Couldn't find %v in map", tests[i].path)
			continue
		}
		if p.Layout != tests[i].layout {
			t.Errorf("Expecting layout %v got %v on %v", tests[i].layout, p.Layout, tests[i].path)
		}
	}
}
