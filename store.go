package main

import (
	"encoding/json"
	"os"
	"sync"

	"bounceshield/models"
)

var jobFile = "data/jobs.json"
var mu sync.Mutex

func saveJob(job models.Job) error {
	mu.Lock()
	defer mu.Unlock()

	var jobs []models.Job
	data, _ := os.ReadFile(jobFile)
	if len(data) > 0 {
		_ = json.Unmarshal(data, &jobs)
	}

	jobs = append(jobs, job)
	newData, _ := json.MarshalIndent(jobs, "", "  ")
	return os.WriteFile(jobFile, newData, 0644)
}

func getJobs(userID string) []models.Job {
	mu.Lock()
	defer mu.Unlock()

	var jobs []models.Job
	data, _ := os.ReadFile(jobFile)
	if len(data) > 0 {
		_ = json.Unmarshal(data, &jobs)
	}

	var filtered []models.Job
	for _, j := range jobs {
		if j.UserID == userID {
			filtered = append(filtered, j)
		}
	}
	return filtered
}
