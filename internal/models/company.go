package models

import (
	"time"
	"gorm.io/gorm"
)

type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	Domain    string         `gorm:"unique;not null" json:"domain"`

	// 👇 DELETED: FormConfig field is gone.

	CreatedBy uint           `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}