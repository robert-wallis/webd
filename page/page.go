// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// Page manages a website content based on YAML files.
package page

import (
	"time"
)

type Page struct {
	Title       string
	SubTitle    string
	URL         string
	External    bool
	Thumbnail   string
	DateUpdated time.Time
	Parent      *Page
	SubPages    []*Page
	Dir         bool
	Layout      string
	Redirects   []string
	Body        []map[string]string
	ListHidden  bool
}

// copyIndex takes the contents of src and puts them in the page.
// This is used for index.yaml pages in the root of the directory, the root is p, and the index.yaml is src.
func (p *Page) copyIndex(src *Page) {
	p.Title = src.Title
	p.SubTitle = src.SubTitle
	if len(p.URL) == 0 {
		p.URL = src.URL
	}
	p.Thumbnail = src.Thumbnail
	p.DateUpdated = src.DateUpdated
	if len(p.SubPages) == 0 && len(src.SubPages) > 0 {
		p.SubPages = src.SubPages
	}
	p.Dir = p.Dir || src.Dir
	p.Layout = src.Layout
	p.Redirects = src.Redirects
	p.Body = src.Body
	p.ListHidden = src.ListHidden
}
