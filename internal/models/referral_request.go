package models

import (
	"time"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReferralRequest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	
	// Relationships
	StudentID   uint           `json:"student_id"`
	Student     User           `gorm:"foreignKey:StudentID" json:"student"`
	
	EmployeeID  uint           `json:"employee_id"`
	Employee    User           `gorm:"foreignKey:EmployeeID" json:"employee"`

	// Status Tracking
	Status      string         `json:"status" gorm:"default:'Pending'"` 
	// Values: 'Pending', 'Accepted', 'Rejected', 'Expired', 'AutoCancelled'

	// The Data (Answers to the Company Form)
	FormData    datatypes.JSON `json:"form_data"` 
	// Example: {"resume": "link...", "job_id": "123", "note": "Hi!"}

	ExpiresAt   time.Time      `json:"expires_at"` // Calculated as CreatedAt + 5 Days

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}