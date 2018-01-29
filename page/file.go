// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
)

func (p *Page) loadSubPage(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("Couldn't load page %v: %v", filename, err)
		return
	}
	defer file.Close()
	buf := bytes.Buffer{}
	if _, err = buf.ReadFrom(file); err != nil {
		err = fmt.Errorf("Couldn't read bytes from page %v: %v", filename, err)
		return
	}
	base := fileBase(filename)
	subPage := &Page{
		Parent: p,
		Layout: "page.html",
	}
	if base == "index" {
		subPage.Layout = "dir.html"
	}
	if err = yaml.Unmarshal(buf.Bytes(), subPage); err != nil {
		err = fmt.Errorf("Couldn't decode yaml for page %v: %v", filename, err)
		return
	}
	if len(subPage.URL) > 0 {
		subPage.External = true
	} else {
		subPage.addRelativeUrl(p, fileBase(filename)+"/")
	}
	if base == "index" {
		p.copyIndex(subPage)
	} else {
		p.SubPages = append(p.SubPages, subPage)
	}
	return
}

func (p *Page) addRelativeUrl(root *Page, addition string) {
	// make a copy of the root url
	u, _ := url.Parse(root.URL)
	u.Path += addition
	p.URL = u.String()
}
