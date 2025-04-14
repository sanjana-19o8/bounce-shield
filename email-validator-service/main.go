package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"

	"bounceshield/email-validator-service/api"
	"bounceshield/email-validator-service/auth"
	"bounceshield/email-validator-service/jobs"
	"bounceshield/email-validator-service/models"
)

func main() {
	log.Println("🚀 Starting server...")

	// Initialize the database
	models.InitDB()

	mux := http.NewServeMux()

	jobs.StartWorkerPool(5)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to BounceShield API!"))
	})

	// API endpoints
	mux.HandleFunc("/api/login", api.LoginHandler)
	mux.HandleFunc("/api/signin", api.RegisterHandler)
	mux.HandleFunc("/api/verify", api.VerifyHandler)
	mux.HandleFunc("/api/batch-verify", api.HandleBatchVerify)
	mux.HandleFunc("/api/jobs", api.JobHandler)

	mux.HandleFunc("/api/queue", auth.AuthMiddleware(api.QueueHandler))
	mux.HandleFunc("/api/status", auth.AuthMiddleware(api.JobStatusHandler))

	// CORS middleware
	handler := cors.Default().Handler(mux)

	log.Println("🌐 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", handler)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
