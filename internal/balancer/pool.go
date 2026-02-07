package balancer

import (
	"sync"
	"sync/atomic"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/models"
)

type ServerPool struct {
	Mux      sync.RWMutex //to protect the backends slice
	Backends []*models.Server
	Current  uint64
	LastServerID  uint64
	Strategy string
}

func (s *ServerPool) GetNextServer() *models.Server {
	if s.Strategy == "round-robin" {
		return WeightedRoundRobinStrategy(s)
	} else {
		return WeightedLeastConnectionStrategy(s)
	}
}

//treats backend servers like a circular queue,
//handing out requests one by one in a fixed order.
func RoundRobinStrategy(s *ServerPool) *models.Server {
	// Atomically increment the counter and get the result.
	// Even if 1000 goroutines do this at once, they each get a unique 'next'.
	next := atomic.AddUint64(&s.Current, 1)
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	l := uint64(len(s.Backends))
	if l == 0 {
		return nil
	}

	// Find the next healthy server.
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
//Upgrade: weighted version :
//Each server has a weight (e.g., 1, 2, 3). A server with weight 3 gets 3 requests for every 1 request to a server with weight 1.
func WeightedRoundRobinStrategy(s *ServerPool) *models.Server {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	var bestServer *models.Server
	totalWeight := 0
	for _, srv := range s.Backends {
		if srv.IsAlive(){
			srv.CurrentWeight += srv.Weight
			totalWeight += srv.Weight
			if bestServer == nil || srv.CurrentWeight > bestServer.CurrentWeight{
				bestServer = srv
			}
		}
	}
	if bestServer == nil{
		return nil
	}
	bestServer.CurrentWeight -= totalWeight
	return bestServer
}

//looks at all available servers and picks the one
//currently handling the fewest active requests.
func LeastConnectionStrategy(s *ServerPool) *models.Server {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	var bestServer *models.Server
	minConn := -1

	for _, server := range s.Backends {
		if server.IsAlive() {
			weight := server.GetWeight()
            if weight <= 0 {
                weight = 1
            }
			conn := server.GetActiveConnections()
			if minConn == -1 || conn < minConn {
				minConn = conn
				bestServer = server
			}
		}
	}

	return bestServer
}
//update: weighted version :
//we calculate a "score" for each server: score = current_active_connections / weight.
//The server with the lowest score gets the next request.
func WeightedLeastConnectionStrategy(s *ServerPool) *models.Server {
	s.Mux.RLock()
	defer s.Mux.RUnlock()

	var bestServer *models.Server
	minScore := -1.0

	for _, server := range s.Backends {
		if server.IsAlive() {
			weight := server.GetWeight()
            if weight <= 0 {
                weight = 1
            }
			score := float64(server.GetActiveConnections()) / float64(weight)
			if minScore == -1 || score < minScore {
				minScore = score
				bestServer = server
			}else if score == minScore && weight > bestServer.GetWeight() {
            	bestServer = server
        }
		}
	}

	return bestServer
}

func (s *ServerPool) GetServerByID(id string) *models.Server{
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	for _ , srv := range s.Backends{
		if srv.ID == id{
			return srv
		}
	}
	return nil
}
