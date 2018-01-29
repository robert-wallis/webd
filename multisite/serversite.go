// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package multisite

import (
	"context"
	"crypto/tls"
	"github.com/robert-wallis/webd/config"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"strings"
)

type serverSite struct {
	errorLog     *log.Logger
	server       server
	runningSites []*runningSite
	hostMap      map[string]*runningSite
	tlsEnabled   bool
}

// server is something that can ListenAndServe and Shutdown.
type server interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
}

// newServerSite creates and initializes an http.Server to go with a list of configs.
func newServerSite(bind string, configs []*config.Config, autoCert bool, infoLog, errorLog *log.Logger) (*serverSite, error) {
	s := &serverSite{
		errorLog:     errorLog,
		runningSites: []*runningSite{},
		hostMap:      make(map[string]*runningSite),
	}
	hs := &http.Server{
		Addr:     bind,
		Handler:  s,
		ErrorLog: errorLog,
	}
	s.server = hs
	for c := range configs {
		cfg := configs[c]
		r, err := newRunningSite(s, cfg, bind, infoLog, errorLog)
		if err != nil {
			return nil, err
		}
		s.appendHostMap(r)
		s.runningSites = append(s.runningSites, r)
	}
	s.initTLS(bind, configs[0].Host, hs, autoCert)
	return s, nil
}

// ListenAndServe starts the underlying http server.
func (s *serverSite) ListenAndServe() error {
	if s.tlsEnabled {
		return s.server.ListenAndServeTLS("", "")
	}
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *serverSite) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Handler takes an incoming request and sends it off to the correct site within MultiSite.
func (s *serverSite) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r, ok := s.hostMap[stripPort(req.Host)]
	if !ok {
		s.errorLog.Println(http.StatusBadGateway, req.Host, req.URL, req.Header)
		http.Error(w, "Site Not Configured", http.StatusBadGateway)
		return
	}
	r.ServeHTTP(w, req)
}

// appendHostMap adds the host names of the site to a map of sites
func (s *serverSite) appendHostMap(site *runningSite) {
	hosts := site.config.HostList()
	for h := range hosts {
		s.hostMap[hosts[h]] = site
	}
}

// stripPort is stolen from url.stripPort, because http.Request.URL doesn't have the hostname
func stripPort(host string) string {
	colon := strings.IndexByte(host, ':')
	if colon == -1 {
		return host
	}
	return host[:colon]
}

// justPort returns only the port portion of the host originally from http.Request.Host
func justPort(host string) string {
	colon := strings.IndexByte(host, ':')
	if colon == -1 {
		return ""
	}
	return host[colon+len(":"):]
}

// initTLS sets up TLS if the bind port is 443
func (s *serverSite) initTLS(bind, host string, hs *http.Server, autoCert bool) {
	if justPort(bind) != "443" {
		return
	}
	hs.TLSConfig = &tls.Config{
		ServerName: host,
	}
	if autoCert {
		s.initAutoCert(hs.TLSConfig)
	}
	s.tlsEnabled = true
}

// initAutoCert sets the GetCertificate function to get the cert automatically from the CA (Let's Encrypt)
func (s *serverSite) initAutoCert(tlsConfig *tls.Config) {
	hosts := hostList(s.runningSites)
	acManager := autocert.Manager{
		Email:      firstEmailFound(s.runningSites),
		Cache:      autocert.DirCache("autocert"),
		Prompt:     acme.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hosts...),
	}
	tlsConfig.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		return acManager.GetCertificate(info)
	}
}

// hostList enumerates all the hosts in the list of sites.
func hostList(sites []*runningSite) (hosts []string) {
	for r := range sites {
		hl := sites[r].HostList()
		for h := range hl {
			hosts = append(hosts, hl[h])
		}
	}
	return
}

// firstEmailFound checks each site in the list for a configured Email, and returns the first one found.
func firstEmailFound(sites []*runningSite) (email string) {
	for r := range sites {
		if sites[r].config.Email != "" {
			return sites[r].config.Email
		}
	}
	return ""
}
