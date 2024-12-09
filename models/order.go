package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID      uint    `gorm:"not null"`
	TotalAmount float64 `gorm:"not null"`
	Status      string  `gorm:"default:'pending'"` // Возможные статусы: pending, paid, cancelled
	PaymentID   string  // ID транзакции из платежной системы
	Items       []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID    uint    `gorm:"not null"`
	ProductID  uint    `gorm:"not null"`
	Quantity   uint    `gorm:"not null"`
	UnitPrice  float64 `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
}

func AutoMigrateOrders(db *gorm.DB) {
	db.AutoMigrate(&Order{})
}
func AutoMigrateOrderItems(db *gorm.DB) {
	db.AutoMigrate(&OrderItem{})
}
