// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Redirect(t *testing.T) {
	hostname, _ := url.Parse("http://test")

	r := redirect{Host: "example.com"}

	req := httptest.NewRequest("GET", hostname.String()+"/some-place/?what#yo", nil)
	w := httptest.NewRecorder()

	r.HTTPSRedirect(w, req)

	if w.Code != 301 {
		t.Error("Expecting 301 got", w.Code)
	}

	location := "https://example.com/some-place/?what#yo"
	header := w.Header().Get("Location")
	if header != location {
		t.Error("Expecting", location, "Got", header)
	}

	if req.URL.Host == "example.com" {
		t.Error("Shouldn't modify existing URL expecting: test got: ", req.URL.Host)
	}
}
