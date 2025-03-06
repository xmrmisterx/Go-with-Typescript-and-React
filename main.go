package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load .env file", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin,Content-Type,Accept",
	}))

	app.Get("/api/todos", getTodo)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5005"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodo(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	err := c.BodyParser(todo)
	if err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo id"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": "true"})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo id"})
	}

	filter := bson.M{"_id": objectId}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": "true"})
}

//// Below is in memory database version
//
//package main
//
//import (
//	"fmt"
//	"github.com/gofiber/fiber/v2"
//	"github.com/joho/godotenv"
//	"log"
//	"os"
//)
//
//type Todo struct {
//	ID        int    `json:"id"`
//	Completed bool   `json:"completed"`
//	Body      string `json:"body"`
//}
//
//func main() {
//	fmt.Println("hello worlds")
//	app := fiber.New()
//
//	err := godotenv.Load(".env")
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	PORT := os.Getenv("PORT")
//
//	var todos []Todo
//
//	//default route
//	app.Get("/api/todos", func(c *fiber.Ctx) error {
//		return c.Status(200).JSON(todos)
//	})
//
//	//add Todo
//	app.Post("/api/todos", func(c *fiber.Ctx) error {
//		todo := &Todo{}
//		if err := c.BodyParser(todo); err != nil {
//			return err
//		}
//
//		if todo.Body == "" {
//			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
//		}
//
//		todo.ID = len(todos) + 1
//		todos = append(todos, *todo)
//
//		return c.Status(201).JSON(todo)
//	})
//
//	//update Todo
//	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
//		id := c.Params("id")
//
//		for i, todo := range todos {
//			if fmt.Sprint(todo.ID) == id {
//				todos[i].Completed = true
//				return c.Status(200).JSON(todos[i])
//			}
//		}
//		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
//	})
//
//	//delete Todo
//	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
//		id := c.Params("id")
//
//		for i, todo := range todos {
//			if fmt.Sprint(todo.ID) == id {
//				todos = append(todos[:i], todos[i+1:]...)
//				return c.Status(200).JSON(fiber.Map{"success": "Todo deleted"})
//			}
//		}
//
//		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
//	})
//
//	log.Fatal(app.Listen(":" + PORT))
//	//log.Fatal(app.Listen(":4000"))
//}
