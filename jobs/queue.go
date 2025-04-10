package jobs

import (
	"bounceshield/models"
	"sync"
)

var (
	jobQueue   = make(chan *models.Job, 100)
	jobStorage = make(map[string]*models.Job)
	mu         sync.RWMutex
)

func EnqueueJob(job *models.Job) {
	mu.Lock()
	defer mu.Unlock()
	jobStorage[job.ID] = job
	jobQueue <- job
}

func GetJob(id string) *models.Job {
	mu.RLock()
	defer mu.RUnlock()
	return jobStorage[id]
}

func GetQueue() <-chan *models.Job {
	return jobQueue
}
