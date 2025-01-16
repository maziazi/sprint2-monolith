package handler

import (
	"errors"
	"fitbyte/internal/middleware"
	"fitbyte/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

func LoginUser(c *gin.Context) {

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := service.Authenticate(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		} else if errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, _ := middleware.GenerateToken(user.Email, user.Id)
	c.JSON(http.StatusOK, gin.H{"email": user.Email, "token": token})
	return
}

func RegisterUser(c *gin.Context) {
	var req AuthRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Proses registrasi user
	user, err := service.RegisterUser(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"}) // 409 Conflict
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{"user": user})
}
