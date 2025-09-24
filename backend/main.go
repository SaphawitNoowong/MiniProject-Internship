package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("mydb").Collection("user")

	app.Get("/api/users", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, map[string]interface{}{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to query database",
				"details": err.Error(),
			})
		}
		defer cursor.Close(ctx)

		var results []map[string]interface{}
		if err := cursor.All(ctx, &results); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to decode results",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"count":   len(results),
			"data":    results,
		})
	})

	log.Println("🚀 Fiber running on :5000")
	log.Fatal(app.Listen(":5000"))
}
