// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/robert-wallis/webd/page"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_Site_staticHandler(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		path string
		code int
	}

	tests := []test{
		{"/favicon.ico", 200},
		{"/noexist", 404},
	}

	for i := range tests {
		req := httptest.NewRequest("GET", address.String()+tests[i].path, nil)
		w := httptest.NewRecorder()
		s.staticHandler(w, req)
		if w.Code != tests[i].code {
			t.Error(tests[i].path, "expected", tests[i].code, "actual", w.Code)
		}
	}
}

func Test_Site_notFoundHandler(t *testing.T) {
	// GIVEN the page path that has no 404 page
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}
	// no 404 page in page content
	s.contentPath = "../page/test_content"
	if s.pageRoot, err = s.loadContent(); err != nil {
		t.Fatal(err)
	}
	s.pageMap = page.MapPages(s.pageRoot)

	// WHEN the request is given to a page that doesn't exist
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/noexist", address), nil)
	w := httptest.NewRecorder()
	s.notFoundHandler(w, req)

	// THEN it should continue to 404
	if w.Code != 404 {
		t.Error("expected 404 actual", w.Code)
	}
}

func Test_Site_notFoundHandler_closed_writer(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	errBuf := &bytes.Buffer{}
	errLog := log.New(errBuf, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, errLog)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", address.String()+"/", nil)
	testError := "test static handler closed"
	rw := mockResponseWriter{
		HeaderReturn:     http.Header{},
		WriteReturnInt:   0,
		WriteReturnError: errors.New(testError),
	}
	s.notFoundHandler(rw, req)

	errStr := errBuf.String()
	if !strings.Contains(errStr, testError) {
		t.Error("Expected", testError, "Got", errStr)
	}
}
