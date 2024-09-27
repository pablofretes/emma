package handlers

import (
	"net/http"

	events "emma/internal/core/domain/events"

	"github.com/gin-gonic/gin"
)

type EventRestAdapter struct {
	EventService events.EventService
}

func NewEventRESTAdapter(service events.EventService) *EventRestAdapter {
	return &EventRestAdapter{EventService: service}
}

func (r *EventRestAdapter) FetchEvent(c *gin.Context) {
	id := c.Param("id")

	event, err := r.EventService.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (r *EventRestAdapter) CreateEvent(c *gin.Context) {
	var event events.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	createdEvent, err := r.EventService.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}

func (r *EventRestAdapter) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var event events.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	updatedEvent, err := r.EventService.UpdateEvent(id, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

func (r *EventRestAdapter) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	err := r.EventService.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (r *EventRestAdapter) GetAllEvents(c *gin.Context) {
	events, err := r.EventService.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, events)
}