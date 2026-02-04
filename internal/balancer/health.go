package balancer

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

func (s *ServerPool) HealthCheck(ctx context.Context) {
	//time.Ticker is a struct in the time package that:
	//Sends the current time on a channel
	//At regular intervals until you stop it
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
        select {
        case <-ticker.C:
            s.PingServers()
        case <-ctx.Done(): // THIS is the listener
            log.Println("Health Check: Received shutdown signal, stopping...")
            return // The function ends here
        }
    }
}

func (s *ServerPool) PingServers() {
	s.Mux.RLock()
	backends := make([]*models.Server, len(s.Backends))
	copy(backends, s.Backends)
	s.Mux.RUnlock()
	var wg sync.WaitGroup
	for _, srv := range backends {
		wg.Add(1)
		go func(s *models.Server){
			defer wg.Done()
			status := isReachable(srv.URL.Host, 2*time.Second)
			if status != srv.IsAlive() {
				srv.SetAlive(status)
				log.Printf("Health Check: Server [%s] (%s) status changed to: %v", srv.ID, srv.URL.Host, status)
			}
		}(srv)
	}
	wg.Wait()
}

func isReachable(host string, timeout time.Duration) bool{
	//instead of sending a full http request, we are performing
	//a raw TCP Dial
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
