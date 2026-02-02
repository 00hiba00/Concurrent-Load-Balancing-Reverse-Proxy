package balancer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)


func (s *ServerPool) GetBackendsHandler(w http.ResponseWriter, r *http.Request){
	s.mux.RLock()
	defer s.mux.RUnlock()
    w.Header().Set("Content-Type", "application/json")
    err := json.NewEncoder(w).Encode(s.Backends)
    if err != nil {
        log.Printf("Failed to encode Backends: %v", err)
        return
    }
}

func (s *ServerPool) GetServerHandler(w http.ResponseWriter, r *http.Request){
	s.mux.RLock()
	defer s.mux.RUnlock()
	IdFromUrl := r.PathValue("id")
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for _, currentServer := range s.Backends {
		if currentServer.ID == IdFromUrl {
			w.Header().Set("Content-Type", "application/json")
        	err := json.NewEncoder(w).Encode(currentServer)
			if err != nil {
				log.Printf("Failed to encode Backends: %v", err)
				return
    		}
        	return
		}
	}
	http.Error(w, "Server not found", http.StatusNotFound)
}

func (s *ServerPool) PostServerHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
        Addr string `json:"addr"`
    }
	err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
	newIDNum := atomic.AddUint64(&s.LastServerID, 1)
    id := fmt.Sprintf("srv-%d", newIDNum)

    srv := newServer(id, data.Addr)

    defer r.Body.Close()
	s.mux.Lock()
	defer s.mux.Unlock()
    s.Backends = append(s.Backends, srv)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(srv)
}

func (s *ServerPool) UpdateServerHandler(w http.ResponseWriter, r *http.Request){
	var data struct {
        Alive bool `json:"alive"`
    }
	err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Invalid JSON data", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
	IdFromUrl := r.PathValue("id")
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, srv := range s.Backends {
		if srv.ID == IdFromUrl {
			srv.SetAlive(data.Alive)
			w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(srv)
            return
		}
	}
	http.Error(w, "Server not found", http.StatusNotFound)
}

func (s *ServerPool) DeleteServerHandler(w http.ResponseWriter, r *http.Request){
	s.mux.Lock()
	defer s.mux.Unlock()
	IdFromUrl := r.PathValue("id")
	if IdFromUrl == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
        return
	}
	for i, currentServer := range s.Backends {
		if currentServer.ID == IdFromUrl {
			copy(s.Backends[i:], s.Backends[i+1:])
			s.Backends[len(s.Backends)-1] = nil
            s.Backends = s.Backends[:len(s.Backends)-1]
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Server not found", http.StatusNotFound)
}

