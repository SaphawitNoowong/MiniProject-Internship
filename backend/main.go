package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

	// Ensure index on users.studentCode (unique) for fast lookups
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer idxCancel()
	usersColl := mongoClient.Database("mydb").Collection("users")
	model := mongo.IndexModel{
		Keys:    bson.D{{Key: "studentCode", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_studentCode_unique"),
	}
	if _, err := usersColl.Indexes().CreateOne(idxCtx, model); err != nil {
		log.Printf("warn: failed to create index on users.studentCode: %v", err)
	}
}

func main() {
	mustConnectMongo()
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		cur, err := coll.Find(ctx, bson.D{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		defer cur.Close(ctx)
		var users []bson.M
		if err := cur.All(ctx, &users); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"ok": true, "data": users})
	})

	app.Get("/users/:studentCode", func(c *fiber.Ctx) error {
		code := c.Params("studentCode")
		if code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "studentCode is required"})
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		var user bson.M
		err := coll.FindOne(ctx, bson.D{{Key: "studentCode", Value: code}}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"ok": false, "error": "not found"})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		return c.JSON(fiber.Map{"ok": true, "data": user})
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		type Student struct {
			StudentCode string    `json:"studentCode" bson:"studentCode"`
			Name        string    `json:"name" bson:"name"`
			Major       string    `json:"major" bson:"major"`
			CreatedAt   time.Time `json:"-" bson:"createdAt"`
		}

		var student Student
		if err := c.BodyParser(&student); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "invalid JSON body",
			})
		}
		if student.StudentCode == "" || student.Name == "" || student.Major == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"ok":    false,
				"error": "studentCode, name, and major are required",
			})
		}

		// reject duplicate by name
		{
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			coll := mongoClient.Database("mydb").Collection("users")
			var existing bson.M
			err := coll.FindOne(ctx, bson.D{{Key: "name", Value: student.Name}}).Decode(&existing)
			if err != nil && err != mongo.ErrNoDocuments {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
			}
			if err == nil {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"ok": false, "error": "name already exists"})
			}
		}

		student.CreatedAt = time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")
		res, err := coll.InsertOne(ctx, student)
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

	// Bulk insert students with duplicate-name validation
	app.Post("/users/many", func(c *fiber.Ctx) error {
		type Student struct {
			StudentCode string    `json:"studentCode" bson:"studentCode"`
			Name        string    `json:"name" bson:"name"`
			Major       string    `json:"major" bson:"major"`
			CreatedAt   time.Time `json:"-" bson:"createdAt"`
		}

		var students []Student
		if err := c.BodyParser(&students); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "invalid JSON body (array required)"})
		}
		if len(students) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "no students provided"})
		}

		// validate required fields and collect names
		nameSet := make(map[string]struct{})
		for i := range students {
			if students[i].StudentCode == "" || students[i].Name == "" || students[i].Major == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "studentCode, name, and major are required for all items"})
			}
			// check duplicate names within payload
			if _, seen := nameSet[students[i].Name]; seen {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "duplicate names in payload"})
			}
			nameSet[students[i].Name] = struct{}{}
			students[i].CreatedAt = time.Now()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		// check existing names in DB with one query
		var names []string
		for n := range nameSet {
			names = append(names, n)
		}
		count, err := coll.CountDocuments(ctx, bson.M{"name": bson.M{"$in": names}})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		if count > 0 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"ok": false, "error": "one or more names already exist"})
		}

		// perform bulk insert
		docs := make([]interface{}, 0, len(students))
		for i := range students {
			docs = append(docs, students[i])
		}
		res, err := coll.InsertMany(ctx, docs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"ok": true, "insertedCount": len(res.InsertedIDs)})
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
