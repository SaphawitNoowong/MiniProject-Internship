package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from Fiber",
		})
	})

	if err := app.Listen(":5000"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
