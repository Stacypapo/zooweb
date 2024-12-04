package controllers

import (
	"ZOOweb/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// Добавить товар в корзину
func AddToCart(c *gin.Context, db *gorm.DB) {
	var form struct {
		ProductID uint `form:"product_id" binding:"required"`
		Quantity  uint `form:"quantity" binding:"required"`
	}

	// Извлекаем данные из формы
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	// Получаем ID пользователя из токена
	userIDRaw, _ := c.Get("userID")
	userID := uint(userIDRaw.(float64))

	// Проверяем, есть ли такой товар
	var product models.Product
	if err := db.First(&product, form.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	// Проверяем, есть ли этот товар уже в корзине
	var cartItem models.CartItem
	if err := db.Where("user_id = ? AND product_id = ?", userID, form.ProductID).First(&cartItem).Error; err == nil {
		// Если товар уже в корзине, обновляем количество
		cartItem.Quantity += form.Quantity
		db.Save(&cartItem)
	} else {
		// Если товара нет в корзине, создаем новую запись
		cartItem = models.CartItem{
			UserID:    userID,
			ProductID: form.ProductID,
			Quantity:  form.Quantity,
		}
		db.Create(&cartItem)
	}
	c.Redirect(301, c.Request.Referer())
}

// Просмотр корзины
func ViewCart(c *gin.Context, db *gorm.DB) {
	userIDRaw, _ := c.Get("userID")
	userID := uint(userIDRaw.(float64))

	var cartItems []models.CartItem
	db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems)

	c.HTML(http.StatusOK, "cart.html", gin.H{"cartItems": cartItems})
}

// Удалить товар из корзины
func RemoveFromCart(c *gin.Context, db *gorm.DB) {
	var form struct {
		CartItemID uint `form:"cart_item_id" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	db.Delete(&models.CartItem{}, form.CartItemID)
	c.Redirect(301, "/user/cart")
}
