package jobs

import (
	"context"
	"fmt"
	"time"

	"bounceshield/email-validator-service/models"
	"bounceshield/email-validator-service/verifier"
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

	resultChan := make(chan *verifier.VerificationResult, 1)

	go func() {
		result := verifier.VerifyEmail(job.EmailToCheck)
		resultChan <- &result // Send the address of result to the channel
		if result.Status == "ok" {
			job.Status = "done"
		} else {
			job.Status = "failed"
		}
		job.Results = append(job.Results, models.Result{
			Email:  result.Email,
			Status: result.Status,
			Reason: result.Reason,
		})
		models.UpdateJob(job) // Ensure this function is implemented in models.go
		close(resultChan)
		fmt.Printf("Job %s processed with result: %s\n", job.ID, result.Status)
	}()

	select {
	case <-ctx.Done():
		job.Status = "done"
		job.Results = append(job.Results, models.Result{
			Email:  job.EmailToCheck,
			Status: "timeout",
			Reason: "SMTP check timed out",
		})
		models.UpdateJob(job)
	case res := <-resultChan:
		job.Status = "done"
		job.Results = append(job.Results, models.Result{
			Email:  res.Email,
			Status: res.Status,
			Reason: res.Reason,
		})
		models.UpdateJob(job)
	}

	fmt.Printf("✅ Processed job %s (%s)\n", job.ID, job.EmailToCheck)
}
