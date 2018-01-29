// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"testing"
)

func Test_Site_loadTemplatesAndContent(t *testing.T) {
	testLog := log.New(&bytes.Buffer{}, "", 0)
	u, _ := url.Parse("http://example.com")
	s, err := New(u, _templatePath, false, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		contentPath string
		err         string
	}
	tests := []test{
		{"../example/content", "<nil>"},
		{"noexist", "Couldn't open content folder noexist: open noexist: no such file or directory"},
	}
	for i := range tests {
		s.contentPath = tests[i].contentPath
		err := s.loadTemplatesAndContent()
		if fmt.Sprintf("%v", err) != tests[i].err {
			t.Errorf("%s expected: %v got: %v", tests[i].contentPath, tests[i].err, err)
		}
	}
}
