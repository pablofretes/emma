package handlers

import (
	"net/http"
	"strings"
	"time"

	"emma/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type EventRestAdapter struct {
	EventRepository domain.EventRepository
}

func NewEventRESTAdapter(service domain.EventRepository) *EventRestAdapter {
	return &EventRestAdapter{EventRepository: service}
}

func (r *EventRestAdapter) FetchEvent(c *gin.Context) {
	id := c.Param("id")
	user := c.Keys["user"].(jwt.MapClaims)

	event, err := r.EventRepository.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if event.Status == "draft" && user["role"] != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (r *EventRestAdapter) CreateEvent(c *gin.Context) {
	var createEventRequest struct {
		Title            string    `json:"title" binding:"required,max=255"`
		ShortDescription string    `json:"short_description" binding:"required,max=500"`
		LongDescription  string    `json:"long_description" binding:"max=5000"`
		Date             string `json:"date" binding:"required"`
		Location         string    `json:"location" binding:"required,max=255"`
		Status           string    `json:"status" binding:"required,oneof=draft published"`
		OrganizerID      string    `json:"organizer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&createEventRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	organizerID, err := uuid.Parse(createEventRequest.OrganizerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID"})
		return
	}

	date, err := time.Parse("2006-01-02T15:04:05Z07:00", createEventRequest.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	event := domain.Event{
		Title:            createEventRequest.Title,
		ShortDescription: createEventRequest.ShortDescription,
		LongDescription:  createEventRequest.LongDescription,
		Date:             date,
		Location:         createEventRequest.Location,
		Status:           createEventRequest.Status,
		OrganizerID:      organizerID,
		Attendees:        []domain.User{},
	}

	createdEvent, err := r.EventRepository.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "data": createdEvent})
}

func (r *EventRestAdapter) UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	
	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	var updateEventRequest struct {
		Title            string    `json:"title" binding:"omitempty,max=255"`
		ShortDescription string    `json:"short_description" binding:"omitempty,max=500"`
		LongDescription  string    `json:"long_description" binding:"omitempty,max=5000"`
		Date             string    `json:"date" binding:"omitempty"`
		Location         string    `json:"location" binding:"omitempty,max=255"`
		Status           string    `json:"status" binding:"omitempty,oneof=draft published"`
		OrganizerID      string    `json:"organizer_id" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&updateEventRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var event domain.Event

	if updateEventRequest.Title != "" {
		event.Title = strings.ToLower(updateEventRequest.Title)
	}
	if updateEventRequest.ShortDescription != "" {
		event.ShortDescription = strings.ToLower(updateEventRequest.ShortDescription)
	}
	if updateEventRequest.LongDescription != "" {
		event.LongDescription = strings.ToLower(updateEventRequest.LongDescription)
	}
	if updateEventRequest.Date != "" {
		date, err := time.Parse("2006-01-02T15:04:05Z07:00", updateEventRequest.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		event.Date = date
	}
	if updateEventRequest.Location != "" {
		event.Location = strings.ToLower(updateEventRequest.Location)
	}
	if updateEventRequest.Status != "" {
		event.Status = strings.ToLower(updateEventRequest.Status)
	}
	if updateEventRequest.OrganizerID != "" {
		organizerID, err := uuid.Parse(updateEventRequest.OrganizerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organizer ID"})
			return
		}
		event.OrganizerID = organizerID
	}

	updatedEvent, err := r.EventRepository.UpdateEvent(id, event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updatedEvent,
	})
}

func (r *EventRestAdapter) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = r.EventRepository.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (r *EventRestAdapter) GetAllEvents(c *gin.Context) {
	var queryParams struct {
		Date   string `form:"date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
		Status string `form:"status" binding:"omitempty,oneof=draft published"`
		Title  string `form:"title" binding:"omitempty,max=255"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	date := queryParams.Date
	status := strings.ToLower(queryParams.Status)
	title := strings.ToLower(queryParams.Title)

	events, err := r.EventRepository.GetAllEvents(status, title, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    events,
	})
}