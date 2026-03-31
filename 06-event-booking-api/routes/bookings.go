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

func getBookings(context *gin.Context) {
	event, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	total, err := models.GetBookingsCountByEventID(event.ID)

	if err != nil {
		log.Println("Error fetching bookings count:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch bookings count.",
		})

		return
	}

	lastID, _ := strconv.ParseInt(context.Query("lastId"), 10, 64)
	bookings, err := models.GetAllBookingsByEventID(event.ID, lastID)

	if err != nil {
		log.Println("Error fetching bookings items:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch bookings items.",
		})

		return
	}

	context.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": bookings,
	})
}

func createBooking(context *gin.Context) {
	event, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	booking := models.Booking{
		EventID: event.ID,
		UserID:  auth.GetUserID(context),
	}

	if err := booking.Save(); err != nil {
		switch err {
		case errs.ErrAlreadyExists:
			context.JSON(http.StatusConflict, gin.H{
				"error": "User is already booked for this event.",
			})

		default:
			log.Println("Error creating booking:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not create booking.",
			})
		}

		return
	}

	context.JSON(http.StatusCreated, booking)
}

func deleteBooking(context *gin.Context) {
	event, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	id, err := strconv.ParseInt(context.Param("bookingId"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not parse booking id.",
		})

		return
	}

	booking, err := models.GetBookingByID(id)

	if err != nil {
		switch err {
		case errs.ErrNotExists:
			context.JSON(http.StatusNotFound, gin.H{
				"error": "Booking not found.",
			})

		default:
			log.Println("Error fetching booking:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not fetch booking.",
			})
		}

		return
	}

	if booking.EventID != event.ID {
		context.JSON(http.StatusNotFound, gin.H{
			"error": "Booking not found for this event.",
		})

		return
	}

	userID := auth.GetUserID(context)

	if event.UserID != userID && booking.UserID != userID {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Only the event owner or booking owner can delete the booking.",
		})

		return
	}

	if err := booking.Delete(); err != nil {
		log.Println("Error deleting booking:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete booking.",
		})

		return
	}

	context.Status(http.StatusNoContent)
}
