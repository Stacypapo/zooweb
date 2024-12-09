package controllers

import (
	"ZOOweb/models"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ViewOrders(c *gin.Context, db *gorm.DB) {
	// Retrieve userID from context
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID := uint(userIDRaw.(float64))

	// Fetch orders for the user from the database
	var orders []models.Order
	if err := db.Where("user_id = ?", userID).Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	// Check if there are no orders
	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No orders found for this user"})
		return
	}

	// Return the orders as JSON
	c.HTML(http.StatusOK, "orders.html", gin.H{"orders": orders})
	//c.JSON(http.StatusOK, orders)
}
func CheckoutHandler(c *gin.Context, db *gorm.DB) {
	userIDRaw, _ := c.Get("userID")
	userID := uint(userIDRaw.(float64))

	var cartItems []models.CartItem

	// Получаем товары из корзины пользователя
	if err := db.Preload("Product").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить данные корзины"})
		return
	}

	// Вычисляем общую стоимость заказа
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Product.Price
	}

	// Создаем новый заказ
	order := models.Order{
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
	}
	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать заказ"})
		return
	}

	// Создаем записи для каждого товара
	for _, item := range cartItems {
		orderItem := models.OrderItem{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.Product.Price,
			TotalPrice: float64(item.Quantity) * item.Product.Price,
		}
		db.Create(&orderItem)
	}

	// Удаляем корзину после оформления заказа
	db.Where("user_id = ?", userID).Delete(&models.CartItem{})

	// Генерируем ссылку для оплаты
	//paymentURL := GeneratePaymentLink(order.TotalAmount, order.ID)

	// Перенаправляем пользователя на страницу оплаты
	c.Redirect(301, "/user/order/success?order_id="+strconv.Itoa(int(order.ID)))
	//c.Redirect(301, paymentURL)
}

func OrderSuccessHandler(c *gin.Context, db *gorm.DB) {
	orderID := c.Query("order_id")

	// Обновляем статус заказа
	db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", "paid")

	c.HTML(http.StatusOK, "order_success.html", gin.H{"message": "Оплата прошла успешно!"})
}

func OrderCancelHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "order_cancel.html", gin.H{"message": "Оплата отменена. Попробуйте снова."})
}

func GeneratePaymentLink(amount float64, orderID uint) string {
	stripe.Key = "sk_test_rfaBmNUu1lWB7VZ0MMSIsYjH"

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Order #" + strconv.Itoa(int(orderID))),
					},
					UnitAmount: stripe.Int64(int64(amount * 100)), // Сумма в центах
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("http://localhost:8080/order/success?order_id=" + strconv.Itoa(int(orderID))),
		CancelURL:  stripe.String("http://localhost:8080/order/cancel"),
	}

	session_, _ := session.New(params)
	return session_.URL
}
