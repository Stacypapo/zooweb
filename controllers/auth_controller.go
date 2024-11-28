package controllers

import (
	"ZOOweb/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"time"
)

var jwtKey = []byte("your_secret_key") // Секретный ключ для JWT

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func Signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}

func Signin(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

// Регистрация нового пользователя
func Register(c *gin.Context, db *gorm.DB) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// Парсинг данных из формы
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H{"error": "Все поля обязательны для заполнения"})
		return
	}

	// Создаем пользователя
	user := models.User{
		Username: form.Username,
	}
	if err := user.HashPassword(form.Password); err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	// Сохраняем в базе данных
	if err := db.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H{"error": "Пользователь уже существует"})
		return
	}
	c.Redirect(301, "/signin")
	//c.HTML(http.StatusOK, "signup.html", gin.H{"success": "Регистрация прошла успешно!"})
}

// Авторизация пользователя
func Login(c *gin.Context, db *gorm.DB) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	// Парсинг данных из формы
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "signin.html", gin.H{"error": "Все поля обязательны для заполнения"})
		return
	}

	// Проверяем наличие пользователя в базе
	var user models.User
	if err := db.Where("username = ?", form.Username).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "signin.html", gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Проверяем пароль
	if err := user.CheckPassword(form.Password); err != nil {
		c.HTML(http.StatusUnauthorized, "signin.html", gin.H{"error": "Неверное имя пользователя или пароль"})
		return
	}

	// Генерируем JWT-токен
	token, err := GenerateJWT(user.ID, user.Username)
	println(token)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "signin.html", gin.H{"error": "Ошибка генерации токена"})
		return
	}
	c.SetCookie(
		"token", // Имя cookie
		token,   // Значение токена
		3600*24, // Время жизни в секундах (1 день)
		"/",     // Доступ для всех путей
		"",      // Хост (оставляем пустым, чтобы использовать текущий)
		false,   // Использовать ли HTTPS (false для тестов, true для продакшна)
		true,    // HttpOnly (защита от XSS-атак)
	)
	// Успех
	c.Redirect(301, "/user/profile")
	//c.HTML(http.StatusOK, "signin.html", gin.H{"success": "Вход выполнен успешно!", "token": token})
}

func GenerateJWT(userID uint, userUsername string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": userUsername,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(jwtKey)
}

func Logout(c *gin.Context) {
	// Удаляем cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.Redirect(301, "/home")
}
