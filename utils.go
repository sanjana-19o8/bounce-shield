package main

import (
	"encoding/csv"
	"os"
	"strings"
	"sync"
)

func readEmailsFromCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var emails []string

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if len(record) > 0 {
			emails = append(emails, strings.TrimSpace(record[0]))
		}
	}
	return emails, nil
}

func writeResultsToCSV(results map[string]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Email", "Status"})
	for email, status := range results {
		writer.Write([]string{email, status})
	}
	return nil
}

func concurrentVerify(emails []string) map[string]string {
	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, email := range emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			status := verifyEmail(email)
			mu.Lock()
			results[email] = status
			mu.Unlock()
		}(email)
	}
	wg.Wait()
	return results
}
