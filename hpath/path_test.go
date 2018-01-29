// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package hpath

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func Test_HtmlNode_Path(t *testing.T) {
	// GIVEN some HTML
	body := "<html><body><p></p></body></html>"

	// WHEN the hpath is /html/body/p
	root, _ := html.Parse(strings.NewReader(body))
	rootNode := &HtmlNode{*root}
	value, err := rootNode.Path("/html/body/p")

	// THEN the value should be the tag name
	if err != nil {
		t.Fatal(err)
	}
	if value != "p" {
		t.Errorf("Expected value p was \"%v\"", value)
	}
}

func Test_HtmlNode_Path_attribute(t *testing.T) {
	// GIVEN some HTML
	doc := "<html><body><div><p name=success></p></div></body></html>"
	htmlRoot, _ := html.Parse(strings.NewReader(doc))
	root := &HtmlNode{*htmlRoot}

	// WHEN the hpath is /html/doc/p@name
	value, err := root.Path("/html/body/div/p@name")

	// THEN the value should be the element value
	if err != nil {
		t.Fatal(err)
	}
	if value != "success" {
		t.Errorf("Expected value \"success\" was \"%v\"", value)
	}
}

func Test_HtmlNode_bad_data(t *testing.T) {
	// GIVEN an empty HTML document
	doc := ""
	htmlRoot, _ := html.Parse(strings.NewReader(doc))
	root := &HtmlNode{*htmlRoot}

	// WHEN Path is called
	value, err := root.Path("/html/doc/p@name")

	// THEN it should fail
	if err == nil {
		t.Errorf("Should have failed with blank doc but passed with value %v", value)
	}
}
