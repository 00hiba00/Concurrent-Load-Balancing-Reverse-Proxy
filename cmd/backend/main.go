package main

import (
	"log"
	"net/http"
	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/backend_logic"
)

func main() {
	// 1. Get the router we built in internal/backend_logic/server.go
	router := backendlogic.NewRouter()

	// 2. Start the server on 8080
	log.Println("Backend API starting on :8080...")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}