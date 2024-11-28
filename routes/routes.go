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
	}

	r.POST("/logout", func(c *gin.Context) {
		controllers.Logout(c)
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/home", func(c *gin.Context) {
		controllers.GetProducts(c, db)
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
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})
}
