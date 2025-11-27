package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nati-d/spill-backend/features/nickname"
	"github.com/nati-d/spill-backend/middleware"
)

func main() {
	// Load environment variables when on dev mode only
	devMode := os.Getenv("DEV_MODE") == "true"
	if devMode {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	// Initialize Supabase client
	InitSupabase() // This will panic if initialization fails

	// Initialize nickname service with Supabase
	if err := nickname.InitSupabase(Supabase()); err != nil {
		log.Fatalf("Failed to initialize nickname service: %v", err)
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Telegram-Init-Data")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	})

	// Public routes (no authentication required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Private routes (authentication required)
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthRequired())
	{
		nickname.RegisterRoutes(authGroup)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	router.Run(":" + port)
}
