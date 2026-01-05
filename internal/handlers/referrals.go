package handlers

import (
	"encoding/json" // <--- Added
	"net/http"
	"time"

	"github.com/Gauravsulegai/careersync/internal/database"
	"github.com/Gauravsulegai/careersync/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes" // <--- Added to fix the type error
)

// SendReferralRequest - Student asks for a referral
// POST /request/referral
func SendReferralRequest(c *gin.Context) {
	userVal, _ := c.Get("user")
	student := userVal.(models.User)

	if student.Role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can send requests"})
		return
	}

	var input struct {
		EmployeeID uint        `json:"employee_id" binding:"required"`
		FormData   interface{} `json:"form_data" binding:"required"` // Generic JSON input
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. SPAM CHECK: Does a PENDING request already exist for this pair?
	var existing models.ReferralRequest
	database.DB.Where("student_id = ? AND employee_id = ? AND status = 'Pending'", student.ID, input.EmployeeID).First(&existing)
	if existing.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "You already have a pending request with this employee."})
		return
	}

	// 2. CONVERT FORM DATA TO JSON (The Fix!)
	formBytes, err := json.Marshal(input.FormData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid form data format"})
		return
	}

	// 3. Create the Request
	req := models.ReferralRequest{
		StudentID:  student.ID,
		EmployeeID: input.EmployeeID,
		Status:     "Pending",
		FormData:   datatypes.JSON(formBytes), // <--- CASTING DONE HERE
		ExpiresAt:  time.Now().Add(5 * 24 * time.Hour), // 5 Days from now
	}

	if err := database.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	// 4. Create In-App Notification for Employee
	notif := models.Notification{
		UserID:  input.EmployeeID,
		Message: "New Referral Request from " + student.Name,
		Type:    "REQUEST_RECEIVED",
	}
	database.DB.Create(&notif)

	c.JSON(http.StatusCreated, gin.H{"message": "Request sent successfully!"})
}

// UpdateRequestStatus - Employee Accepts/Rejects
// PUT /request/:id/status
func UpdateRequestStatus(c *gin.Context) {
	reqID := c.Param("id")
	userVal, _ := c.Get("user")
	employee := userVal.(models.User)

	var input struct {
		Status string `json:"status" binding:"required"` // "Accepted" or "Rejected"
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Find the Request
	var req models.ReferralRequest
	if err := database.DB.First(&req, reqID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	// 2. Security Check: Is this THIS employee's request?
	if req.EmployeeID != employee.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot manage this request"})
		return
	}

	// 3. Update Status
	req.Status = input.Status
	database.DB.Save(&req)

	// 4. NOTIFY STUDENT
	notif := models.Notification{
		UserID:  req.StudentID,
		Message: "Your referral request was " + input.Status + " by " + employee.Name,
		Type:    "STATUS_UPDATE",
	}
	database.DB.Create(&notif)

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}