package jobs

import (
	"context"
	"fmt"
	"time"

	"bounceshield/models"
	"bounceshield/verifier"
)

func StartWorkerPool(n int) {
	for i := 0; i < n; i++ {
		go worker()
	}
}

func worker() {
	for job := range GetQueue() {
		processJob(job)
	}
}

func processJob(job *models.Job) {
	job.Status = "processing"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resultChan := make(chan *models.VerificationResult, 1)

	go func() {
		result := verifier.VerifyEmail(job.Email)
		resultChan <- result
	}()

	select {
	case <-ctx.Done():
		job.Status = "done"
		job.Result = &models.VerificationResult{
			Email:  job.Email,
			Status: "timeout",
			Reason: "SMTP check timed out",
		}
	case res := <-resultChan:
		job.Status = "done"
		job.Result = res
	}

	fmt.Printf("✅ Processed job %s (%s)\n", job.ID, job.Email)
}
