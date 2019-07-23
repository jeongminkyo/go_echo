package main

import (
	"net/http"
	"context"
    "fmt"
    "log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// User
type User struct {
	ID    string `json:"id" form:"id" query:"id"`
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

var collection *mongo.Collection

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	
	if err != nil {
		log.Fatal(err)
	}
	
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("gofkadb").Collection("users")
	
	fmt.Println("Connected to MongoDB!")

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

// e.POST("/users", saveUser)
func saveUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Inserted a single document: ", insertResult, u)

	return c.JSON(http.StatusCreated, u)
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")

	var user User
	filter := bson.D{{"id", id}}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return c.String(http.StatusNotFound, id + " is not exist")
	}

	return c.JSON(http.StatusOK, &user)

}

// e.PUT("/users/:id", updateUser)
func updateUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	name := c.QueryParam("name")
	email := c.QueryParam("email")

	filter := bson.D{{"id", id}}

	update := bson.D{
		{"$set", bson.D{
			{"name", name},
			{"email", email},
		}},
	}

	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("update user information", updateResult)

	return c.JSON(http.StatusOK, id)
}

// e.DELETE("/users/:id", deleteUser)
func deleteUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")

	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{"id", id}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)

	return c.String(http.StatusOK, id+" is deleted")
}
