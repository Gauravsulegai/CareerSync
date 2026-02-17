package handlers

import (
	"net/http"
	"time"

	"careersync/internal/database"
	"careersync/internal/models"
	"github.com/gin-gonic/gin"
)

// Define the exact input we expect from the frontend
type ReferralInput struct {
	// Asking for CompanyID (Broadcast Mode)
	CompanyID uint   `json:"company_id" binding:"required"`
	
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Mobile     string `json:"mobile" binding:"required"`
	LinkedIn   string `json:"linkedin_url" binding:"required,url"`
	Resume     string `json:"resume_url" binding:"required,url"`
	JobLink    string `json:"job_link" binding:"required"` 
	Motivation string `json:"motivation" binding:"required,max=1000"` 
}

// 1. SEND REQUEST (Broadcast to Company - No Notifications)
func SendReferralRequest(c *gin.Context) {
	userVal, _ := c.Get("user")
	student := userVal.(models.User)

	if student.Role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can send requests"})
		return
	}

	var input ReferralInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the Request (Unassigned: EmployeeID is nil)
	req := models.ReferralRequest{
		StudentID:          student.ID,
		CompanyID:          input.CompanyID, 
		EmployeeID:         nil,             
		Status:             "Pending",
		
		CandidateFirstName: input.FirstName,
		CandidateLastName:  input.LastName,
		CandidateEmail:     input.Email,
		CandidateMobile:    input.Mobile,
		LinkedInURL:        input.LinkedIn,
		ResumeURL:          input.Resume,
		JobLink:            input.JobLink,
		Motivation:         input.Motivation,
		ExpiresAt:          time.Now().Add(5 * 24 * time.Hour),
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Request broadcasted to company network."})
}

// UpdateRequestStatus - Employee Accepts (Claims) or Rejects
func UpdateRequestStatus(c *gin.Context) {
	reqID := c.Param("id")
	userVal, _ := c.Get("user")
	employee := userVal.(models.User)

	var input struct {
		Status string `json:"status" binding:"required"` 
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.ReferralRequest
	if err := database.DB.First(&req, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// Logic: If I accept it, I CLAIM it.
	if input.Status == "Accepted" {
		if req.EmployeeID != nil && *req.EmployeeID != employee.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Another employee has already accepted this request!"})
			return
		}
		req.EmployeeID = &employee.ID
	} else if input.Status == "Rejected" {
		if req.EmployeeID != nil && *req.EmployeeID != employee.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot reject a request claimed by someone else"})
			return
		}
		req.EmployeeID = &employee.ID
	}

	req.Status = input.Status
	database.DB.Save(&req)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// GetEmployeeRequests (Shows My Claims + Unclaimed Pool)
func GetRequests(c *gin.Context) {
	userVal, _ := c.Get("user")
	employee := userVal.(models.User)

	if employee.Role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only employees can view requests"})
		return
	}

	var requests []models.ReferralRequest
	
	// Show requests assigned to ME OR (belong to my company AND are unassigned)
	result := database.DB.Preload("Student").
		Where("(employee_id = ?) OR (company_id = ? AND employee_id IS NULL)", employee.ID, employee.CompanyID).
		Find(&requests)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// GetStudentRequests
func GetStudentRequests(c *gin.Context) {
	userVal, _ := c.Get("user")
	student := userVal.(models.User)

	var requests []models.ReferralRequest
	result := database.DB.Preload("Company").Preload("Employee").Where("student_id = ?", student.ID).Find(&requests)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, requests)
}