// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"fmt"
	"github.com/robert-wallis/webd/page"
	"html/template"
	"time"
)

func (s *Site) loadTemplatesAndContent() (err error) {
	templatesCompiled, err := s.loadTemplates()
	if err != nil {
		return
	}

	root, err := s.loadContent()
	if err != nil {
		return
	}

	// saving only if successful
	s.templates = templatesCompiled
	s.pageRoot = root
	s.pageMap = page.MapPages(root)
	s.redirectMap = page.MapRedirects(page.Walk(root), s.base)
	return
}

// Builds all the templates in the {Site.templatePath}/layouts/*.html path.
func (s *Site) loadTemplates() (templatesCompiled *template.Template, err error) {
	funcMap := template.FuncMap{
		"mod": func(a int, b int) int {
			return a % b
		},
		"sub": func(a int, b int) int {
			return a - b
		},
		"time": func() time.Time {
			return time.Now()
		},
		"html": func(a ...interface{}) template.HTML {
			return template.HTML(fmt.Sprint(a...))
		},
	}
	layoutPattern := fmt.Sprintf("%s/layouts/*.html", s.templatePath)
	templatesCompiled, err = template.New("site").Funcs(funcMap).ParseGlob(layoutPattern)
	if err != nil {
		return
	}
	return
}

func (s *Site) loadContent() (root *page.Page, err error) {
	root, err = page.LoadRoot(s.contentPath, s.base)
	if err != nil {
		return
	}
	return
}
