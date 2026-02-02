package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL *url.URL `json:"url"`
	Alive  bool  `json:"alive"`
	ActiveConnections  int `json:"active_connections"`
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