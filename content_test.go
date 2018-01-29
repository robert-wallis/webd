// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package main

import (
	"bytes"
	"fmt"
	"github.com/robert-wallis/webd/hpath"
	"github.com/robert-wallis/webd/site"
	"golang.org/x/net/html"
	"log"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var _testPath = "./example"

func Test_Handler_content(t *testing.T) {
	// GIVEN the index page has been shown
	hostname, _ := url.Parse("http://test")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, _ := site.New(hostname, _testPath, false, false, testLog, testLog)
	req := httptest.NewRequest("GET", hostname.String()+"/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	// WHEN the body is evaluated for specific content strings
	contentStrings := []string{
		"<title>Example.com</title>",
		"<h1>Example.com</h1>",
		"This is an example site.",
		"Mix Body",
		"Just HTML Body",
		"Trip",
		"/thumbs/trip.jpg",
		"A lame site",
		"Robert's Dot Files",
		"https://github.com/robert-wallis/dotfiles",
		"test@example.com",
		"@example",
		"twitter.com/example",
		"github.com/robert-wallis",
		"_gaq.push(['_setAccount', googleAnalyticsId])",
		fmt.Sprintf("Copyright © %d YOUR NAME", time.Now().Year()),
	}

	// THEN each of the strings should be rendered on the page
	body := w.Body.String()
	for e := range contentStrings {
		expected := contentStrings[e]
		if !strings.Contains(body, expected) {
			t.Errorf("Expected but not found in the index page for example/ : \"%v\"", expected)
		}
	}
	if t.Failed() {
		t.Logf("body:", body)
	}
}

func Test_HTML_rightness(t *testing.T) {
	// GIVEN the index page is shown
	hostname, _ := url.Parse("https://test")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, _ := site.New(hostname, _testPath, false, false, testLog, testLog)
	req := httptest.NewRequest("GET", hostname.String()+"/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	// WHEN checking for HTML correctness
	type htmlCheck struct {
		path  string
		value string
	}
	htmlChecks := []htmlCheck{
		{"/html/head/base@href", hostname.String()},
	}

	// THEN the HTML should be correct
	root, err := html.Parse(w.Body)
	if err != nil {
		t.Fatal("Couldn't parse HTML from root", err)
	}
	pathRoot := &hpath.HtmlNode{Node: *root}
	for e := range htmlChecks {
		expected := htmlChecks[e]
		value, err := pathRoot.Path(expected.path)
		if err != nil {
			t.Errorf("Expected HTML path not found :\"%v\" %v", expected.path, err)
		}
		if value != expected.value {
			t.Errorf("Expected HTML value \"%v\" in %v but was %v", expected.value, expected.path, value)
		}
	}
}
