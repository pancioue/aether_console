package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"aether-console/internal/db"
	"aether-console/internal/http/handlers"
	"aether-console/internal/services/health"
	"aether-console/internal/todo"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mysqlCfg := db.LoadConfig()
	mysqlDB, err := db.Open(mysqlCfg)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer mysqlDB.Close()

	// services
	healthSvc := health.New()

	// repos
	todoRepo := todo.Repo{DB: mysqlDB}

	// handlers
	healthH := handlers.NewHealthHandler(healthSvc)
	todoH := handlers.NewTodoHandler(todoRepo)

	// router (net/http)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/healthz", healthH.Healthz)
	mux.HandleFunc("/todos", todoH.List)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}
