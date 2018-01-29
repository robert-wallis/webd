// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// Package to run multiple sites or virtual-hosts on the same computer.
package multisite

import (
	"context"
	"github.com/robert-wallis/webd/config"
	"log"
	"sync"
)

// MultiSite manages multiple different sites.
type MultiSite struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	sites    []*serverSite
}

// New loads a sites.yaml file and creates servers for unique binds internally.
func New(configFilename string, autoCert bool, infoLog, errorLog *log.Logger) (*MultiSite, error) {
	sites, err := config.Load(configFilename)
	if err != nil {
		return nil, err
	}

	httpSites := config.GroupServers(sites)
	m := &MultiSite{
		infoLog:  infoLog,
		errorLog: errorLog,
		sites:    []*serverSite{},
	}
	for bind, list := range httpSites {
		s, err := newServerSite(bind, list, autoCert, infoLog, errorLog)
		if err != nil {
			return nil, err
		}
		m.sites = append(m.sites, s)
	}
	return m, nil
}

// ListenAndServe starts each server in MultiSite, blocks until all inner ListenAndServe return.
func (m *MultiSite) ListenAndServe() error {
	wg := sync.WaitGroup{}
	wg.Add(len(m.sites))
	var err error
	for s := range m.sites {
		site := m.sites[s]
		for r := range site.runningSites {
			m.infoLog.Println("starting", site.runningSites[r].config.Host, "on", site.runningSites[r].bind)
		}
		go func(site *serverSite) {
			err = site.ListenAndServe()
			wg.Done()
		}(m.sites[s])
	}
	wg.Wait()
	return err
}

// Shutdown gracefully shuts down all the servers.
func (m *MultiSite) Shutdown(ctx context.Context) {
	wg := sync.WaitGroup{}
	for s := range m.sites {
		site := m.sites[s]
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := site.Shutdown(ctx); err != nil {
				m.errorLog.Println("Error: shutdown", site.runningSites[0].config.Host, err)
			}
		}()
	}
	wg.Wait()
}
