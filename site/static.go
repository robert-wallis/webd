// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"net/http"
	"os"
	"path"
)

func (s *Site) staticHandler(w http.ResponseWriter, req *http.Request) {
	filePath := path.Join(s.staticPath, req.URL.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.notFoundHandler(w, req)
		return
	}
	s.fileHandler.ServeHTTP(w, req)
}

func (s *Site) notFoundHandler(w http.ResponseWriter, req *http.Request) {
	s.infoLog.Println(404, req.Host, req.URL)
	p, found, _ := s.contentPage("/404/")
	if !found {
		s.errLog.Println(404, req.Host, req.URL, "Error: 404.yaml template not found")
		http.Error(w, "Resource Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(404)
	if err := s.writePage(w, p); err != nil {
		s.errLog.Println(500, req.Host, req.URL, "Template Execute", err)
		http.Error(w, "Template Execute Error", http.StatusInternalServerError)
		return
	}
}
