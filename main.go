// Copyright (C) 2018 Robert A. Wallis, All Rights Reserved.

// webd multiple site web server.
// Serves static pages, or pages built with html templates and yaml content.
// Uses Let's Encrypt, and automatically renews HTTPS certificates.
package main

import (
	"flag"
	"fmt"
	"github.com/robert-wallis/webd/multisite"
	"github.com/robert-wallis/webd/site"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const VERSION = "2019-01-28"

var _bind = flag.String("bind", ":80", "bind ip and port")
var _hostname = flag.String("hostname", "example.com", "outside hostname for site")
var _liveRefresh = flag.Bool("live-refresh", false, "Should reload all templates each request?")
var _autoCert = flag.Bool("auto-cert", true, "Automatically get and renew TLS/SSL certificates?")

const (
	ExitSingleSiteInit = iota
	ExitSingleParam
	ExitSingleSiteRuntime
	ExitMultiSiteInit
	ExitMultiSiteRuntime
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\nVersion %s\n\n", os.Args[0], VERSION)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	infoLog := log.New(os.Stdout, "", log.LstdFlags)
	errorLog := log.New(os.Stderr, "", log.LstdFlags)
	basePath := filepath.Base(os.Args[0])

	if flag.NArg() == 0 {
		infoLog.Println("Starting", basePath, VERSION, "Single Site Mode")
		singleSite(infoLog, errorLog)
	} else {
		infoLog.Println("Starting", basePath, VERSION, "Multiple Site Mode")
		multiSite(flag.Arg(0), infoLog, errorLog)
	}
}

func singleSite(infoLog, errorLog *log.Logger) {
	base, err := url.Parse(*_hostname)
	if err != nil {
		errorLog.Println(err)
		os.Exit(ExitSingleParam)
		return
	}
	s, err := site.New(base, ".", *_liveRefresh, false, infoLog, errorLog)
	if err != nil {
		errorLog.Println(err)
		os.Exit(ExitSingleSiteInit)
		return
	}
	if err = http.ListenAndServe(*_bind, s); err != nil {
		errorLog.Printf("Server Error: %v\n", err)
		os.Exit(ExitSingleSiteRuntime)
	}
}

func multiSite(siteConfigFile string, infoLog, errorLog *log.Logger) {
	ms, err := multisite.New(siteConfigFile, *_autoCert, infoLog, errorLog)
	if err != nil {
		errorLog.Println(err)
		os.Exit(ExitMultiSiteInit)
	}
	if err = ms.ListenAndServe(); err != nil {
		errorLog.Println(err)
		os.Exit(ExitMultiSiteRuntime)
	}
}
