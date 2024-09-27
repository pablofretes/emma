package handlers

import (
	"net/http"

	users "emma/internal/core/domain/users"

	"github.com/gin-gonic/gin"
)

type UserRestAdapter struct {
	UserService users.UserService
}

func NewUserRESTAdapter(service users.UserService) *UserRestAdapter {
	return &UserRestAdapter{UserService: service}
}

func (r *UserRestAdapter) FetchUser(c *gin.Context) {
	username := c.Param("username")

	user, err := r.UserService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (r *UserRestAdapter) CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := r.UserService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}