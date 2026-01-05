package models

import (
	"time"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Company struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null" json:"name"`   // e.g., "Google"
	Domain      string         `gorm:"unique;not null" json:"domain"` // e.g., "google.com" (The Lock)
	
	// The Shared "Master Form" Logic (JSON)
	// Example: {"require_resume": true, "custom_questions": ["Why us?"]}
	FormConfig  datatypes.JSON `json:"form_config"` 

	CreatedBy   uint           `json:"created_by"` // ID of the First Settler

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}