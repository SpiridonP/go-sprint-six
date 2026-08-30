package main

import (
	"log"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	logger := log.Default()
	srv := server.NewServer(logger)

	logger.Printf("Starting server on :8080")
	if err := srv.Srv.ListenAndServe(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
