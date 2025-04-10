package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/jwt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"backend/db"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	db.InitDB()

	app.Post("/register", Register)
	app.Post("/login", Login)

	api := app.Group("/api", jwt.New(jwt.Config{
		SigningKey: jwtSecret,
	}))

	api.Post("/save-job", SaveJob)
	api.Get("/job-history", GetJobHistory)

	app.Listen(":3000")
}


func handleCSVUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("❌ File upload failed")
	}

	path := "./uploads/" + file.Filename
	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).SendString("❌ File save failed")
	}

	emails, err := readEmailsFromCSV(path)
	if err != nil {
		return c.Status(500).SendString("❌ Could not parse CSV")
	}

	results := concurrentVerify(emails)

	outputFile := "./verified_" + file.Filename
	if err := writeResultsToCSV(results, outputFile); err != nil {
		return c.Status(500).SendString("❌ Failed to save results")
	}

	return c.JSON(results) // send results as JSON
}
