package handlers

import (
	"net/http"
	"careersync/internal/database"
	"careersync/internal/models"
	"github.com/gin-gonic/gin"
)

// SearchCompanies - Returns a list of companies matching the query
func SearchCompanies(c *gin.Context) {
	query := c.Query("query")
	var companies []models.Company

	result := database.DB.Where("name ILIKE ?", "%"+query+"%").Find(&companies)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, companies)
}

// 👇 DELETED: UpdateCompanyForm is gone.