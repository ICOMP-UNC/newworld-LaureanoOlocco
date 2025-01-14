package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	//_ "github.com/ICOMP-UNC/newworld-rodriguezzfran/docs" // Import generated docs
	_ "github.com/lib/pq"

	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/database"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/core/services"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/handlers"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/repositories"
	"github.com/ICOMP-UNC/newworld-LaureanoOlocco/internal/server"
)

// @title Fiber API for new word project
// @version 1.2
// @description This Api makes the CRUD operations for testing the new world project
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host api.docker.localhost
// @BasePath /
func main() {

	// Load environment variables
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading .env file: %q", err)
	// }

	// check if the environment variables PORT and HOST are not empty
	var connStr string
	if os.Getenv("RUN_LOCAL") == "true" {
		connStr = fmt.Sprintf(
			"user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	} else {
		p := os.Getenv("DB_PORT")
		port, err := strconv.ParseUint(p, 10, 32) // Converting port string to int
		if err != nil {
			fmt.Println("Error parsing port str to int")
		}
		connStr = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			port,
		)
	}

	log.Printf("Environment variables:")
	log.Printf("RUN_LOCAL: '%s'", os.Getenv("RUN_LOCAL"))
	log.Printf("DB_HOST: '%s'", os.Getenv("DB_HOST"))
	log.Printf("DB_PORT: '%s'", os.Getenv("DB_PORT"))
	log.Printf("DB_NAME: '%s'", os.Getenv("DB_NAME"))
	log.Printf("DB_USER: '%s'", os.Getenv("DB_USER"))

	log.Printf("Connecting to database with connection string: %s", strings.Replace(connStr, os.Getenv("DB_PASSWORD"), "[REDACTED]", 1))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	defer db.Close()

	// Asure tables exist
	if err := database.EnsureTablesExist(db); err != nil {
		log.Fatalf("Error ensuring tables exist: %q", err)
	}

	// Init repositories, services and handlers
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandlers := handlers.NewUserHandlers(userService)

	// For offers
	offerRepository := repositories.NewOfferRepository(db)
	offerService := services.NewOfferService(offerRepository, userService)
	offerHandlers := handlers.NewOfferHandlers(offerService)

	// For admin
	adminRepository := repositories.NewAdminRepository(db)
	adminService := services.NewAdminService(adminRepository, offerService)
	adminHandlers := handlers.NewAdminHandlers(adminService)

	// Init server
	server := server.NewServer(userHandlers, offerHandlers, adminHandlers)

	// Start the server
	server.Start()
}
