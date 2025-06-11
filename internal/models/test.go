package models

import (
	"time"
)

type Test struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UUID      string    `gorm:"type:uuid;uniqueIndex" json:"uuid"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedAt time.Time `gorm:"type:timestamp with time zone;index" json:"deleted_at"`
}
