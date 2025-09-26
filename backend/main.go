package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		userId := c.Query("id")
		if userId == "1" {
			type User struct {
				Name      string
				StudentID string
				Major     string
			}
			user := []User{
				{
					Name:      "Supphawit Noowong",
					StudentID: "65122632",
					Major:     "COEAI"},
				{
					Name:      "Jimmy",
					StudentID: "65133437",
					Major:     "COEAI",
				},
			}

			return c.JSON(user)
		} else {
			return c.SendString("Don't have UserId : " + userId)
		}

	})

	app.Post("/users/information", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Users Postman")
	})

	app.Put("/users/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Users Put")
	})

	app.Delete("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Users Delete")
	})

	app.Patch("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Users Patch")
	})

	log.Println("🚀 Fiber running on :5000")
	log.Fatal(app.Listen(":5000"))
}
