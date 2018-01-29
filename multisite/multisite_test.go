// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package multisite

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test_NewMultiSite(t *testing.T) {
	// GIVEN a valid config
	configFilename := "../test_data/combine_sites.yaml"
	testLog := log.New(&bytes.Buffer{}, "", 0)

	// WHEN a MultiSite is created
	ms, err := New(configFilename, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	var prodSite, testSite *serverSite
	for s := range ms.sites {
		if _, ok := ms.sites[s].hostMap["test.example.com"]; ok {
			testSite = ms.sites[s]
		} else {
			prodSite = ms.sites[s]
		}
	}

	// THEN the right number of sub objects should be made
	if len(prodSite.hostMap) != 3 {
		t.Errorf("Expecting 3 hosts got %v", len(ms.sites[0].hostMap))
	}
	if len(testSite.hostMap) != 1 {
		t.Errorf("Expecting 1 host got %v", len(ms.sites[1].hostMap))
	}
	if len(ms.sites) != 3 {
		t.Errorf("Expecting 3 running sites got %v", len(ms.sites))
	}

	// THEN the host map should point to the right objects
	listedSites := make(map[*serverSite]bool)
	for s := range ms.sites {
		for host, site := range ms.sites[s].hostMap {
			listedSites[ms.sites[s]] = true
			if site.config.Host == host {
				continue
			}
			var found bool
			for a := range site.config.Aliases {
				if site.config.Aliases[a] == host {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Host %v is incorrectly mapped to site %v", host, site)
			}
		}
	}

	// THEN the list of sites should all be mapped
	for s := range ms.sites {
		site := ms.sites[s]
		if _, ok := listedSites[site]; !ok {
			t.Errorf("Site %v not in hostMap", site)
		}
	}

	// THEN the logs should be the ones passed in
	if ms.infoLog != testLog || ms.errorLog != testLog {
		t.Errorf("Expecting %x for logs, they were %x and %x", testLog, ms.infoLog, ms.errorLog)
	}
}

func Test_MultiSite_https_redirect(t *testing.T) {

	// GIVEN a valid config
	configFilename := "../test_data/combine_sites.yaml"
	testLog := log.New(&bytes.Buffer{}, "", 0)

	// WHEN a MultiSite is created
	ms, err := New(configFilename, false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}
	prodSite, _, _ := mapSites(ms)

	// THEN the http://example.com site should redirect to the https version.
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "example.com:8101"
	rr := httptest.NewRecorder()
	prodSite.ServeHTTP(rr, req)

	if rr.Code != 301 {
		t.Errorf("Expecting 301 got %v", rr.Code)
	}
	location := rr.HeaderMap.Get("Location")
	if location != "https://example.com/" {
		t.Errorf("Expecting https, got %v", location)
	}
}

func Test_NewMultiSite_error(t *testing.T) {
	// GIVEN a config file that doesn't exist
	configFilename := "noexist"
	testLog := log.New(&bytes.Buffer{}, "", 0)

	// WHEN a MultiSite is created
	ms, err := New(configFilename, false, testLog, testLog)

	// THEN it should error
	if err == nil {
		t.Error("Expected error that the config file doesn't exist.")
	}
	if err != nil && ms != nil {
		t.Error("When there's an error, it shouldn't return a MultiSite")
	}
}

func Test_MultiSite_Serve(t *testing.T) {
	// GIVEN a successfully loaded multi-site
	testLog := log.New(&bytes.Buffer{}, "", 0)
	ms, err := New("../test_data/combine_sites.yaml", false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	prodUrl, testUrl, secureUrl := injectTestServers(ms)

	// WHEN the MultiServer starts ListenAndServe
	go ms.ListenAndServe()

	// THEN all the sites should be up and running
	type requestTest struct {
		url  string
		host string
		code int
	}
	tests := []requestTest{
		{prodUrl + "/files.example.com.txt", "files.example.com", 200},
		{testUrl + "/test.example.com.txt", "test.example.com", 200},
		{secureUrl + "/example.com.txt", "example.com", 200},
		{secureUrl + "/secure.example.com.txt", "secure.example.com", 200},
		{prodUrl + "/nope", "example.com", 404},
		{prodUrl + "/nope", "files.example.com", 404},
		{testUrl + "/nope", "test.example.com", 404},
		{secureUrl + "/nope", "example.com", 404},
		{prodUrl + "/", "not-configured.example.com", 502},
		{prodUrl + "/", "", 502},
		{testUrl + "/", "example.com", 502},
		{testUrl + "/", "files.example.com", 502},
		{prodUrl + "/", "test.example.com", 502},
		{secureUrl + "/", "files.example.com", 502},
	}
	for i := range tests {
		test := tests[i]
		req, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Error(err)
			continue
		}
		req.Host = test.host
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			continue
		}
		res.Body.Close()
		if res.StatusCode != test.code {
			t.Errorf("Expecting response %v got %v for %v", test.code, res.StatusCode, test)
		}
	}
	ms.Shutdown(nil)
}

func Test_MultiSite_Shutdown_error(t *testing.T) {
	// GIVEN a successfully loaded multi-site
	infoLog := log.New(&bytes.Buffer{}, "", 0)
	errorBuf := &bytes.Buffer{}
	errorLog := log.New(errorBuf, "", 0)
	ms, err := New("../test_data/combine_sites.yaml", false, infoLog, errorLog)
	if err != nil {
		t.Fatal(err)
	}
	prodUrl, _, _ := injectTestServers(ms)
	go ms.ListenAndServe()
	if errorBuf.Len() > 0 {
		t.Fatal(errorBuf)
	}

	// WHEN it is shut down with an error
	u, _ := url.Parse(prodUrl)
	hangingClient, err := net.Dial("tcp", u.Host)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now())
	ms.Shutdown(ctx)
	hangingClient.Close()
	cancel()

	// THEN the logs should have an error
	errStr := errorBuf.String()
	if !strings.Contains(errStr, "deadline exceeded") {
		t.Errorf(`Expected error "deadline exceeded" not found in error log: "%v"`, errStr)
	}
}

func injectTestServers(ms *MultiSite) (prodUrl, testUrl, secureUrl string) {
	sites := map[string]string{}
	for s := range ms.sites {
		site := ms.sites[s]
		ts := newTestServer(site)
		site.server = ts
		siteName := fmt.Sprintf("%v:%v", ms.sites[s].runningSites[0].config.Host, ms.sites[s].runningSites[0].bind)
		sites[siteName] = ts.testServer.URL
	}
	return sites["example.com:localhost:8101"], sites["test.example.com:localhost:8202"], sites["example.com:localhost:8443"]
}

func mapSites(ms *MultiSite) (prod, test, secure *serverSite) {
	sites := map[string]*serverSite{}
	for s := range ms.sites {
		siteName := fmt.Sprintf("%v:%v", ms.sites[s].runningSites[0].config.Host, ms.sites[s].runningSites[0].bind)
		sites[siteName] = ms.sites[s]
	}
	return sites["example.com:localhost:8101"], sites["test.example.com:localhost:8202"], sites["example.com:localhost:8443"]
}
