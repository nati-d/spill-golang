package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nati-d/spill-backend/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Supabase client
	InitSupabase() // This will panic if initialization fails

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
		// TODO: Add your route handlers here
		// Example:
		// authGroup.GET("/api/user", getUserHandler)
		// nickname.RegisterRoutes(authGroup)
		// auth.RegisterRoutes(authGroup)
		// confession.RegisterRoutes(authGroup)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port %s", port)
	router.Run(":" + port)
}
