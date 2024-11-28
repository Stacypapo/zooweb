package main

import (
	"ZOOweb/models"
	"ZOOweb/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	// Параметры подключения к PostgreSQL
	dsn := "user=papo password=12345 dbname=shop host=109.120.185.19 port=5432 sslmode=disable"

	// Подключение к PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	println("Database connected!")

	// Автоматическое создание таблицы
	models.AutoMigrateProducts(db)
	models.AutoMigrateUsers(db)

	// Инициализация Gin
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	r.LoadHTMLGlob("templates/*") // Путь к папке с шаблонами

	// Настройка маршрутов
	routes.SetupRoutes(r, db)

	// Запуск сервера
	r.Run(":8080")
}
