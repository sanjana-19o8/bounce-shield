package api

import (
	"encoding/json"
	"net/http"
	"bounceshield/verifier"
)

type requestPayload struct {
	Email string `json:"email"`
}

type Job struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

var jobStore = []Job{}

func saveJob(job Job) {
	jobStore = append(jobStore, job)
}

func getJobs(userID string) []Job {
	var userJobs []Job
	for _, job := range jobStore {
		if job.UserID == userID {
			userJobs = append(userJobs, job)
		}
	}
	return userJobs
}

// POST /api/jobs
func SaveJobHandler(w http.ResponseWriter, r *http.Request) {
	var job Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	saveJob(job)
	json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
}

// GET /api/jobs?user=xyz
func GetJobsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	jobs := getJobs(userID)
	json.NewEncoder(w).Encode(jobs)
}

// GET /api/verify
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload requestPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.Email == "" {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	result := verifier.VerifyEmail(payload.Email)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}