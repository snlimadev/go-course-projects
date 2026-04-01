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
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access token not provided.",
			})

			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		claims, err := auth.ValidateAccessToken(tokenString)

		if err != nil {
			log.Println("Error validating access token:", err)

			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access token invalid or expired.",
			})

			return
		}

		auth.SetContext(context, claims)
		context.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := context.GetHeader("Authorization")

		if authHeader != "" && strings.HasPrefix(authHeader, bearerPrefix) {
			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			claims, err := auth.ValidateAccessToken(tokenString)

			if err == nil {
				auth.SetContext(context, claims)
			}
		}

		context.Next()
	}
}
