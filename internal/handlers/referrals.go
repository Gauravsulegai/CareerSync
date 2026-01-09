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
	EmployeeID uint   `json:"employee_id" binding:"required"`
	
	// The Standard Fields
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Mobile     string `json:"mobile" binding:"required"`
	LinkedIn   string `json:"linkedin_url" binding:"required,url"`
	Resume     string `json:"resume_url" binding:"required,url"`
	
	// 👇 NEW INPUT FIELD
	JobLink    string `json:"job_link" binding:"required"` 

	Motivation string `json:"motivation" binding:"required,max=1000"` 
}

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

	// 1. SPAM CHECK
	var existing models.ReferralRequest
	database.DB.Where("student_id = ? AND employee_id = ? AND status = 'Pending'", student.ID, input.EmployeeID).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "You already have a pending request with this employee."})
		return
	}

	// 2. Create the Request
	req := models.ReferralRequest{
		StudentID:          student.ID,
		EmployeeID:         input.EmployeeID,
		Status:             "Pending",
		
		// Map Input to DB Columns
		CandidateFirstName: input.FirstName,
		CandidateLastName:  input.LastName,
		CandidateEmail:     input.Email,
		CandidateMobile:    input.Mobile,
		LinkedInURL:        input.LinkedIn,
		ResumeURL:          input.Resume,
		
		// 👇 SAVE THE JOB LINK
		JobLink:            input.JobLink,

		Motivation:         input.Motivation,

		ExpiresAt:          time.Now().Add(5 * 24 * time.Hour),
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	// 3. Notify Employee
	notif := models.Notification{
		UserID:  input.EmployeeID,
		Message: "New Referral Request from " + input.FirstName + " " + input.LastName,
		Type:    "REQUEST_RECEIVED",
	}
	database.DB.Create(&notif)

	c.JSON(http.StatusCreated, gin.H{"message": "Request sent successfully!"})
}

// UpdateRequestStatus - Employee Accepts/Rejects
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

	if req.EmployeeID != employee.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot manage this request"})
		return
	}

	req.Status = input.Status
	database.DB.Save(&req)

	notif := models.Notification{
		UserID:  req.StudentID,
		Message: "Your referral request was " + input.Status + " by " + employee.Name,
		Type:    "STATUS_UPDATE",
	}
	database.DB.Create(&notif)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

// GetEmployeeRequests
func GetRequests(c *gin.Context) {
	userVal, _ := c.Get("user")
	employee := userVal.(models.User)

	if employee.Role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only employees can view requests"})
		return
	}

	var requests []models.ReferralRequest
	// Fetch requests that belong to this employee
	result := database.DB.Preload("Student").Where("employee_id = ?", employee.ID).Find(&requests)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}