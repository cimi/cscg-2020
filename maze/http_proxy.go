package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// simple http proxy to allow login and initial setup
// we don't need these values at the moment so responses are unchanged
func newHTTPProxy() {
	origin, _ := url.Parse("http://" + serverIP + "/")
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}
	proxy := &httputil.ReverseProxy{Director: director}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
	fmt.Println("Starting HTTP proxy")
	log.Fatal(http.ListenAndServe(":80", nil))
	// https://www.integralist.co.uk/posts/golang-reverse-proxy/#3
	// use Transport to modify responses: https://stackoverflow.com/a/31536962
	// use httputil.DumpRequest(req, true) to debug
}
