// POST /api/jobs
func saveJobHandler(w http.ResponseWriter, r *http.Request) {
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
func getJobsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user")
	jobs := getJobs(userID)
	json.NewEncoder(w).Encode(jobs)
}
