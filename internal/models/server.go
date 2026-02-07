package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	ID   string  `json:"id"`
	URL *url.URL `json:"url"`
	Alive  bool  `json:"alive"`
	ActiveConnections  int `json:"active_connections"`
	Weight int `json:"weight"`
	CurrentWeight int `json:"current_weight"`
	// These fields are for internal use only
	// We don't export them to JSON because they contain functions/locks
	Mux  sync.RWMutex  `json:"-"`  //Multiple Readers OR One Writer
	ReverseProxy *httputil.ReverseProxy  `json:"-"`
}

func (s *Server) IsAlive() bool {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return s.Alive
}

func (s *Server) SetAlive(alive bool){
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.Alive = alive
}

func (s *Server) GetActiveConnections() int {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return s.ActiveConnections
}

func (s *Server) IncrementConnections() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.ActiveConnections++
}

func (s *Server) DecrementConnections() {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	// Safety check to ensure we never go below zero
	if s.ActiveConnections > 0 {
		s.ActiveConnections--
	}
}


func (s *Server) GetWeight() int {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return s.Weight
}

func (s*Server) GetCurrentWeight() int {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return s.CurrentWeight
}