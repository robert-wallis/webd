// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type mockResponseWriter struct {
	http.ResponseWriter
	HeaderReturn      http.Header
	WriteInBytes      []byte
	WriteReturnInt    int
	WriteReturnError  error
	WriteHeaderInCode int
}

func (m mockResponseWriter) Header() http.Header {
	return m.HeaderReturn
}

func (m mockResponseWriter) Write(bytes []byte) (int, error) {
	m.WriteInBytes = bytes
	return m.WriteReturnInt, m.WriteReturnError
}

func (m mockResponseWriter) WriteHeader(code int) {
	m.WriteHeaderInCode = code
}

func Test_Site_ServeHTTP(t *testing.T) {
	// GIVEN an initialized site instance
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN a request comes in
	req := httptest.NewRequest("GET", address.String()+"/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	// THEN the response should be successful
	if w.Code != 200 {
		t.Error("Index handler code expected 200 actual", w.Code)
	}
}

func Test_Site_ServeHTTP_liveRefresh_error(t *testing.T) {
	// GIVEN an initialized site instance with liveRefresh enabled
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	errBuf := &bytes.Buffer{}
	errLog := log.New(errBuf, "", 0)
	s, err := New(address, _templatePath, true, false, testLog, errLog)
	if err != nil {
		t.Fatal(err)
	}

	// WHEN a the templates are broken and a new request comes in
	s.templatePath = "noexist"
	req := httptest.NewRequest("GET", address.String()+"/", nil)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, req)

	// THEN the response should be successful
	if w.Code != 200 {
		t.Error("Index handler code expected 200 actual", w.Code)
	}

	// THEN the error log should contain an error about the templates
	errStr := errBuf.String()
	if !strings.Contains(errStr, "pattern matches no files") {
		t.Error(errStr)
	}
}

func Test_Site_ServeHTTP_sad_path(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	errBuf := &bytes.Buffer{}
	errLog := log.New(errBuf, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, errLog)
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		path     string
		code     int
		location string
	}

	tests := []test{
		{"/privacy.html", 301, "/privacy/"},
		{"/privacy", 301, "/privacy/"},
		{"/noexist", 404, ""},
	}

	for i := range tests {
		req := httptest.NewRequest("GET", address.String()+tests[i].path, nil)
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)
		if w.Code != tests[i].code {
			t.Error(tests[i].path, "expected", tests[i].code, "actual", w.Code, errBuf.String())
		}
		loc := w.HeaderMap.Get("location")
		if w.Code >= 300 && w.Code < 400 && loc != tests[i].location {
			t.Error(tests[i].path, "expected", tests[i].location, "redirect location:", loc, w.Code)
		}
	}
}

func Test_Site_ServeHTTP_https_redirect(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	errBuf := &bytes.Buffer{}
	errLog := log.New(errBuf, "", 0)
	s, err := New(address, _templatePath, false, true, testLog, errLog)
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		path     string
		code     int
		location string
	}

	// should just reflect request without thinking about it to https
	tests := []test{
		{"/", 301, "https://localhost:8009/"},
		{"/privacy.html", 301, "https://localhost:8009/privacy.html"},
		{"/privacy", 301, "https://localhost:8009/privacy"},
		{"/noexist", 301, "https://localhost:8009/noexist"},
	}

	for i := range tests {
		req := httptest.NewRequest("GET", address.String()+tests[i].path, nil)
		w := httptest.NewRecorder()
		s.ServeHTTP(w, req)
		if w.Code != tests[i].code {
			t.Error(tests[i].path, "expected", tests[i].code, "actual", w.Code, errBuf.String())
		}
		loc := w.HeaderMap.Get("location")
		if w.Code >= 300 && w.Code < 400 && loc != tests[i].location {
			t.Error(tests[i].path, "expected", tests[i].location, "redirect location:", loc, w.Code)
		}
	}
}

func Test_Site_ServeHTTP_closed_writer(t *testing.T) {
	address, _ := url.Parse("http://localhost:8009")
	testLog := log.New(&bytes.Buffer{}, "", 0)
	errBuf := &bytes.Buffer{}
	errLog := log.New(errBuf, "", 0)
	s, err := New(address, _templatePath, false, false, testLog, errLog)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", address.String()+"/", nil)
	testError := "test handler closed"
	rw := mockResponseWriter{
		HeaderReturn:     http.Header{},
		WriteReturnInt:   0,
		WriteReturnError: errors.New(testError),
	}
	s.ServeHTTP(rw, req)

	errStr := errBuf.String()
	if !strings.Contains(errStr, testError) {
		t.Errorf(`expected "%v" got "%v"`, testError, errStr)
	}
}
