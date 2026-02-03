package main

import (
	"context"
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
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
	go balancer.Pool.HealthCheck(ctx)

	go RunAdminServer(balancer.Pool)
	log.Println("Load Balancer starting on :8080...")
	log.Printf("Current Strategy: %s", balancer.Pool.Strategy)

	// 2. Use a HandlerFunc to catch every request
	err := http.ListenAndServe(":8080", http.HandlerFunc(proxyHandler))
	if err != nil {
		log.Fatalf("Proxy server failed: %s", err)
	}
}

func proxyHandler(w http.ResponseWriter, r *http.Request){
	for range 3{
		targetServer := balancer.Pool.GetNextServer()
		if targetServer == nil{
			http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
			break
		}
		sucessCnx := true
		targetServer.IncrementConnections()
		defer targetServer.DecrementConnections()
		targetServer.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			targetServer.SetAlive(false)
			sucessCnx = false
		}
		targetServer.ReverseProxy.ServeHTTP(w, r)
		if sucessCnx {
			return
		}
	}
	http.Error(w, "All backends are currently unreachable", http.StatusServiceUnavailable)
}

func RunAdminServer(pool *balancer.ServerPool) {
    adminMux := http.NewServeMux()

    adminMux.HandleFunc("GET /backends", pool.GetBackendsHandler)
    adminMux.HandleFunc("GET /backends/{id}", pool.GetServerHandler)
    adminMux.HandleFunc("POST /backends", pool.PostServerHandler)
    adminMux.HandleFunc("PATCH /backends/{id}", pool.UpdateServerHandler)
    adminMux.HandleFunc("DELETE /backends/{id}", pool.DeleteServerHandler)

    log.Println("Admin API is running on :9090")

    if err := http.ListenAndServe(":9090", adminMux); err != nil {
        log.Fatalf("Admin server failed: %v", err)
    }
}