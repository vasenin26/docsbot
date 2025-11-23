package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vasenin26/docsbot/internal/server"
)

func main() {
	port := os.Getenv("METRICS_PORT")
	if port == "" {
		port = "9090"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting docsbot metrics server on %s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
