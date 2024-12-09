package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем токен из cookies
		tokenString, err := c.Cookie("token")
		if err != nil {
			// Если токен отсутствует, перенаправляем на страницу входа
			c.Redirect(http.StatusFound, "/signin") // Редирект на страницу входа
			c.Abort()
			return
		}

		// Проверяем токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{
					Err:  http.ErrAbortHandler,
					Type: gin.ErrorTypePublic,
					Meta: "Invalid signing method",
				}
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			// Если токен недействителен, перенаправляем на страницу входа
			c.Redirect(http.StatusFound, "/signin") // Редирект на страницу входа
			c.Abort()
			return
		}

		// Извлекаем данные пользователя из токена и сохраняем в контексте
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["id"])
		}

		// Продолжаем обработку запроса
		c.Next()
	}
}
