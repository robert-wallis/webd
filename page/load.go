// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// LoadRoot takes a root folder and returns a Page with it's full tree of content loaded in sub-pages.
func LoadRoot(contentFolder string, baseUrl *url.URL) (root *Page, err error) {
	dir, err := os.Open(contentFolder)
	if err != nil {
		err = fmt.Errorf("Couldn't open content folder %v: %v", contentFolder, err)
		return
	}
	defer dir.Close()
	root = &Page{
		URL: baseUrl.String(),
	}
	err = root.addAllPages(dir)
	return
}

func (p *Page) addAllPages(dir *os.File) (err error) {
	fileInfos, err := dir.Readdir(0)
	if err != nil {
		err = fmt.Errorf("Couldn't read %v folder: %v", dir.Name(), err)
		return
	}
	for f := range fileInfos {
		fileInfo := fileInfos[f]
		filename := fileInfo.Name()
		fullFilename := filepath.Join(dir.Name(), filename)
		if fileInfo.IsDir() {
			if err = p.loadSubDir(fullFilename); err != nil {
				return
			}
			continue
		}
		if filepath.Ext(filename) == ".yaml" {
			err = p.loadSubPage(fullFilename)
			if err != nil {
				err = fmt.Errorf("Couldn't load page: %v %v", filename, err)
				return
			}
		}
	}
	return
}
