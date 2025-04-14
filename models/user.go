package models

import (
	"errors"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func SaveUser(username, password string) error {
	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err := DB.Exec(query, username, password)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("user already exists")
		}
		return err
	}
	return nil
}

func GetUserByUsername(username string) *User {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	row := DB.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil
	}
	return &user
}
