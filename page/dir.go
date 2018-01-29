// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"fmt"
	"os"
)

func (p *Page) loadSubDir(dirName string) (err error) {
	subDir, err := os.Open(dirName)
	if err != nil {
		err = fmt.Errorf("Couldn't read subdir %v %v", dirName, err)
		return
	}
	defer subDir.Close()
	subPage := &Page{
		Title:  dirName,
		Parent: p,
		Dir:    true,
		Layout: "dir.html",
	}
	subPage.addRelativeUrl(p, fileBase(dirName)+"/")
	err = subPage.addAllPages(subDir)
	if err != nil {
		err = fmt.Errorf("Couldn't load subpages for %v: %v", dirName, err)
	}
	sortDateDesc(subPage.SubPages)
	p.SubPages = append(p.SubPages, subPage)
	return
}
