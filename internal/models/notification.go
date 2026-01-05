package models

import (
	"time"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"` // Who receives this?
	Message   string         `json:"message"` // e.g., "New referral request from Gaurav"
	Type      string         `json:"type"`    // "REQUEST_RECEIVED", "STATUS_UPDATE", "SYSTEM"
	IsRead    bool           `json:"is_read" gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}