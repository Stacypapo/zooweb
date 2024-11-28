package controllers

import (
	"ZOOweb/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Getprofile(c *gin.Context, db *gorm.DB) {
	var user models.User

	// Извлечение userID из контекста
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Логируем значение userIDRaw для отладки
	fmt.Printf("Raw userID: %#v\n", userIDRaw)

	// Преобразуем userIDRaw (float64) в uint
	userID := uint(userIDRaw.(float64)) // Преобразуем из float64 в uint

	// Логируем преобразованный userID
	fmt.Printf("Converted userID: %d\n", userID)

	// Запрос к базе данных
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Логируем найденного пользователя для отладки
	fmt.Printf("Found user: %+v\n", user)

	// Передача данных в шаблон
	c.HTML(http.StatusOK, "profile.html", gin.H{"user": user})
}
