package models

import "gorm.io/gorm"

// Модель для корзины
type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`             // Пользователь, которому принадлежит корзина
	ProductID uint    `gorm:"not null"`             // ID товара
	Quantity  uint    `gorm:"not null"`             // Количество товара
	Product   Product `gorm:"foreignKey:ProductID"` // Связь с товаром
}

// Автоматическая миграция для корзины
func AutoMigrateCart(db *gorm.DB) {
	db.AutoMigrate(&CartItem{})
}
