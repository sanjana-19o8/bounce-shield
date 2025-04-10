package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"bounceshield/jobs"
	"bounceshield/models"
)

type queueRequest struct {
	Email string `json:"email"`
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload queueRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Email == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	job := &models.Job{
		ID:       uuid.New().String(),
		Email:    payload.Email,
		Status:   "queued",
		QueuedAt: time.Now(),
	}

	jobs.EnqueueJob(job)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func JobStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing job ID", http.StatusBadRequest)
		return
	}

	job := jobs.GetJob(id)
	if job == nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}
