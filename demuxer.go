package main

import (
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
)

type Demuxer struct {
	conf  Config
	proxy *httputil.ReverseProxy
}

func NewDemuxer(conf Config) *Demuxer {
	var d Demuxer
	d.conf = conf
	d.proxy = &httputil.ReverseProxy{Director: func(req *http.Request) {
		host := d.getHostname(req.Host)
		req.URL.Scheme = "http"
		req.URL.Host = d.conf.Forwards[host]
	}}

	return &d
}

func (d Demuxer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	host := d.getHostname(req.Host)
	forward := d.conf.Forwards[host]
	redirect := d.conf.Redirects[host]

	if forward != "" {
		d.proxy.ServeHTTP(rw, req)
	} else if redirect != "" {
		http.Redirect(rw, req, redirect, http.StatusMovedPermanently)
	} else {
		http.NotFound(rw, req)
	}
}

func (d Demuxer) getHostname(host string) string {
	if m, _ := regexp.MatchString(":[0-9]+$", host); m {
		host, _, _ = net.SplitHostPort(host)
	}
	return host
}