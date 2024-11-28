package models

import (
	"gorm.io/gorm"
)

// Product модель для базы данных
type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:100" json:"name"`
	Slug        string  `gorm:"unique" json:"slug"`
	Stock       int     `json:"stock"`
	Available   bool    `json:"available"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

// AutoMigrateProducts создает таблицу для продуктов
func AutoMigrateProducts(db *gorm.DB) {
	db.AutoMigrate(&Product{})
}
