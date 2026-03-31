package middlewares

import (
	"log"
	"net/http"
	"strings"

	"example.com/event-booking-api/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := context.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token not provided.",
			})

			context.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := auth.ValidateAccessToken(tokenString)

		if err != nil {
			log.Println("Error validating access token:", err)

			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Access token invalid or expired.",
			})

			context.Abort()
			return
		}

		context.Set("userID", claims.UserID)
		context.Set("name", claims.Name)
		context.Set("email", claims.Email)

		context.Next()
	}
}
