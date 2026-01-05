package handlers

import (
	"net/http"

	"github.com/Gauravsulegai/careersync/internal/database"
	"github.com/Gauravsulegai/careersync/internal/models"
	"github.com/gin-gonic/gin"
)

// SearchCompanies - Returns a list of companies matching the query
// Example: GET /companies/search?query=goog
func SearchCompanies(c *gin.Context) {
	query := c.Query("query") // Get the text from URL
	var companies []models.Company

	// 1. Search DB (Case-insensitive fuzzy search)
	// LIKE %query% finds any company containing the text
	result := database.DB.Where("name ILIKE ?", "%"+query+"%").Find(&companies)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, companies)
}

// UpdateCompanyForm - Allows an employee to update the shared JSON form
// PUT /company/form
func UpdateCompanyForm(c *gin.Context) {
	// 1. Get User from Context (set by middleware)
	// We need to verify the user is logged in (which the middleware does),
	// but also that they belong to a company.
	userVal, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	currentUser := userVal.(models.User)

	// 2. Check if User is linked to a company
	if currentUser.CompanyID == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You must belong to a company to edit the form"})
		return
	}

	// 3. Bind the new JSON Config
	var input struct {
		FormConfig interface{} `json:"form_config" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. Update the Company Table
	result := database.DB.Model(&models.Company{}).
		Where("id = ?", *currentUser.CompanyID).
		Update("form_config", input.FormConfig)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update form"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Referral Form Updated Successfully!"})
}