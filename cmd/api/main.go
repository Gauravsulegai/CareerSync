package main

import (
	"github.com/Gauravsulegai/careersync/internal/database"
	"github.com/Gauravsulegai/careersync/internal/handlers"
	"github.com/Gauravsulegai/careersync/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	app := gin.Default()

	// Public Routes
	app.POST("/signup", handlers.Signup)
	app.POST("/login", handlers.Login)
	app.GET("/companies/search", handlers.SearchCompanies) 

	// Temporary Verification Route
	app.GET("/verify/:id", func(c *gin.Context) {
		id := c.Param("id")
		database.DB.Model(&database.DB).Table("users").Where("id = ?", id).Update("is_verified", true)
		c.JSON(200, gin.H{"message": "Account Verified Successfully! You can now log in."})
	})

	// Protected Routes (Requires Login)
	protected := app.Group("/")
	protected.Use(middleware.RequireAuth) 
	{
		// Company Management
		protected.PUT("/company/form", handlers.UpdateCompanyForm)

		// Referral Logic
		protected.POST("/request/referral", handlers.SendReferralRequest)
		protected.PUT("/request/:id/status", handlers.UpdateRequestStatus)

		// Profile (Testing)
		protected.GET("/profile", func(c *gin.Context) {
			user, _ := c.Get("user")
			c.JSON(200, gin.H{"user": user})
		})
	}

	app.Run(":8080")
}