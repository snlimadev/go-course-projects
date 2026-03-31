package routes

import (
	"log"
	"net/http"

	"example.com/event-booking-api/auth"
	"example.com/event-booking-api/config"
	"example.com/event-booking-api/errs"
	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func login(context *gin.Context) {
	var credentials models.UserCredentials

	if err := context.ShouldBindJSON(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := credentials.GetAuthenticatedUser()

	if err != nil {
		switch err {
		case errs.ErrInvalidCredentials:
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user credentials.",
			})

		default:
			log.Println("Error validating user credentials:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not validate user credentials.",
			})
		}

		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)

	if err != nil {
		log.Println("Error generating refresh token:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate refresh token.",
		})

		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Name, user.Email)

	if err != nil {
		log.Println("Error generating access token:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate access token.",
		})

		return
	}

	maxAge := config.JWTRefreshDurationDays * 24 * 3600

	setRefreshTokenCookie(context, refreshToken, maxAge)
	context.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func refresh(context *gin.Context) {
	refreshToken, err := context.Cookie("refreshToken")

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token not found. Please login again.",
		})

		return
	}

	claims, err := auth.ValidateRefreshToken(refreshToken)

	if err != nil {
		log.Println("Error validating refresh token:", err)
		setRefreshTokenCookie(context, "", -1)

		context.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token invalid or expired.",
		})

		return
	}

	user, err := models.GetUserByID(claims.UserID)

	if err != nil {
		log.Println("Error fetching user data:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch user data.",
		})

		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Name, user.Email)

	if err != nil {
		log.Println("Error refreshing access token:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not refresh access token. Please login again.",
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func logout(context *gin.Context) {
	setRefreshTokenCookie(context, "", -1)
	context.Status(http.StatusNoContent)
}

func setRefreshTokenCookie(context *gin.Context, token string, maxAge int) {
	isProd := config.Environment == "production"
	context.SetCookie("refreshToken", token, maxAge, "/", "", isProd, true)
}
