// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

package config

import (
	"path/filepath"
	"testing"
)

func Test_configSite_LoadConfigs(t *testing.T) {
	// GIVEN a test config
	configFilename := "../test_data/sites.yaml"

	// WHEN the sites are loaded
	sites, err := Load(configFilename)
	if err != nil {
		t.Fatal(err)
	}

	// THEN the expected data should be loaded
	if len(sites) != 2 {
		t.Fatalf("Expecting 2 sites, got %v", len(sites))
	}
	if sites[0].Host != "example.com" {
		t.Errorf("Expecting example.com host got %v", sites[0].Host)
	}
	if sites[0].Aliases[0] != "www.example.com" {
		t.Errorf("Expecting www.example.com Alias got %v", sites[0].Aliases[0])
	}
	if sites[0].Email != "test@example.com" {
		t.Errorf("Expecting test@example.com got %v", sites[0].Email)
	}
	if sites[0].Static != false {
		t.Errorf("Expecting non static, got %v", sites[0].Static)
	}
	if sites[0].Path != filepath.Clean("../example") {
		t.Errorf("Expecting modified path based on file ../example got %v", sites[0].Path)
	}
	if sites[0].LetsEncrypt != true {
		t.Errorf("Expecingt LetsEncrypt got %v", sites[0].LetsEncrypt)
	}
	if sites[0].Bind.HTTP != ":80" {
		t.Errorf("Expecting :80 got %v", sites[0].Bind.HTTP)
	}
	if sites[0].Bind.HTTPS != ":443" {
		t.Errorf("Expecting :443 got %v", sites[0].Bind.HTTPS)
	}
	if sites[1].Static != true {
		t.Errorf("Expecting static got %v", sites[1].Static)
	}
	if sites[1].Path != filepath.Clean("../test_data/files.example.com") {
		t.Errorf("Expecting modified path based on file ../test_data/files.example.com got %v", sites[1].Path)
	}
	if sites[1].Bind.HTTP != ":80" {
		t.Errorf("Expecting :80 got %v", sites[1].Bind.HTTP)
	}
	if sites[1].Bind.HTTPS != "" {
		t.Errorf("Expecting \"\" got \"%v\"", sites[1].Bind.HTTPS)
	}
}

func Test_configSite_LoadConfig_errors(t *testing.T) {
	// GIVEN a nonexistent file, WHEN the sites are loaded THEN it should error
	if _, err := Load("noexist"); err == nil {
		t.Error("Should have errored that the file doesn't exist.")
	}

	// GIVEN a non-yaml file, WHEN the sites are loaded, THEN it should have a parse error
	if _, err := Load("../test_data/not_yaml"); err == nil {
		t.Error("SHould have had a parse error with yaml.")
	}
}

func Test_configSite_CombineServers(t *testing.T) {
	// GIVEN the test list of servers
	sites, err := Load("../test_data/combine_sites.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// WHEN the servers are combined
	configs := GroupServers(sites)

	// THEN they should be grouped as expected
	if len(configs) != 3 {
		t.Log(configs)
		t.Fatalf("Expecting 3 configs binds, got %v", len(configs))
	}
	for bind, list := range configs {
		if bind != "localhost:8101" && bind != "localhost:8202" && bind != "localhost:8443" {
			t.Errorf("Unexpected bind: %v", bind)
		}
		if bind == "localhost:8101" && len(list) != 2 {
			t.Log(list)
			t.Errorf("Expected 2 configs on :8101, got %v", len(list))
		}
		if bind == "localhost:8202" && len(list) != 1 {
			t.Log(list)
			t.Errorf("Expected 1 config on localhost:8202, got %v", len(list))
		}
		if bind == "localhost:8443" && len(list) != 2 {
			t.Log(list)
			t.Errorf("Expected 2 configs on localhost:8443, got %v", len(list))
		}
		for l := range list {
			config := list[l]
			if config.Host != "example.com" &&
				config.Host != "files.example.com" &&
				config.Host != "test.example.com" &&
				config.Host != "secure.example.com" {
				t.Errorf("Unexpected host: %v", config.Host)
			}
		}
	}
}

func Test_configSite_HostList(t *testing.T) {
	// GIVEN a list of test sites
	sites, err := Load("../test_data/sites.yaml")
	if err != nil {
		t.Error(err)
	}

	// WHEN the hosts are enumerated
	var hosts []string
	for s := range sites {
		hs := sites[s].HostList()
		for h := range hs {
			hosts = append(hosts, hs[h])
		}
	}

	// THEN they should contain the .Host and .Aliases fields.
	expected := map[string]bool{
		"example.com":       false,
		"other.example.com": false,
		"www.example.com":   false,
		"files.example.com": false,
	}

	for h := range hosts {
		host := hosts[h]
		if _, ok := expected[host]; !ok {
			t.Errorf("Unexpected host %v", host)
		} else {
			expected[host] = true
		}
	}
	for host, v := range expected {
		if v == false {
			t.Errorf("Host not returned %v", host)
		}
	}
}
