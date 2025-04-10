package main

import (
	"encoding/json"
	"os"
	"sync"
)

var jobFile = "data/jobs.json"
var mu sync.Mutex

func saveJob(job Job) error {
	mu.Lock()
	defer mu.Unlock()

	var jobs []Job
	data, _ := os.ReadFile(jobFile)
	if len(data) > 0 {
		_ = json.Unmarshal(data, &jobs)
	}

	jobs = append(jobs, job)
	newData, _ := json.MarshalIndent(jobs, "", "  ")
	return os.WriteFile(jobFile, newData, 0644)
}

func getJobs(userID string) []Job {
	mu.Lock()
	defer mu.Unlock()

	var jobs []Job
	data, _ := os.ReadFile(jobFile)
	if len(data) > 0 {
		_ = json.Unmarshal(data, &jobs)
	}

	var filtered []Job
	for _, j := range jobs {
		if j.UserID == userID {
			filtered = append(filtered, j)
		}
	}
	return filtered
}
