package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	StudentCode string    `json:"studentCode" bson:"studentCode"`
	Name        string    `json:"name" bson:"name"`
	Major       string    `json:"major" bson:"major"`
	CreatedAt   time.Time `json:"-" bson:"createdAt"`
}

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

	// Ensure index on users.studentCode (unique) for fast lookups
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer idxCancel()
	usersColl := mongoClient.Database("mydb").Collection("users")
	model := mongo.IndexModel{
		Keys:    bson.D{{Key: "studentCode", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_studentCode_unique"),
	}
	if _, err := usersColl.Indexes().CreateOne(idxCtx, model); err != nil {
		log.Printf("warn: failed to create index on users.studentCode: %v", err.Error())
	}
}

func main() {
	mustConnectMongo()
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		major := c.Query("major")
		name := c.Query("name")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		filter := bson.M{}
		if major != "" {
			filter["major"] = major
		}
		if name != "" {
			filter["name"] = name
		}

		cur, err := coll.Find(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Have some problem"})
		}
		defer cur.Close(ctx)
		var users []bson.M
		if err := cur.All(ctx, &users); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Something went wrong"})
		}
		if users == nil {
			return c.JSON(fiber.Map{"ok": false, "message": "Don't have user that you want"})
		}
		return c.JSON(fiber.Map{"ok": true, "data": users})
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		var students []Student
		if err := c.BodyParser(&students); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "invalid JSON body (array required)"})
		}
		if len(students) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "no students provided"})
		}
		// validate required fields and collect names
		for i := range students {
			if students[i].StudentCode == "" || students[i].Name == "" || students[i].Major == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "studentCode, name, and major are required for all items"})
			}
			students[i].CreatedAt = time.Now()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		// perform bulk insert
		docs := make([]interface{}, 0, len(students))
		for i := range students {
			docs = append(docs, students[i])
		}
		res, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Already have this student code."})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true, "insertedCount": len(res.InsertedIDs)})
	})

	app.Put("/users", func(c *fiber.Ctx) error {
		var students []Student
		if err := c.BodyParser(&students); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "invalid JSON body (array required)",
			})
		}
		if len(students) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "no students provided",
			})
		}

		// validate required fields
		for i := range students {
			if students[i].StudentCode == "" || students[i].Name == "" || students[i].Major == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"ok":    false,
					"error": "studentCode, name, and major are required for all items",
				})
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		var upsertedCount int64
		var modifiedCount int64
		var matchedCount int64

		for _, s := range students {
			filter := bson.M{"studentCode": s.StudentCode}
			update := bson.M{
				"$set": bson.M{
					"name":  s.Name,
					"major": s.Major,
				},
				"$currentDate": bson.M{
					"UpdatedAt": true,
				},
				"$onInsert": bson.M{
					"createdAt": true,
				},
			}

			res, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ok":    false,
					"error": err.Error(),
				})
			}

			matchedCount += res.MatchedCount
			modifiedCount += res.ModifiedCount
			upsertedCount += res.UpsertedCount
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":            true,
			"matchedCount":  matchedCount,
			"modifiedCount": modifiedCount,
			"upsertedCount": upsertedCount,
			"data":          students,
		})
	})

	app.Delete("/users", func(c *fiber.Ctx) error {
		studentCode := c.Query("studentCode")
		if studentCode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "studentCode is required",
			})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		filter := bson.M{"studentCode": studentCode}

		res, err := coll.DeleteOne(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Don't have this studentCode"})
		}
		if res.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ok":    false,
				"error": "Don't have this studentCode",
			})
		}
		return c.JSON(fiber.Map{"ok": true, "message": "Delete Succesful"})

	})

	app.Patch("/users", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! Users Patch")
	})

	log.Println("🚀 Fiber running on :5000")
	log.Fatal(app.Listen(":5000"))
}
