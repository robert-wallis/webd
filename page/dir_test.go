// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"testing"
)

func Test_Page_loadDir(t *testing.T) {

	// should error if it can't open the folder
	p := &Page{URL: "http://test/"}
	err := p.loadSubDir("noexist")
	if err == nil {
		t.Fatal("Should have failed to load non-existing tree.")
	}

	// should not error on directores that exist
	p = &Page{URL: "http://test/"}
	err = p.loadSubDir("test_content")
	if err != nil {
		t.Fatal("Couldn't load conent test tree.", err)
	}
}
