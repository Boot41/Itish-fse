package main

import (
	"io"
	"log"
	"os"

	"itish41/doctor_ai_assistant/initializers"
	"itish41/doctor_ai_assistant/middleware"
	"itish41/doctor_ai_assistant/route"

	"github.com/gin-gonic/gin"
)

func init() {
	if err := initializers.LoadEnv(); err != nil {
		log.Fatalf("[CRITICAL] Failed to load env: %s", err)
	}
	if err := initializers.ConnectDB(); err != nil {
		log.Fatalf("[CRITICAL] Failed to initialize database connection: %s", err)
	}
	// Uncomment to run migrations
	if err := initializers.Migrate(); err != nil {
		log.Fatalf("[CRITICAL] Failed to run database migrations: %s", err)
	}
}

func main() {
	// Create router with default logger and recovery middleware
	router := gin.New()

	// Add essential middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	f, err := os.Create("logFile.log") // creating a logger file to store the logger output
	if err != nil {
		log.Println("Trouble creating a logger")
	}
	defer f.Close()

	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // connecting gin default writer to write to file and and terminal

	route.SetupRoutes(router) // accessing the endpoints

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	// Set trusted proxies
	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Printf("Warning: Failed to set trusted proxies: %v\n", err)
	}

	// Start server
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
