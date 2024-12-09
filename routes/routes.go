package routes

import (
	"ZOOweb/controllers"
	"ZOOweb/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	secretKey := "your_secret_key"
	authGroup := r.Group("/user")
	authGroup.Use(middleware.JWTAuthMiddleware(secretKey)) // Используем middleware
	{
		authGroup.GET("/profile", func(c *gin.Context) {
			controllers.Getprofile(c, db)
		})

		authGroup.POST("/cart/add", func(c *gin.Context) {
			controllers.AddToCart(c, db)
		})

		authGroup.GET("/cart", func(c *gin.Context) {
			controllers.ViewCart(c, db)
		})

		authGroup.POST("/cart/remove", func(c *gin.Context) {
			controllers.RemoveFromCart(c, db)
		})

	}
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})

	r.GET("/home", func(c *gin.Context) {
		controllers.GetProducts(c, db)
	})

	r.GET("/search", func(c *gin.Context) {
		controllers.SearchProducts(c, db)
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})

	r.GET("/signin", func(c *gin.Context) {
		c.HTML(200, "signin.html", nil)
	})

	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})

	r.POST("/logout", func(c *gin.Context) {
		controllers.Logout(c)
	})
	
	r.GET("/password/forgot", func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgot_password.html", nil)
	})

	r.POST("/password/forgot", func(c *gin.Context) {
		controllers.RequestPasswordReset(c, db)
	})

	r.GET("/password/reset/:token", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reset_password.html", gin.H{"token": c.Param("token")})
	})

	r.POST("/password/reset", func(c *gin.Context) {
		controllers.ResetPassword(c, db)
	})

	r.GET("/product/:slug", func(c *gin.Context) {
		controllers.GetProduct(c, db)
	})

	r.POST("/newproduct", func(c *gin.Context) {
		controllers.CreateProduct(c, db)
	})

	r.GET("/create", func(c *gin.Context) {
		controllers.NewProduct(c)
	})

	r.POST("/decrease_stock/:id/:amount", func(c *gin.Context) {
		controllers.DecreaseStock(c, db)
	})

	r.GET("/terms", func(c *gin.Context) {
		c.HTML(http.StatusOK, "terms.html", nil)
	})
}
