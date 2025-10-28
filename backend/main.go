package main

import (
	"context"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Student struct {
	StudentCode string    `json:"studentCode" bson:"studentCode"`
	Name        string    `json:"name" bson:"name"`
	Major       string    `json:"major" bson:"major"`
	CreatedAt   time.Time `json:"-" bson:"createdAt"`
	Password    string    `json:"-" bson:"password"`
	Role        string    `json:"role" bson:"role"`
}

type CreateStudentPayload struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Major    string `json:"major"`
}

type LoginPayload struct {
	StudentCode string `json:"studentCode"`
	Password    string `json:"password"`
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

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, https://mini-project-internship.vercel.app",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// app.Get("/users", func(c *fiber.Ctx) error {
	// 	keyword := c.Query("keyword")
	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()
	// 	coll := mongoClient.Database("mydb").Collection("users")
	// 	filter := bson.M{}
	// 	if keyword != "" {
	// 		filter = bson.M{
	// 			"$or": []bson.M{
	// 				{"name": keyword},
	// 				{"major": keyword},
	// 			},
	// 		}
	// 	}

	// 	cur, err := coll.Find(ctx, filter)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Have some problem"})
	// 	}
	// 	defer cur.Close(ctx)
	// 	var users []bson.M
	// 	if err := cur.All(ctx, &users); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "Something went wrong"})
	// 	}
	// 	if users == nil {
	// 		return c.JSON(fiber.Map{"ok": false, "message": "Don't have user that you want"})
	// 	}
	// 	return c.JSON(fiber.Map{"ok": true, "data": users})
	// })

	app.Get("/users", func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))

		// 1. รับค่า search จาก query string
		searchQuery := c.Query("search")

		skip := (page - 1) * limit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		// 2. สร้าง filter สำหรับ MongoDB
		filter := bson.M{"role": "nisit"} // ถ้าไม่มีการค้นหา, filter จะเป็น role ที่เป็น nisit เท่านั้น
		if searchQuery != "" {
			searchFilter := bson.M{
				"$or": []bson.M{
					{"studentCode": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"name": bson.M{"$regex": searchQuery, "$options": "i"}},
					{"major": bson.M{"$regex": searchQuery, "$options": "i"}},
				},
			}
			filter = bson.M{
				"$and": []bson.M{
					{"role": "nisit"},
					searchFilter,
				},
			}
		}
		// 3. ใช้ filter ที่สร้างขึ้นในการค้นหาและนับจำนวน
		var students []Student
		findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

		cursor, err := coll.Find(ctx, filter, findOptions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &students); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}

		totalRecords, err := coll.CountDocuments(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}
		totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

		return c.JSON(fiber.Map{
			"ok":   true,
			"data": students,
			"pagination": fiber.Map{
				"currentPage":  page,
				"totalPages":   totalPages,
				"totalRecords": totalRecords,
				"limit":        limit,
			},
		})
	})

	app.Get("/users/latest", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		var lastStudent Student
		// ค้นหาแค่ 1 record โดยเรียงจาก studentCode มากไปน้อย
		findOptions := options.FindOne().SetSort(bson.D{{Key: "studentCode", Value: -1}})
		err := coll.FindOne(ctx, bson.M{}, findOptions).Decode(&lastStudent)

		// กรณีที่ไม่เจอข้อมูลเลย (เป็นนิสิตคนแรกของระบบ)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// ส่งรหัสเริ่มต้นกลับไป (เช่น รหัส 0 เพื่อให้ frontend นำไป +1)
				return c.JSON(fiber.Map{
					"ok":   true,
					"data": fiber.Map{"studentCode": "66000000"},
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "database error"})
		}

		// ถ้าเจอข้อมูล, ส่งข้อมูลของนิสิตคนล่าสุดกลับไป
		return c.JSON(fiber.Map{
			"ok":   true,
			"data": lastStudent,
		})
	})

	// app.Post("/users", func(c *fiber.Ctx) error {
	// 	var students []Student
	// 	if err := c.BodyParser(&students); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "invalid JSON body (array required)"})
	// 	}
	// 	if len(students) == 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "no students provided"})
	// 	}
	// 	// validate required fields and collect names
	// 	for i := range students {
	// 		if students[i].StudentCode == "" || students[i].Name == "" || students[i].Major == "" {
	// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "studentCode, name, and major are required for all items"})
	// 		}
	// 	}

	// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// 	defer cancel()
	// 	coll := mongoClient.Database("mydb").Collection("users")

	// 	var upsertedCount int64

	// 	for _, s := range students {
	// 		filter := bson.M{"studentCode": s.StudentCode}
	// 		update := bson.M{
	// 			"$set": bson.M{
	// 				"name":  s.Name,
	// 				"major": s.Major,
	// 			},
	// 			"$currentDate": bson.M{
	// 				"createAt": true,
	// 			},
	// 		}

	// 		res, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	// 		if err != nil {
	// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 				"ok":    false,
	// 				"error": err.Error(),
	// 			})
	// 		}

	// 		upsertedCount += res.UpsertedCount
	// 	}

	// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 		"ok":            true,
	// 		"upsertedCount": upsertedCount,
	// 		"data":          students,
	// 	})
	// })

	app.Post("/users", func(c *fiber.Ctx) error {
		var payload CreateStudentPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "invalid JSON body"})
		}

		// Validate ข้อมูลที่ได้รับมา
		if payload.Password == "" || payload.Name == "" || payload.Major == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "password, name and major are required"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		// ค้นหานิสิตที่มีรหัสล่าสุดในฐานข้อมูล
		var lastStudent Student
		findOptions := options.FindOne().SetSort(bson.D{{Key: "studentCode", Value: -1}})
		err := coll.FindOne(ctx, bson.M{}, findOptions).Decode(&lastStudent)

		var nextStudentCodeStr string
		if err != nil {
			if err == mongo.ErrNoDocuments {
				nextStudentCodeStr = "65000000" // กำหนดรหัสเริ่มต้น
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "database error on find"})
			}
		} else {
			lastCode, _ := strconv.Atoi(lastStudent.StudentCode)
			nextStudentCode := lastCode + 1
			nextStudentCodeStr = strconv.Itoa(nextStudentCode)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": "failed to hash password"})
		}
		// สร้าง struct Student ที่สมบูรณ์เพื่อเตรียมบันทึก
		newStudent := Student{
			StudentCode: nextStudentCodeStr,
			Password:    string(hashedPassword),
			Name:        payload.Name,
			Major:       payload.Major,
			CreatedAt:   time.Now(),
			Role:        "nisit",
		}

		_, err = coll.InsertOne(ctx, newStudent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"ok": false, "error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"ok":   true,
			"data": newStudent,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var payload LoginPayload
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"ok": false, "error": "invalid JSON"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		// 4. ค้นหา User ด้วย StudentCode
		var user Student
		filter := bson.M{"studentCode": payload.StudentCode}
		err := coll.FindOne(ctx, filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// เคสที่หา studentCode นี้ไม่เจอ
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "error": "Didn't have this Student Code"})
				// คุณสามารถ return 404 Not Found หรือจัดการตาม logic ของคุณ
			} else {
				// เคสที่เกิด error อื่นๆ
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"ok": false, "error": "Have some problem please try again later"})
				// จัดการ error ทั่วไป
			}
		}
		// 5. ตรวจสอบรหัสผ่าน (สำคัญมาก: ในระบบจริงควรใช้ bcrypt.CompareHashAndPassword)
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"ok":    false,
				"error": "Your password is incorrect", // (แนะนำให้ใช้ error message ที่เป็นกลาง)
			})
		}

		// (ในระบบจริง: ควรสร้าง JWT Token ส่งกลับไปแทน)

		// 6. ถ้าสำเร็จ, ส่ง Role กลับไปให้ Frontend
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":   true,
			"role": user.Role, // <-- ส่ง Role กลับไป
		})
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

		var modifiedCount int64

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
			}

			res, err := coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ok":    false,
					"error": err.Error(),
				})
			}

			modifiedCount += res.ModifiedCount
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":            true,
			"modifiedCount": modifiedCount,
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
		var students []map[string]interface{}
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

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		coll := mongoClient.Database("mydb").Collection("users")

		var matchedCount int64
		var modifiedCount int64

		for _, s := range students {
			studentCode, ok := s["studentCode"].(string)
			if !ok || studentCode == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"ok":    false,
					"error": "studentCode is required for each student",
				})
			}

			// เอา studentCode ออกจาก payload กันไม่ให้แก้
			delete(s, "studentCode")

			update := bson.M{
				"$set": s,
				"$currentDate": bson.M{
					"updatedAt": true,
				},
			}

			res, err := coll.UpdateOne(ctx, bson.M{"studentCode": studentCode}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ok":    false,
					"error": err.Error(),
				})
			}

			matchedCount += res.MatchedCount
			modifiedCount += res.ModifiedCount
		}
		if matchedCount < int64(len(students)) {
			return c.JSON(fiber.Map{
				"ok":    false,
				"error": "One or more studentCode not found",
			})
		}
		return c.JSON(fiber.Map{
			"ok":            true,
			"matchedCount":  matchedCount,
			"modifiedCount": modifiedCount,
			"data":          students,
		})
	})

	log.Printf("🚀 Fiber running on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
