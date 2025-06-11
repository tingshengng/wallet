package main

import (
	"log"
	"os"

	"wallet/internal/database"
	"wallet/internal/migrations"
	"wallet/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	migrator := migrations.NewMigrator(db)
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	// POST /api/tests
	app.Post("/api/tests", func(c *fiber.Ctx) error {
		var test models.Test
		if err := c.BodyParser(&test); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if test.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name is required",
			})
		}
		test.UUID = uuid.New().String()

		if err := db.Create(&test).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(test)
	})

	// GET /api/tests
	app.Get("/api/tests", func(c *fiber.Ctx) error {
		var tests []models.Test
		if err := db.Find(&tests).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(tests)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
