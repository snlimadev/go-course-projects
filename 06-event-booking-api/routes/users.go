package routes

import (
	"log"
	"net/http"
	"strconv"

	"example.com/event-booking-api/auth"
	"example.com/event-booking-api/errs"
	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func createUser(context *gin.Context) {
	var user models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Save(); err != nil {
		switch err {
		case errs.ErrAlreadyExists:
			context.JSON(http.StatusConflict, gin.H{
				"error": "Email is already in use.",
			})

		default:
			log.Println("Error creating user:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not create user.",
			})
		}

		return
	}

	context.JSON(http.StatusCreated, user)
}

func deleteUser(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not parse user id.",
		})

		return
	}

	if id != auth.GetUserID(context) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Only the owner can delete the account.",
		})

		return
	}

	user := models.User{ID: auth.GetUserID(context)}

	if err := user.Delete(); err != nil {
		switch err {
		case errs.ErrNotExists:
			context.JSON(http.StatusNotFound, gin.H{
				"error": "User not found.",
			})

		default:
			log.Println("Error deleting user:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not delete user.",
			})
		}

		return
	}

	setRefreshTokenCookie(context, "", -1)
	context.Status(http.StatusNoContent)
}
