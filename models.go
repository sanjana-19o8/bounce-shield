package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Job struct {
	ID        string   `json:"id"`
	UserID    string   `json:"user_id"`
	Emails    []string `json:"emails"`
	Results   []Result `json:"results"`
	Timestamp int64    `json:"timestamp"`
}

type Result struct {
	Email   string `json:"email"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		filename TEXT,
		timestamp DATETIME,
		status TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id)
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
