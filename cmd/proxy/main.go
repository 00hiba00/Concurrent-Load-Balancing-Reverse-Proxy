package main

import (
	"context"
	"log"
	"net/http"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/admin"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/balancer"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
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
	var targetServer *models.Server
	cookie, err := r.Cookie("LB_STICKY_SESSION")
	for i := range 3{
		if i==0 && err == nil{
			targetServer = balancer.Pool.GetServerByID(cookie.Value)
			if targetServer == nil || !targetServer.IsAlive() {
				targetServer = nil
		}
		}
		if targetServer == nil{
			targetServer = balancer.Pool.GetNextServer()
		}
		if targetServer == nil{
			http.Error(w, "Service Unavailable: No healthy backends", http.StatusServiceUnavailable)
			return
		}
		sucessCnx := true
		targetServer.IncrementConnections()
		//if a backend dies before the health check sees it,
		//we want to mark it as dead immediately and try another one.
		targetServer.ReverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			targetServer.SetAlive(false)
			sucessCnx = false
		}
		http.SetCookie(w, &http.Cookie{
				Name:     "LB_STICKY_SESSION",
				Value:    targetServer.ID,
				Path:     "/",
				HttpOnly: true,  // Security : make it invisible to frontend code
				MaxAge:   1800,  // Each time reset 30 minutes
    	})
		targetServer.ReverseProxy.ServeHTTP(w, r)
		targetServer.DecrementConnections()
		if sucessCnx {
			return
		}
	}
	http.Error(w, "All backends are currently unreachable", http.StatusServiceUnavailable)
}

func RunAdminServer(pool *balancer.ServerPool) {
    adminMux := http.NewServeMux()

    adminMux.HandleFunc("GET /dashboard", admin.DashboardHandler(pool))

	adminMux.HandleFunc("GET /backends", admin.GetBackendsHandler(pool))
	adminMux.HandleFunc("GET /backends/{id}", admin.GetServerHandler(pool))
	adminMux.HandleFunc("POST /backends", admin.PostServerHandler(pool))
	adminMux.HandleFunc("PATCH /backends/{id}", admin.UpdateServerHandler(pool))
	adminMux.HandleFunc("DELETE /backends/{id}", admin.DeleteServerHandler(pool))

	log.Println("Admin Dashboard available at http://localhost:9090/dashboard")
    log.Println("Admin API is running on :9090")

    if err := http.ListenAndServe(":9090", adminMux); err != nil {
        log.Fatalf("Admin server failed: %v", err)
    }
}