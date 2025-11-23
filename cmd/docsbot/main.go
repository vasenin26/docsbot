package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vasenin26/docsbot/internal/server"
	"github.com/vasenin26/docsbot/internal/telegram"
)

func main() {
	port := os.Getenv("METRICS_PORT")
	if port == "" {
		port = "9090"
	}
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting docsbot metrics server on %s", addr)

	// Telegram optional
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if token != "" {
		if err := telegram.StartInBackground(ctx, token); err != nil {
			log.Fatalf("failed to start telegram bot: %v", err)
		}
	}

	// Handle OS signals for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Printf("shutdown signal received")
		cancel()
		// allow some time for goroutines to stop
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()

	if err := server.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
