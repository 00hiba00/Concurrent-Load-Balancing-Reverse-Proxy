package balancer

import (
	"sync"
	"sync/atomic"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

type ServerPool struct {
	mux      sync.RWMutex //to protect the backends slice
	Backends []*models.Server
	Current  uint64
	LastServerID  uint64
	Strategy string
}

func (s *ServerPool) GetNextServer() *models.Server {
	if s.Strategy == "round-robin" {
		return RoundRobinStrategy(s)
	} else {
		return LeastConnectionStrategy(s)
	}
}

//treats backend servers like a circular queue,
//handing out requests one by one in a fixed order.
func RoundRobinStrategy(s *ServerPool) *models.Server {
	// 1. Atomically increment the counter and get the result.
	// Even if 1000 goroutines do this at once, they each get a unique 'next'.
	next := atomic.AddUint64(&s.Current, 1)
	s.mux.RLock()
	defer s.mux.RUnlock()

	l := uint64(len(s.Backends))
	if l == 0 {
		return nil
	}

	// 3. Find the next healthy server.
	for i := range l {
		idx := (next + i) % l
		target := s.Backends[idx]

		if target.IsAlive() {
			// Optimization: If we skipped dead servers, we can update
			// 'Current' to 'idx' to help the next request start from a healthy spot.
			if i > 0 {
				atomic.StoreUint64(&s.Current, idx)
			}
			return target
		}
	}

	return nil
}

//looks at all available servers and picks the one
//currently handling the fewest active requests.
func LeastConnectionStrategy(s *ServerPool) *models.Server {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var bestServer *models.Server
	minConnections := -1

	for _, server := range s.Backends {
		if server.IsAlive() {
			conn := server.GetActiveConnections()

			if minConnections == -1 || conn < minConnections {
				minConnections = conn
				bestServer = server
			}
		}
	}

	return bestServer
}

func (s *ServerPool) HealthCheck() {

}

func (s *ServerPool) AddServer() {

}
func (s *ServerPool) RemoveServer() {

}

func (s *ServerPool) GetStatus() {

}
