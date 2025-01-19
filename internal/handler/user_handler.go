// File: handler/auth_handler.go
package handler

import (
	"errors"
	"fitbyte/internal/middleware"
	"fitbyte/internal/model"
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

	token, _ := middleware.GenerateToken(user.Email, user.Id)
	c.JSON(http.StatusOK, gin.H{"email": user.Email, "token": token})
}

func RegisterUser(c *gin.Context) {
	log.Println("Handler RegisterUser hit")
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input validated: %+v", req)
	user, err := service.RegisterUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Service error: %v", err)
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	token, _ := middleware.GenerateToken(user.Email, user.Id)
	if err != nil {
		log.Printf("Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	log.Println("User registered successfully")
	c.JSON(http.StatusCreated, gin.H{"email": user.Email, "token": token})
}

func GetUserProfileHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	user, err := service.GetUser(userIDInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"preference": user.Preference,
		"weightUnit": user.WeightUnit,
		"heightUnit": user.HeightUnit,
		"weight":     user.Weight,
		"height":     user.Height,
		"email":      user.Email,
		"name":       user.Name,
		"imageUri":   user.ImageUri,
	})
}

func UpdateUserProfileHandler(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDInt, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	var updatedProfile model.UserProfile
	if err := c.ShouldBindJSON(&updatedProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.PatchUser(userIDInt, &updatedProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{

		"preference": (*user).Preference,
		"weightUnit": (*user).WeightUnit,
		"heightUnit": (*user).HeightUnit,
		"weight":     (*user).Weight,
		"height":     (*user).Height,
		"email":      (*user).Email,
		"name":       (*user).Name,
		"imageUri":   (*user).ImageUri,
	})
}
