package models

import (
	"time"
	"gorm.io/gorm"
)

type ReferralRequest struct {
	ID        uint `gorm:"primaryKey" json:"id"`

	// Relationships
	StudentID uint  `json:"student_id"`
	Student   *User `gorm:"foreignKey:StudentID" json:"student"`

	// 👇 CHANGED: Employee is now OPTIONAL (Pointer *uint)
	// This allows it to be 'nil' (null) when the request is first sent.
	EmployeeID *uint `json:"employee_id"` 
	Employee   *User `gorm:"foreignKey:EmployeeID" json:"employee"`

	// 👇 NEW: We link request to the Company
	CompanyID uint     `json:"company_id"`
	Company   *Company `gorm:"foreignKey:CompanyID" json:"company"`

	Status string `json:"status" gorm:"default:'Pending'"`

	// Candidate Fields
	CandidateFirstName string `json:"first_name"`
	CandidateLastName  string `json:"last_name"`
	CandidateEmail     string `json:"email"`
	CandidateMobile    string `json:"mobile"`
	LinkedInURL        string `json:"linkedin_url"`
	ResumeURL          string `json:"resume_url"`
	
	// New Field
	JobLink            string `json:"job_link"`

	Motivation         string `json:"motivation"`

	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}