// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// Site manages templates and content for a site, and even serves the content.
package site

import (
	"github.com/robert-wallis/webd/page"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
)

// Site controls the handling of HTTP traffic to a site.
type Site struct {
	base          *url.URL
	templates     *template.Template
	templatePath  string
	contentPath   string
	staticPath    string
	pageRoot      *page.Page
	pageMap       map[string]*page.Page
	redirectMap   map[string]string
	fileHandler   http.Handler
	liveRefresh   bool
	redirectHttps bool
	infoLog       *log.Logger
	errLog        *log.Logger
}

// New creates and configures a Site.  It loads the templates and content.
// `templatePath` is the place that contains the `layouts` folder.
// `templatePath` contains the `content` folder that is turned into Page objects.
func New(base *url.URL, templatePath string, liveRefresh bool, redirectHttps bool, infoLog, errLog *log.Logger) (s *Site, err error) {
	staticPath := path.Join(templatePath, "static")
	s = &Site{
		base:          base,
		templatePath:  templatePath,
		contentPath:   path.Join(templatePath, "content"),
		staticPath:    staticPath,
		fileHandler:   http.FileServer(http.Dir(staticPath)),
		liveRefresh:   liveRefresh,
		infoLog:       infoLog,
		redirectHttps: redirectHttps,
		errLog:        errLog,
	}
	if err = s.loadTemplatesAndContent(); err != nil {
		return nil, err
	}
	return
}

// contentPage finds the page that matches the url
func (s *Site) contentPage(path string) (page *page.Page, found, folderRedirect bool) {
	if page, found = s.pageMap[path]; !found {
		slashed := path + "/"
		if page, found = s.pageMap[slashed]; found {
			folderRedirect = true
		}
		return
	}
	return
}
