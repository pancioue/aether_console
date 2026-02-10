package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"aether-console/internal/http/handlers"
	"aether-console/internal/services/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// services
	healthSvc := health.New()

	// handlers
	healthH := handlers.NewHealthHandler(healthSvc)

	// router (net/http)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthH.Healthz)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}