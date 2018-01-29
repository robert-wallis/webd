// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"github.com/robert-wallis/webd/page"
	"io"
	"net/http"
	"net/url"
)

// ServeHTTP processes requests for the site.  Including dynamic and static content.
func (s *Site) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if s.liveRefresh {
		err := s.loadTemplatesAndContent()
		if err != nil {
			s.errLog.Println("liveRefresh Error:", err)
			// continue to old good version
		}
	}
	if s.redirectHttps && s.base.Scheme == "http" {
		r := redirect{Host: s.base.Host}
		s.infoLog.Println("301 to https", req.Host, req.URL)
		r.HTTPSRedirect(w, req)
		return
	}
	if loc, ok := s.redirectMap[req.URL.Path]; ok {
		s.infoLog.Println("301", req.Host, req.URL)
		http.Redirect(w, req, loc, http.StatusMovedPermanently)
		return
	}
	p, found, folderRedirect := s.contentPage(req.URL.Path)
	if !found {
		s.staticHandler(w, req)
		return
	}
	if folderRedirect {
		u, _ := url.Parse(p.URL)
		s.infoLog.Println("301", req.Host, req.URL, "to", u.Path)
		http.Redirect(w, req, page.RelativeBaseOrFullUrl(s.base, p.URL), http.StatusMovedPermanently)
		return
	}
	if err := s.writePage(w, p); err != nil {
		s.errLog.Println(500, req.Host, req.URL, "Template Execute", err)
		http.Error(w, "Template Execute Error", http.StatusInternalServerError)
		return
	}
	return
}

func (s *Site) writePage(w io.Writer, p *page.Page) (err error) {
	err = s.templates.ExecuteTemplate(w, p.Layout, p)
	return
}
