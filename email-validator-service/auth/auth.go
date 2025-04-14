package auth

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"bounceshield/email-validator-service/models"
)

var jwtSecret = []byte("supersecretkey")

// Register handles user registration
func Register(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Save the user to the database
	err = models.SaveUser(body.Username, string(hash))
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(400).JSON(fiber.Map{"error": "Username already exists"})
		}
		log.Printf("Error saving user: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// Login handles user login
func Login(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body Body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Fetch user details from the database
	user := models.GetUserByUsername(body.Username)
	if user == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{"token": signed})
}

// RegisterUser is a helper function to register a user
func RegisterUser(username, password string) error {
	// Check if the user already exists
	existingUser := models.GetUserByUsername(username)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Save the user to the database
	return models.SaveUser(username, string(hash))
}
