package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"bounceshield/email-validator-service/auth"
	"bounceshield/email-validator-service/jobs"
	"bounceshield/email-validator-service/models"
	"bounceshield/email-validator-service/verifier"
)

type queueRequest struct {
	Email string `json:"email"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&payload)

	if payload.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJWT(payload.Email)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Username == "" || payload.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := auth.RegisterUser(payload.Username, payload.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// email := r.Header.Get("User-Email")

	var payload queueRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Email == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	job := &models.Job{
		ID:           uuid.New().String(),
		EmailToCheck: payload.Email,
		Status:       "queued",
		Timestamp:    time.Now().Unix(),
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

func HandleSingleVerify(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
	}
	var req Req
	_ = json.NewDecoder(r.Body).Decode(&req)

	result := verifier.VerifyEmail(req.Email)

	res := map[string]string{
		"email":  req.Email,
		"status": result.Status,
	}
	json.NewEncoder(w).Encode(res)
}

func HandleBatchVerify(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		UserID string   `json:"user_id"`
		Emails []string `json:"emails"`
	}
	var req Req
	_ = json.NewDecoder(r.Body).Decode(&req)

	var results []models.Result
	for _, email := range req.Emails {
		result := verifier.VerifyEmail(email)
		results = append(results, models.Result{
			Email:  email,
			Status: result.Status,
			Reason: result.Reason,
		})
	}

	job := models.Job{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Emails:    req.Emails,
		Results:   results,
		Timestamp: time.Now().Unix(),
	}

	models.SaveJob(job.EmailToCheck, job.Owner)
	json.NewEncoder(w).Encode(job)
}

func JobHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := r.URL.Query().Get("user")
		json.NewEncoder(w).Encode(getJobs(userID))
	case http.MethodPost:
		var job models.Job
		_ = json.NewDecoder(r.Body).Decode(&job)

		models.SaveJob(job.EmailToCheck, job.Owner)
		json.NewEncoder(w).Encode(map[string]string{"status": "saved"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
