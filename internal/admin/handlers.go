package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync/atomic"
	"text/template"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/balancer"
)


func GetBackendsHandler(s *balancer.ServerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		s.Mux.RLock()
		defer s.Mux.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(s.Backends)
		if err != nil {
			log.Printf("Failed to encode Backends: %v", err)
			return
		}
	}

}


func GetServerHandler(s *balancer.ServerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		s.Mux.RLock()
		defer s.Mux.RUnlock()
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
}

func PostServerHandler(s *balancer.ServerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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

    srv := balancer.NewServer(id, data.Addr)

    defer r.Body.Close()
	s.Mux.Lock()
	defer s.Mux.Unlock()
    s.Backends = append(s.Backends, srv)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(srv)
	}

}


func UpdateServerHandler(s *balancer.ServerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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
		s.Mux.Lock()
		defer s.Mux.Unlock()
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
}


func DeleteServerHandler(s *balancer.ServerPool) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		s.Mux.Lock()
		defer s.Mux.Unlock()
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
				if r.Header.Get("HX-Request") != "" {
    				w.WriteHeader(http.StatusOK) // HTMX needs 200 to swap content
				} else{
					w.WriteHeader(http.StatusNoContent)
				}
				return
			}
		}
		http.Error(w, "Server not found", http.StatusNotFound)
	}
}

func DashboardHandler(pool *balancer.ServerPool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Define the path to your HTML file
        // Note: This path is relative to where you run the 'go run' command
        tmplPath := filepath.Join("internal", "admin", "dashboard.html")

        // 2. Parse the template
        tmpl, err := template.ParseFiles(tmplPath)
        if err != nil {
            http.Error(w, "Failed to load dashboard: "+err.Error(), http.StatusInternalServerError)
            return
        }

        // 3. Execute: This maps your ServerPool data to the {{.Backends}} in HTML
        // We set the Content-Type to HTML so the browser renders it
        w.Header().Set("Content-Type", "text/html")
        err = tmpl.Execute(w, pool)
        if err != nil {
            http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
        }
    }
}