package auth

import "github.com/gin-gonic/gin"

func GetUserID(context *gin.Context) int64 {
	return getFromContext(context, "userID", int64(0))
}

func GetName(context *gin.Context) string {
	return getFromContext(context, "name", "")
}

func GetEmail(context *gin.Context) string {
	return getFromContext(context, "email", "")
}

func getFromContext[T any](context *gin.Context, key string, defaultValue T) T {
	value, exists := context.Get(key)

	if !exists {
		return defaultValue
	}

	return value.(T)
}
