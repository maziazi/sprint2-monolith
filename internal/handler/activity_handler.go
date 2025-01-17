package handler

import (
	"fitbyte/internal/model"
	"fitbyte/internal/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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

	log.Printf("Received activity data: %+v", activity) // Debug log: Data yang diterima
	response, err := service.CreateActivity(activity)
	if err != nil {
		log.Printf("Service error: %v", err) // Debug log: Error dari service
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Activity created successfully: %+v", response) // Debug log: Respons berhasil
	c.JSON(http.StatusCreated, response)
}
