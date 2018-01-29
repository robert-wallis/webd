// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"container/list"
)

func (p *Page) Breadcrumbs() (out []*Page) {
	l := list.New()
	var c = p
	for c != nil {
		l.PushFront(c)
		c = c.Parent
	}
	for e := l.Front(); e != nil; e = e.Next() {
		out = append(out, e.Value.(*Page))
	}
	return
}
