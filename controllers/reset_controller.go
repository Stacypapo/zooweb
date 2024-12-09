package controllers

import (
	"ZOOweb/models"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func RequestPasswordReset(c *gin.Context, db *gorm.DB) {
	var form struct {
		Email string `form:"email" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный email"})
		return
	}

	// Проверяем, существует ли пользователь
	var user models.User
	if err := db.Where("email = ?", form.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Генерируем токен и сохраняем его в базе
	token := generateRandomToken()
	expirationTime := time.Now().Add(1 * time.Hour) // Токен действителен 1 час

	resetToken := models.ResetToken{
		Email:     form.Email,
		Token:     token,
		ExpiresAt: expirationTime,
	}
	db.Create(&resetToken)

	// Отправляем email
	resetLink := fmt.Sprintf("http://localhost:8080/password/reset/%s", token)
	sendEmail(form.Email, "Восстановление пароля", "Пройдите по ссылке для восстановления пароля: "+resetLink)

	c.JSON(http.StatusOK, gin.H{"message": "Ссылка для восстановления отправлена на email"})
}

// Генерация случайного токена
func generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Пример отправки email
func sendEmail(to, subject, body string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "stacypapo@mail.ru")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.mail.ru", 587, "stacypapo@mail.ru", "RyZ5uqF3M1Pw24zfLPXi")
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Ошибка при отправке письма:", err)
	}
}

func ResetPassword(c *gin.Context, db *gorm.DB) {
	var form struct {
		Token           string `form:"token" binding:"required"`
		NewPassword     string `form:"password" binding:"required"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Проверяем токен в базе
	var resetToken models.ResetToken
	if err := db.Where("token = ?", form.Token).First(&resetToken).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен недействителен"})
		return
	}

	// Проверяем срок действия токена
	if time.Now().After(resetToken.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен истёк"})
		return
	}

	// Обновляем пароль пользователя
	var user models.User
	if err := db.Where("email = ?", resetToken.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при шифровании пароля"})
		return
	}

	user.Password = string(hashedPassword)
	db.Save(&user)

	// Удаляем использованный токен
	db.Delete(&resetToken)

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлён"})
}
