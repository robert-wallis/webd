// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// Package to manage site configuration files.
package config

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Config represents the data of a single site in a sites.yaml file that describes how to configure websites.
type Config struct {
	Host        string   // the main hostname of this site
	Aliases     []string // listed hosts will redirect here
	Email       string   // admin to contact, used for acme
	Static      bool     // true if path points directly to static content, false if it's a dynamic site
	Path        string
	Bind        ConfigBind
	LetsEncrypt bool // use "Let's Encrypt" free auto CA to renew the SSL certificates
}

// ConfigBind is the host and port to bind a TCP socket to.
type ConfigBind struct {
	HTTP  string
	HTTPS string
}

// Load opens the config file at the location in `configFile` and returns all the Config found within that file.
func Load(configFile string) (sites []*Config, err error) {
	var stream *os.File

	// load file
	stream, err = os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("Couldn't Load Config: %v", err)
	}
	defer stream.Close()
	data := &bytes.Buffer{}
	if _, err = data.ReadFrom(stream); err != nil {
		return nil, fmt.Errorf("Error reading %v: %v", configFile, err)
	}
	if err = yaml.Unmarshal(data.Bytes(), &sites); err != nil {
		return nil, fmt.Errorf("Error parsing yaml in %v: %v", configFile, err)
	}

	// fix paths
	dir := filepath.Dir(configFile)
	for c := range sites {
		sites[c].Path = filepath.Join(dir, sites[c].Path)
	}

	return
}

// GroupServers combines configs by bind strings to determine which servers need to start on which ports.
// This method is needed to host multiple hostnames on a single port.
func GroupServers(sites []*Config) (configs map[string][]*Config) {
	configs = make(map[string][]*Config)
	for s := range sites {
		site := sites[s]
		if len(site.Bind.HTTP) > 0 {
			makeAppendSite(&configs, site.Bind.HTTP, site)
		}
		if len(site.Bind.HTTPS) > 0 {
			makeAppendSite(&configs, site.Bind.HTTPS, site)
		}
	}
	return
}

// makeAppendSite creates the first in the list at the `bind` key, or appends to the list at the `bind` key.
func makeAppendSite(m *map[string][]*Config, bind string, site *Config) {
	var l []*Config
	var ok bool
	if l, ok = (*m)[bind]; !ok {
		l = make([]*Config, 1)
		l[0] = site
	} else {
		l = append(l, site)
	}
	(*m)[bind] = l
}

// HostList returns a list of hosts from the site .Host and .Aliases
func (site *Config) HostList() (hosts []string) {
	hosts = append(hosts, site.Host)
	for a := range site.Aliases {
		hosts = append(hosts, site.Aliases[a])
	}
	return
}
