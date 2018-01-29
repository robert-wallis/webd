// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package site

import (
	"net/http"
	"net/url"
)

type redirect struct {
	Host string
}

func (r *redirect) HTTPSRedirect(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String()) // make a copy
	u.Host = r.Host
	u.Scheme = "https"
	http.Redirect(w, req, u.String(), http.StatusMovedPermanently)
}
