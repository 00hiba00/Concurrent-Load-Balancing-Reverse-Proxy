package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	backendURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		log.Fatal("Invalid backend URL:", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
	    req.Header.Set("X-Forwarded-Proto", "http")
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err == nil {
    		req.Header.Set("X-Forwarded-For", ip)
		}
	}
	log.Println("Reverse Proxy starting on :8080...")
	log.Println("Forwarding requests to Backend on :8081")
	err = http.ListenAndServe(":8080", proxy)
	if err != nil {
		log.Fatalf("Proxy server failed: %s", err)
	}
}