package models

import (
	"gorm.io/gorm"
	"time"
)

type ResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func AutoMigrateResetToken(db *gorm.DB) {
	db.AutoMigrate(&ResetToken{})
}
