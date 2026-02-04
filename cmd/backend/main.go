package main

import (
	"log"
	"net/http"
	"os"

	"github.com/00hiba00/Concurrent-Load-Balancing-Reverse-Proxy/internal/backend_logic"
)

func main() {
	// 1. Get the router we built in internal/backend_logic/server.go
	router := backendlogic.NewRouter()
	port := os.Args[1]
	// 2. Start the server on 8081
	log.Println("Backend API starting on :" + port +"...")
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}