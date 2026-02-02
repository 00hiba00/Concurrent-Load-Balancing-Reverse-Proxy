package main

import (
	"log"
	"net/http"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/balancer"
)

//Proxy for one backend server
/*
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
}*/

//Load Balancing proxy
func main(){
	log.Println("Load Balancer starting on :8080...")
	log.Printf("Current Strategy: %s", balancer.Pool.Strategy)

	// 2. Use a HandlerFunc to catch every request
	err := http.ListenAndServe(":8080", http.HandlerFunc(proxyHandler))
	if err != nil {
		log.Fatalf("Proxy server failed: %s", err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request){
	targetServer := balancer.Pool.GetNextServer()
	if targetServer == nil{
		http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
		return
	}
	targetServer.IncrementConnections()
	defer targetServer.DecrementConnections()
	targetServer.ReverseProxy.ServeHTTP(w, r)
}