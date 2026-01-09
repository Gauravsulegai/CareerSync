package models

import (
	"time"
	"gorm.io/gorm"
)

type ReferralRequest struct {
	ID        uint `gorm:"primaryKey" json:"id"`

	// Relationships (Pointers)
	StudentID uint  `json:"student_id"`
	Student   *User `gorm:"foreignKey:StudentID" json:"student"`

	EmployeeID uint  `json:"employee_id"`
	Employee   *User `gorm:"foreignKey:EmployeeID" json:"employee"`

	// Status Tracking
	Status string `json:"status" gorm:"default:'Pending'"`

	// --- 🔴 NEW STANDARDIZED FIELDS 🔴 ---
	CandidateFirstName string `json:"first_name"`
	CandidateLastName  string `json:"last_name"`
	CandidateEmail     string `json:"email"`
	CandidateMobile    string `json:"mobile"`
	LinkedInURL        string `json:"linkedin_url"`
	ResumeURL          string `json:"resume_url"`
	
	// 👇 NEW FIELD ADDED HERE
	JobLink            string `json:"job_link"` 

	Motivation         string `json:"motivation"` // "Why fit?" (Max 100 words)

	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}