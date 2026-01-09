package main

import (
	"log"

	"careersync/internal/database"
	"careersync/internal/handlers"
	"careersync/internal/middleware" 

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
		
		// 👇 NEW ROUTE ADDED HERE
		protected.GET("/requests", handlers.GetRequests)
	}

	log.Println("Server running on port 8080")
	r.Run(":8080")
}