// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"bytes"
	"log"
	"net/url"
	"testing"
)

var _templatePath = "../example"

func Test_Site_New(t *testing.T) {
	// GIVEN valid starting params
	u, _ := url.Parse("http://localhost:8009")

	// WHEN the server is created
	testLog := log.New(&bytes.Buffer{}, "", 0)
	_, err := New(u, _templatePath, false, false, testLog, testLog)

	// THEN it shouldn't generate an error
	if err != nil {
		t.Error(err)
	}
}

func Test_Site_New_badPath(t *testing.T) {
	// GIVEN valid starting inputs
	u, _ := url.Parse("http://localhost:8009")

	// WHEN the test path is invalid
	templatesPath := "noexist"
	testLog := log.New(&bytes.Buffer{}, "", 0)
	_, err := New(u, templatesPath, false, false, testLog, testLog)

	// THEN it should fail
	if err == nil {
		t.Error("Should have failed.")
	}
}

func Test_Site_contentPage(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		location string
		title    string
		found    bool
		redirect bool
	}
	tests := []test{
		{"/blog/trip/", "Trip", true, false},
		{"/blog/", "Post List Example", true, false},
		{"/blog", "Post List Example", true, true},
		{"/", "Example.com", true, false},
		{"/noexist", "", false, false},
	}
	for i := range tests {
		p, found, redirect := s.contentPage(tests[i].location)
		if found != tests[i].found {
			t.Errorf("%v expected found %v got %v", tests[i].location, tests[i].found, found)
		}
		if redirect != tests[i].redirect {
			t.Errorf("%v expected redirect %v got %v", tests[i].location, tests[i].redirect, redirect)
		}
		if found && p.Title != tests[i].title {
			t.Errorf("%v expected title %v got %v", tests[i].location, tests[i].title, p.Title)
		}
	}
}
