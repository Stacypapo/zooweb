package controllers

import (
	"ZOOweb/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// Получить все продукты
func GetProducts(c *gin.Context, db *gorm.DB) {
	var products []models.Product
	db.Where("available = ?", true).Find(&products)
	c.HTML(http.StatusOK, "index.html", gin.H{"products": products})
}

// Получить детали продукта
func GetProduct(c *gin.Context, db *gorm.DB) {
	slug := c.Param("slug")
	var product models.Product
	if err := db.Where("slug = ?", slug).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.HTML(http.StatusOK, "product.html", gin.H{"product": product})
}

// Добавить новый продукт
func CreateProduct(c *gin.Context, db *gorm.DB) {
	var product models.Product
	var err error
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product.Name = c.Request.Form.Get("name")
	product.Slug = slug.Make(product.Name)
	stockStr := c.Request.Form.Get("stock")
	product.Stock, err = strconv.Atoi(stockStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock value"})
		return
	}
	product.Available = product.Stock > 0
	priceStr := c.Request.Form.Get("price")
	product.Price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price value"})
		return
	}
	product.Description = c.Request.Form.Get("description")
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload image"})
		return
	}

	filePath := fmt.Sprintf("uploads/%s", file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	product.Image = filePath
	db.Create(&product)
	c.Redirect(http.StatusFound, "/home")
}
func NewProduct(c *gin.Context) {
	c.HTML(http.StatusOK, "newproduct.html", nil)
}

// Уменьшить количество продукта
func DecreaseStock(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	amountStr := c.Param("amount")
	amount, _ := strconv.Atoi(amountStr)

	var product models.Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if product.Stock >= amount {
		product.Stock -= amount
		product.Available = product.Stock > 0
		db.Save(&product)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Not enough stock"})
	}
}
