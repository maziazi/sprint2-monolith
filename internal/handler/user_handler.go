// File: handler/auth_handler.go
package handler

import (
	"errors"
	"fitbyte/internal/service"
	"github.com/gin-gonic/gin"
	"log"
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

	token, _ := service.GenerateToken(user.Email, user.Id)
	c.JSON(http.StatusOK, gin.H{"email": user.Email, "token": token})
}

func RegisterUser(c *gin.Context) {
	log.Println("Handler RegisterUser hit") // Tambahkan log untuk melacak
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Service error: %v", err) // Tambahkan log error
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, err := service.GenerateToken(user.Email, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err) // Log jika token gagal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User registered successfully")
	c.JSON(http.StatusCreated, gin.H{"email": user.Email, "token": token})
}
