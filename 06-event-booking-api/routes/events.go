package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/event-booking-api/auth"
	"example.com/event-booking-api/errs"
	"example.com/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	lastID, _ := strconv.ParseInt(context.Query("lastId"), 10, 64)
	hostID, _ := strconv.ParseInt(context.Query("hostId"), 10, 64)
	guestID, _ := strconv.ParseInt(context.Query("guestId"), 10, 64)

	events, err := models.GetAllEvents(lastID, hostID, guestID)

	if err != nil {
		log.Println("Error fetching events:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not fetch events.",
		})

		return
	}

	context.JSON(http.StatusOK, events)
}

func getEventDetails(context *gin.Context) {
	event, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {
	var event models.Event

	if err := context.ShouldBindJSON(&event); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.UserID = auth.GetUserID(context)

	if err := event.Save(); err != nil {
		log.Println("Error creating event:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create event.",
		})

		return
	}

	eventDetails := models.EventDetails{
		Event:               event,
		UserName:            auth.GetName(context),
		LoggedUserHasBooked: false,
		LoggedUserBookingID: nil,
	}

	context.Header("Location", fmt.Sprintf("/events/%d", event.ID))
	context.JSON(http.StatusCreated, eventDetails)
}

func updateEvent(context *gin.Context) {
	oldEvent, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	if oldEvent.UserID != auth.GetUserID(context) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Only the owner can update the event.",
		})

		return
	}

	var event models.Event

	if err := context.ShouldBindJSON(&event); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.ID = oldEvent.ID
	event.UserID = oldEvent.UserID
	event.CreatedAt = oldEvent.CreatedAt

	if err := event.Update(); err != nil {
		log.Println("Error updating event:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not update event.",
		})

		return
	}

	eventDetails := models.EventDetails{
		Event:               event,
		UserName:            oldEvent.UserName,
		LoggedUserHasBooked: oldEvent.LoggedUserHasBooked,
		LoggedUserBookingID: oldEvent.LoggedUserBookingID,
	}

	context.JSON(http.StatusOK, eventDetails)
}

func deleteEvent(context *gin.Context) {
	event, ok := fetchEventOrAbort(context)

	if !ok {
		return
	}

	if event.UserID != auth.GetUserID(context) {
		context.JSON(http.StatusForbidden, gin.H{
			"error": "Only the owner can delete the event.",
		})

		return
	}

	if err := event.Delete(); err != nil {
		log.Println("Error deleting event:", err)

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not delete event.",
		})

		return
	}

	context.Status(http.StatusNoContent)
}

func fetchEventOrAbort(context *gin.Context) (*models.EventDetails, bool) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not parse event id.",
		})

		return nil, false
	}

	event, err := models.GetEventDetailsByID(id, auth.GetUserID(context))

	if err != nil {
		switch err {
		case errs.ErrNotExists:
			context.JSON(http.StatusNotFound, gin.H{
				"error": "Event not found.",
			})

		default:
			log.Println("Error fetching event:", err)

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not fetch event.",
			})
		}

		return nil, false
	}

	return event, true
}
