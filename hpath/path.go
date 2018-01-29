// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// HPaths are like XPath for HTML, with less features.
package hpath

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// HtmlNode wraps an html.Node to add HPath functionality.
type HtmlNode struct {
	html.Node
}

// Path returns a value in an HTML document.
// Using the path: /html/head/title
// This will return the <title> node in an HTML document.  Or an error if it couldn't find it.
func (n *HtmlNode) Path(path string) (value string, err error) {
	var current *html.Node = &n.Node
	dir := strings.Split(path, "/")
	var d int = 0
	for {
		if d >= len(dir) {
			err = fmt.Errorf("Couldn't find path %v", path)
			return
		}
		d_name := dir[d]

		if len(d_name) == 0 {
			// skip empty paths
			d++
			continue
		}

		if current == nil {
			err = fmt.Errorf("Couldn't find path %v at %v %d", path, d_name, d)
			return
		}

		if current.Type == html.DocumentNode {
			current = current.FirstChild
			continue
		}

		if current.Type == html.DoctypeNode {
			current = current.NextSibling
			continue
		}

		if current.Type != html.ElementNode {
			current = current.NextSibling
			continue
		}

		element := d_name
		attr := ""
		if split := strings.Index(d_name, "@"); split != -1 {
			element = d_name[:split]
			attr = d_name[split+1:]
		}

		if element == current.Data {
			if d == len(dir)-1 {
				if len(attr) > 0 {
					for a := range current.Attr {
						at := current.Attr[a]
						if at.Key == attr {
							// SUCCESS attr found
							value = at.Val
							return
						}
					}
					err = fmt.Errorf("Couldn't find attribute %v in %v HTML path %v", attr, element, path)
					return
				}
				// SUCCESS element found
				value = current.Data
				return
			}

			// going up
			d++
			current = current.FirstChild
			continue
		}
		// not the right element at this level, try next sibling
		current = current.NextSibling
	}
}
