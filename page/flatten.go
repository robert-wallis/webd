// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import "sort"

// Takes all the sub-pages of a page, and their children, and returns them as a single list.
func (p *Page) Flatten() (list []*Page) {
	todo := []*Page{p}
	for len(todo) > 0 {
		current := todo[0]
		todo = todo[1:]
		if current != p && (len(current.Body) != 0 || !current.Dir) && !current.ListHidden {
			list = append(list, current)
		}
		for s := range current.SubPages {
			sub := current.SubPages[s]
			todo = append(todo, sub)
		}
	}
	sortDateDesc(list)
	return list
}

// sortDateDesc puts the newest pages on the top.
func sortDateDesc(list []*Page) {
	sort.Slice(list, func(i, j int) bool {
		return list[j].DateUpdated.Before(list[i].DateUpdated)
	})
}
