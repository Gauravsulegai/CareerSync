package main

import (
	"log"
	"os"

	"careersync/internal/database"
	"careersync/internal/handlers"
	"careersync/internal/middleware"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 👇 FIX: Allow the request to come from ANYWHERE (Dynamic Origin)
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Fallback for tools like Postman that don't send an Origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	database.ConnectDB()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// --- PUBLIC ROUTES ---
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.GET("/companies/search", handlers.SearchCompanies)

	// --- PROTECTED ROUTES ---
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth)
	{
		protected.POST("/request/referral", handlers.SendReferralRequest)
		protected.PUT("/request/:id/status", handlers.UpdateRequestStatus)
		protected.GET("/requests", handlers.GetRequests)
		protected.GET("/my-requests", handlers.GetStudentRequests)
	}

	// Use the port Render gives us, or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Println("Server running on port " + port)
	r.Run(":" + port)
}