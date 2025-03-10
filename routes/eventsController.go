package routes

import (
	"net/http"
	"strconv"

	"events.com/rest-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch events."})
		return
	}
	// context.JSON(http.StatusOK, events)
	context.JSON(http.StatusOK, gin.H{"data": events})
}

func createEvent(context *gin.Context) {
	var event models.Event
	err := context.ShouldBindJSON(&event)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	userId:=context.GetInt64("userId")
	// event.ID = 1S
	event.UserId = userId
	err = event.Save()

	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create event."})

	}
	context.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func getEvent(context *gin.Context){
	// eventId:= context.Param("id")
	eventId, err:= strconv.ParseInt(context.Param("id"), 10, 64)
	if err !=nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse eventId"})
		return
	}
	event, err := models.GetEventByID(eventId)
	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fecth event"})
		return
	}
	context.JSON(http.StatusOK, event)
}

func updateEvent(context *gin.Context){
	eventId, err:= strconv.ParseInt(context.Param("id"), 10, 64)
	if err !=nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse eventId"})
		return
	}

	userId:=context.GetInt64("userId") //userId from the token, was passed through the authorization
	event, err := models.GetEventByID(eventId)

	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch event"})
		return
	}

	if event.UserId != userId{
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to update this event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)

	if err !=nil{
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data" })
	}

	updatedEvent.ID= eventId
	err=updatedEvent.Update()
	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update Event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event updated successfully"})
}

func deleteEvent(context *gin.Context){
	// eventId:= context.Param("id")
	eventId, err:= strconv.ParseInt(context.Param("id"), 10, 64)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}

	userId:=context.GetInt64("userId") //userId from the token, was passed through the authorization

	event, err := models.GetEventByID(eventId)

	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fecth event"})
		return
	}

	if event.UserId != userId{
		context.JSON(http.StatusUnauthorized, gin.H{"message": "You are not authorized to delete this event"})
		return
	}

	err=event.Delete()
	if err !=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update Event"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "event deleted successfully"})
}