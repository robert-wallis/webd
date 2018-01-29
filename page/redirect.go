// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package page

import "net/url"

func MapRedirects(pages chan *Page, base *url.URL) (m map[string]string) {
	m = make(map[string]string)
	for p := range pages {
		if len(p.Redirects) > 0 {
			for i := range p.Redirects {
				m[p.Redirects[i]] = RelativeBaseOrFullUrl(base, p.URL)
			}
		}
	}
	return
}

// RelativeBaseOrFullUrl returns /page.html if pageURL is example.com/page.html and base is example.com, otherwise it just returns pageUrl
func RelativeBaseOrFullUrl(base *url.URL, pageUrl string) string {
	u, _ := url.Parse(pageUrl)
	if base.Host == u.Host {
		return u.Path
	}
	return pageUrl
}
