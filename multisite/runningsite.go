// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package multisite

import (
	"github.com/robert-wallis/webd/config"
	"github.com/robert-wallis/webd/site"
	"log"
	"net/http"
	"net/url"
)

// runningSite maps a Config to a server.
type runningSite struct {
	config     *config.Config
	serverSite *serverSite
	handler    http.Handler
	site       *site.Site
	bind       string
}

func newRunningSite(serverSite *serverSite, config *config.Config, bind string, infoLog, errorLog *log.Logger) (r *runningSite, err error) {
	r = &runningSite{
		config:     config,
		serverSite: serverSite,
		bind:       bind,
	}
	switch {
	case config.Static:
		r.handler = http.FileServer(http.Dir(config.Path))
	default:
		base, err := baseUrl(config, bind)
		if err != nil {
			return nil, err
		}
		if r.site, err = site.New(base, config.Path, false, shouldRedirectHttps(config), infoLog, errorLog); err != nil {
			return nil, err
		}
		r.handler = r.site
	}
	return
}

func (r *runningSite) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r *runningSite) HostList() (hosts []string) {
	hl := r.config.HostList()
	for h := range hl {
		hosts = append(hosts, hl[h])
	}
	return
}

// baseUrl returns the base url for the configuration to be used as a rel-canonical link or HTML5 base.
func baseUrl(config *config.Config, bind string) (u *url.URL, err error) {
	proto := "http"
	if len(bind) > 0 && config.Bind.HTTPS == bind {
		proto = "https"
	}
	u, err = url.Parse(proto + "://" + config.Host)
	return
}

func shouldRedirectHttps(config *config.Config) bool {
	return len(config.Bind.HTTPS) > 0
}
