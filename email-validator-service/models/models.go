package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Job struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	Emails       []string `json:"emails"`
	Results      []Result `json:"results"`
	Timestamp    int64    `json:"timestamp"`
	EmailToCheck string   `json:"email_to_check"`
	Owner        string   `json:"owner"`
	Status       string   `json:"status"`
}

type Result struct {
	Email  string `json:"email"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func SaveJob(emailToCheck, owner string) Job {
	return Job{
		ID:           uuid.New().String(),
		EmailToCheck: emailToCheck,
		Status:       "queued",
		Owner:        owner,
		Timestamp:    time.Now().Unix(),
	}
}

func UpdateJob(job *Job) {
    query := `UPDATE jobs SET status = ?, timestamp = ? WHERE id = ?`
    _, err := DB.Exec(query, job.Status, job.Timestamp, job.ID)
    if err != nil {
        log.Printf("Failed to update job %s: %v", job.ID, err)
    }
}

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	createUserTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		password TEXT
	);`

	createJobTable := `CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		user_id TEXT,
		email_to_check TEXT,
		owner TEXT,
		status TEXT,
		timestamp INTEGER
	);`

	_, err = DB.Exec(createUserTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(createJobTable)
    if err != nil {
        log.Fatal("Failed to create jobs table:", err)
    }
}
