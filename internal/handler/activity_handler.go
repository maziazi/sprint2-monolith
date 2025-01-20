package handler

import (
	"fitbyte/internal/model"
	"fitbyte/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func CreateActivity(c *gin.Context) {
	log.Println("CreateActivity handler hit") // Debug log: Handler dipanggil

	var activity model.ActivityRequest
	if err := c.ShouldBindJSON(&activity); err != nil {
		log.Printf("Failed to bind JSON: %v", err) // Debug log: Error binding JSON
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !model.IsValidActivityType(activity.ActivityType) {
		log.Printf("Invalid ActivityType: %s", activity.ActivityType)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activityType"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println("UserID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Debug nilai dan tipe userID
	log.Printf("UserID raw value: %+v, Type: %T", userID, userID)

	// Konversi userID ke uint
	var userIDUint uint
	switch v := userID.(type) {
	case string:
		userIDInt, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("Failed to convert UserID from string to integer: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		userIDUint = uint(userIDInt)
	case float64:
		userIDUint = uint(v)
	case int:
		userIDUint = uint(v)
	default:
		log.Println("UserID is not a valid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Log hasil konversi
	log.Printf("Converted UserID: %d", userIDUint)

	// Lanjutkan proses dengan userIDUint
	log.Printf("Received activity data: %+v", activity) // Debug log: Data yang diterima
	response, err := service.CreateActivity(activity, userIDUint)
	if err != nil {
		log.Printf("Service error: %v", err) // Debug log: Error dari service
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Activity created successfully: %+v", response) // Debug log: Respons berhasil
	c.JSON(http.StatusCreated, response)
}
