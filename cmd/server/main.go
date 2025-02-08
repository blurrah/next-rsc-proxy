package main

import (
	"log"

	"github.com/blurrah/next-rsc-proxy/internal/proxy"
)

func main() {
	server := proxy.NewServer()
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Server started on port 8080")
}
