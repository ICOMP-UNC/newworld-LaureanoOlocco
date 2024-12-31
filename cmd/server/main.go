package main

import (
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/database"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	//"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/repositories"
	"log"
)

// @tittle Fiber API Example
// @version 1
// @description This is a simple API example using Fiber
// @host localhost:3001
// @BasePath /
// @schemes http
// @produces json
// @consumes json
// contact:
//   name: Juan Diego

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// Connect to the database
	dbPool, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbPool.Close()

	// Initialize the database
	err = database.InitDB(dbPool)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString("¡Hola, mundo!")
		return nil
	})
	// user
	app.Post("/auth/register", handlers.Register)

	app.Listen(":3001")
}
