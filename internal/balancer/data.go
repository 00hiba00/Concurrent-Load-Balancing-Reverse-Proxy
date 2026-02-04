package balancer

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

func NewServer(id string, rawURL string) *models.Server {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
        rawURL = "http://" + rawURL
    }
	backendURL, err := url.Parse(rawURL)
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

	return &models.Server{
		ID:           id,
		URL:          backendURL,
		Alive:        true,
		ReverseProxy: proxy,
	}
}


var Pool = &ServerPool{
	Backends: []*models.Server{},
	Strategy: "round-robin",
}
