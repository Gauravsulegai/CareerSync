package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	Email     string         `gorm:"unique;not null" json:"email"` // Personal Email
	Password  string         `json:"-"`
	Role      string         `json:"role"` // "student" or "employee"

	// --- NEW FIELDS FOR EMPLOYEES ---
	CompanyID   *uint        `json:"company_id"` // Link to Company Table (Pointer allows null for students)
	Company     *Company     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"company,omitempty"`
	
	Position    string       `json:"position"`   // e.g., "SDE II"
	WorkEmail   string       `json:"work_email"` // e.g., "gaurav@google.com"
	IsVerified  bool         `json:"is_verified" gorm:"default:false"` // Have they clicked the link?

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}