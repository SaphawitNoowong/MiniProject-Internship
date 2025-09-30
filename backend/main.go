package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func mustConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to create mongo client: %v", err)
	}
	// Verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping mongo: %v", err)
	}
	mongoClient = client
	log.Println("✅ Connected to MongoDB")
}

func main() {
	mustConnectMongo()
	app := fiber.New()

	app.Get("/health/db", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := mongoClient.Ping(ctx, nil); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ok":    false,
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{"ok": true})
	})

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
		type Student struct {
			StudentCode string    `json:"studentCode" bson:"studentCode"`
			Name        string    `json:"name" bson:"name"`
			Major       string    `json:"major" bson:"major"`
			CreatedAt   time.Time `json:"-" bson:"createdAt"`
		}

		var in Student
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "invalid JSON body",
			})
		}
		if in.StudentCode == "" || in.Name == "" || in.Major == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "studentCode, name, and major are required",
			})
		}

		in.CreatedAt = time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		res, err := coll.InsertOne(ctx, in)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":    false,
				"error": err.Error(),
			})
		}

		var id string
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			id = oid.Hex()
		} else {
			id = ""
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"ok": true,
			"id": id,
		})
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
