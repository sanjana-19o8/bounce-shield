package main

import (
	"backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func getUserID(c *fiber.Ctx) int {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["user_id"].(float64))
}

func SaveJob(c *fiber.Ctx) error {
	userID := getUserID(c)
	type Payload struct {
		Filename string `json:"filename"`
		Status   string `json:"status"`
	}
	var body Payload
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).SendString("Invalid payload")
	}
	_, err := db.DB.Exec("INSERT INTO jobs(user_id, filename, timestamp, status) VALUES (?, ?, ?, ?)",
		userID, body.Filename, time.Now(), body.Status)
	if err != nil {
		return c.Status(500).SendString("Failed to save job")
	}
	return c.SendStatus(200)
}

func GetJobHistory(c *fiber.Ctx) error {
	userID := getUserID(c)
	rows, err := db.DB.Query("SELECT filename, timestamp, status FROM jobs WHERE user_id = ? ORDER BY timestamp DESC", userID)
	if err != nil {
		return c.Status(500).SendString("DB error")
	}
	defer rows.Close()

	var jobs []map[string]interface{}
	for rows.Next() {
		var fname, status string
		var ts string
		rows.Scan(&fname, &ts, &status)
		jobs = append(jobs, map[string]interface{}{
			"filename":  fname,
			"timestamp": ts,
			"status":    status,
		})
	}
	return c.JSON(jobs)
}
