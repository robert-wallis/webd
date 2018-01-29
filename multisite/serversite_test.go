// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package multisite

import (
	"bytes"
	"context"
	"fmt"
	"github.com/robert-wallis/webd/config"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	server
	testServer *httptest.Server
}

func newTestServer(handler http.Handler) *testServer {
	return &testServer{
		testServer: httptest.NewServer(handler),
	}
}

func (s *testServer) ListenAndServe() error {
	return nil
}

func (s *testServer) Shutdown(ctx context.Context) (err error) {
	s.testServer.Close()
	if ctx != nil {
		err = ctx.Err()
	}
	return
}

func Test_MultiSite_serverSite_ServeHTTP(t *testing.T) {
	// GIVEN a successfully loaded multi-site
	testLog := log.New(&bytes.Buffer{}, "", 0)
	ms, err := New("../test_data/combine_sites.yaml", false, testLog, testLog)
	if err != nil {
		t.Fatal(err)
	}

	site := map[string]*serverSite{}
	for s := range ms.sites {
		siteName := fmt.Sprintf("%v:%v", ms.sites[s].runningSites[0].config.Host, ms.sites[s].runningSites[0].bind)
		site[siteName] = ms.sites[s]
	}

	// WHEN a request comes in for a specific site
	type requestTest struct {
		host     string
		path     string
		code     int
		siteName string
	}
	tests := []requestTest{
		{"example.com", "/example.com.txt", 301, "example.com:localhost:8101"},
		{"example.com", "/files.example.com.txt", 301, "example.com:localhost:8101"},
		{"example.com", "/example.com.txt", 200, "example.com:localhost:8443"},
		{"example.com", "/files.example.com.txt", 404, "example.com:localhost:8443"},
		{"files.example.com", "/files.example.com.txt", 200, "example.com:localhost:8101"},
		{"files.example.com", "/example.com.txt", 404, "example.com:localhost:8101"},
		{"test.example.com", "/test.example.com.txt", 200, "test.example.com:localhost:8202"},
		{"test.example.com", "/example.com.txt", 404, "test.example.com:localhost:8202"},
		{"example.com", "/example.com.txt", 200, "example.com:localhost:8443"},
		{"example.com", "/files.example.com.txt", 404, "example.com:localhost:8443"},
		{"not-configured.example.com", "/", 502, "example.com:localhost:8101"},
		{"not-configured.example.com", "/", 502, "example.com:localhost:8443"},
		{"not-configured.example.com", "/", 502, "test.example.com:localhost:8202"},
		{"example.com:80", "/example.com.txt", 301, "example.com:localhost:8101"},
		{"example.com:443", "/example.com.txt", 200, "example.com:localhost:8443"},
	}
	for i := range tests {
		// THEN the response should be from the correct site folder
		test := tests[i]
		req := httptest.NewRequest("GET", fmt.Sprintf("https://%s%s", test.host, test.path), nil)
		w := httptest.NewRecorder()
		site[test.siteName].ServeHTTP(w, req)
		if w.Code != test.code {
			t.Errorf("Expecting response %v got %v for #%v %v", test.code, w.Code, i, test)
		}
	}
}

func Test_serverSite_initTLS(t *testing.T) {
	// GIVEN a partially configured serverSite
	hs := &http.Server{}
	ss := &serverSite{
		runningSites: []*runningSite{{
			config: &config.Config{
				Host:    "test-host",
				Aliases: []string{"host-alias"},
			},
		}, {
			config: &config.Config{
				Host:  "site2",
				Email: "site2@example.com",
			},
		},
		},
		server: hs,
	}

	// WHEN the serverSite is configured for TLS
	ss.initTLS(":443", "localhost", hs, true)

	// THEN the TLSConfig should be setup
	if ss.tlsEnabled == false {
		t.Fatal("Expecting tlsEnabled to be true")
	}

	if hs.TLSConfig == nil {
		t.Fatal("TLSConfig was nil")
	}

	if hs.TLSConfig.GetCertificate == nil {
		t.Error("TLSCOnfig.GetCertificate function was not set.")
	}
}

func Test_stripPort(t *testing.T) {
	type test struct {
		host     string
		expected string
	}
	tests := []test{
		{"host:23", "host"},
		{":23", ""},
		{"host:", "host"},
		{"host", "host"},
		{":", ""},
		{"", ""},
	}
	for i := range tests {
		tst := tests[i]
		result := stripPort(tst.host)
		if result != tst.expected {
			t.Errorf(`Host "%v" expecting "%v" got "%v"`, tst.host, tst.expected, result)
		}
	}
}

func Test_justPort(t *testing.T) {
	type test struct {
		host     string
		expected string
	}
	tests := []test{
		{"host:23", "23"},
		{":23", "23"},
		{"host:", ""},
		{"host", ""},
		{":", ""},
		{"", ""},
	}
	for i := range tests {
		tst := tests[i]
		result := justPort(tst.host)
		if result != tst.expected {
			t.Errorf(`Host "%v" expecting "%v" got "%v"`, tst.host, tst.expected, result)
		}
	}
}

func Test_hostList(t *testing.T) {
	// GIVEN a partially configured serverSite
	hs := &http.Server{}
	ss := &serverSite{
		runningSites: []*runningSite{{
			config: &config.Config{
				Host:    "test-host",
				Aliases: []string{"host-alias"},
			},
		}, {
			config: &config.Config{
				Host:  "site2",
				Email: "site2@example.com",
			},
		},
		},
		server: hs,
	}

	// WHEN this host list is generated
	hosts := hostList(ss.runningSites)

	// THEN it should match the test
	if len(hosts) != 3 {
		t.Errorf(`Host len 3 expected, got %d`, len(hosts))
	}
	tests := []string{"test-host", "host-alias", "site2"}
	for i := range tests {
		testHost := tests[i]
		found := false
		for hl := range hosts {
			var listHost = hosts[hl]
			if testHost == listHost {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Couldn't find %s in list %v", testHost, hostList)
		}
	}
}
