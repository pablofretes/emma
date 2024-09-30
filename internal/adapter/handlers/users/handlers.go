package handlers

import (
	"net/http"
	"strings"

	"emma/internal/domain"
	"emma/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type UserRestAdapter struct {
	UserRepository domain.UserRepository
}

func NewUserRESTAdapter(repository domain.UserRepository) *UserRestAdapter {
	return &UserRestAdapter{UserRepository: repository}
}

func (r *UserRestAdapter) FetchUser(c *gin.Context) {
	id := c.Param("id")

	user, err := r.UserRepository.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *UserRestAdapter) CreateUser(c *gin.Context) {
	var createUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&createUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	hashedPassword, err := utils.HashPassword(createUserRequest.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	createUserRequest.Password = string(hashedPassword)

	user := domain.User{
		Username: createUserRequest.Username,
		Password: string(hashedPassword),
		Role:     createUserRequest.Role,
		AttendingEvents: []domain.Event{},
		OrganizedEvents: []domain.Event{},
	}

	err = r.UserRepository.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (r *UserRestAdapter) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"success": false,
		})
		return
	}

	user, err := r.UserRepository.GetByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"success": false,
		})
		return
	}

	match := utils.VerifyPassword(loginRequest.Password, user.Password)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "wrong username or password",
			"success": false,
		})
		return
	}

	claims := utils.SignTokenClaims{
		Id:       user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
	}

	token, err := utils.SignToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate token",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":   token,
	})
}

func (r *UserRestAdapter) UpdateUser(c *gin.Context) {
	eventId := c.Param("eventId")
	user := c.Keys["user"].(jwt.MapClaims)

	userId, ok := user["id"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Invalid user ID in token",
			"success": false,
		})
		return
	}

	err := r.UserRepository.Update(eventId, userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Failed to update user",
			"success": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
	})
}

func (r *UserRestAdapter) GetUsersEvents(c *gin.Context) {
	userID := c.Param("id")
	status := c.Query("status")
	status = strings.ToLower(status)
	if status != "" && status != "active" && status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid status parameter. Must be either 'active' or 'completed'",
			"success": false,
		})
		return
	}

	events, err := r.UserRepository.GetUsersEvents(userID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch events",
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    events,
	})
}
