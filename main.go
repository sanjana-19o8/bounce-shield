package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/verify", handleSingleVerify)
	mux.HandleFunc("/api/batch-verify", handleBatchVerify)
	mux.HandleFunc("/api/jobs", jobHandler)

	// CORS middleware
	handler := cors.Default().Handler(mux)

	log.Println("🌐 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

// --- HANDLERS ---

func handleSingleVerify(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
	}
	var req Req
	_ = json.NewDecoder(r.Body).Decode(&req)

	status, reason := verifyEmail(req.Email)

	res := map[string]string{"email": req.Email, "status": status, "reason": reason}
	json.NewEncoder(w).Encode(res)
}

func handleBatchVerify(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		UserID string   `json:"user_id"`
		Emails []string `json:"emails"`
	}
	var req Req
	_ = json.NewDecoder(r.Body).Decode(&req)

	var results []Result
	for _, email := range req.Emails {
		status, reason := verifyEmail(email)
		results = append(results, Result{Email: email, Status: status, Reason: reason})
	}

	job := Job{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Emails:    req.Emails,
		Results:   results,
		Timestamp: time.Now().Unix(),
	}

	saveJob(job)
	json.NewEncoder(w).Encode(job)
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("user")
		json.NewEncoder(w).Encode(getJobs(userID))
	case http.MethodPost:
		var job Job
		_ = json.NewDecoder(r.Body).Decode(&job)
		saveJob(job)
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
