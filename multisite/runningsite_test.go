// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package multisite

import (
	"github.com/robert-wallis/webd/config"
	"testing"
)

func Test_baseUrl(t *testing.T) {
	cfg := &config.Config{
		Host: "example.com",
		Bind: config.ConfigBind{
			HTTP:  ":80",
			HTTPS: ":443",
		},
	}
	type test struct {
		bind     string
		expected string
	}
	tests := []test{
		{":80", "http://example.com"},
		{":443", "https://example.com"},
	}
	for i := range tests {
		tst := tests[i]
		got, err := baseUrl(cfg, tst.bind)
		if err != nil {
			t.Error(err)
		}
		if got.String() != tst.expected {
			t.Errorf(`bind "%v" expected "%v" got "%v"`, tst.bind, tst.expected, got)
		}
	}
}

func Test_HostList(t *testing.T) {
	// GIVEN a partially configured serverSite
	runningSites := []*runningSite{
		{config: &config.Config{Host: "test-host", Aliases: []string{"host-alias"}}},
		{config: &config.Config{Host: "site2", Email: "site2@example.com"}},
	}

	// WHEN this host list is generated
	hosts := hostList(runningSites)

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

func Test_ShouldRedirectHttps(t *testing.T) {
	type test struct {
		Bind  config.ConfigBind
		Https bool
	}

	tests := []test{
		{config.ConfigBind{HTTP: "", HTTPS: ""}, false},
		{config.ConfigBind{HTTP: ":80", HTTPS: ""}, false},
		{config.ConfigBind{HTTP: ":8080", HTTPS: ""}, false},
		{config.ConfigBind{HTTP: ":443", HTTPS: ""}, false},
		{config.ConfigBind{HTTP: "", HTTPS: ":443"}, true},
		{config.ConfigBind{HTTP: ":80", HTTPS: ":443"}, true},
		{config.ConfigBind{HTTP: "", HTTPS: ":80"}, true},
		{config.ConfigBind{HTTP: "", HTTPS: ":8443"}, true},
		{config.ConfigBind{HTTP: "", HTTPS: ":8080"}, true},
	}

	for i := range tests {
		tst := tests[i]
		cfg := &config.Config{
			Bind: tst.Bind,
		}
		result := shouldRedirectHttps(cfg)
		if result != tst.Https {
			t.Errorf("Expecting %v got %v for Http:%v Https:%v", tst.Https, result, tst.Bind.HTTP, tst.Bind.HTTPS)
		}
	}
}
