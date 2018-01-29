// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"container/list"
	"net/url"
)

// MapPages gives a path to sub-page listing of the Page tree.
func MapPages(root *Page) map[string]*Page {
	m := make(map[string]*Page)
	for current := range Walk(root) {
		if current.External {
			// skip external links
			continue
		}
		if u, err := url.Parse(current.URL); err == nil {
			path := u.Path
			if path == "" {
				path = "/"
			}
			m[path] = current
		}
	}
	return m
}

// Walk returns channel of every Page and SubPage by walking the tree breadth-first.
func Walk(root *Page) chan *Page {
	out := make(chan *Page)
	go func() {
		todo := list.New()
		todo.PushBack(root)
		for todo.Len() > 0 {
			c := todo.Front()
			todo.Remove(c)
			current := c.Value.(*Page)
			out <- current
			for s := range current.SubPages {
				todo.PushBack(current.SubPages[s])
			}
		}
		close(out)
	}()
	return out
}
