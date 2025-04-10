package main

import (
	"backend/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var jwtSecret = []byte("supersecretkey")

func Register(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	_, err := db.DB.Exec("INSERT INTO users(username, password) VALUES (?, ?)", body.Username, string(hash))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
	}
	return c.JSON(fiber.Map{"message": "User registered"})
}

func Login(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	var userID int
	var hashed string
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", body.Username).Scan(&userID, &hashed)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hashed), []byte(body.Password)) != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	signed, _ := token.SignedString(jwtSecret)
	return c.JSON(fiber.Map{"token": signed})
}
