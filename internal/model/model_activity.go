package model

import "time"

var ValidActivityTypes = map[string]bool{
	"Walking":    true,
	"Yoga":       true,
	"Stretching": true,
	"Cycling":    true,
	"Swimming":   true,
	"Dancing":    true,
	"Hiking":     true,
	"Running":    true,
	"HIIT":       true,
	"JumpRope":   true,
}

type ActivityRequest struct {
	ActivityType      string `json:"activityType" binding:"required"`
	DoneAt            string `json:"doneAt" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	DurationInMinutes int    `json:"durationInMinutes" binding:"required,min=1"`
}

type ActivityResponse struct {
	ActivityID        string    `json:"activityId"`
	ActivityType      string    `json:"activityType"`
	DoneAt            string    `json:"doneAt"`
	DurationInMinutes int       `json:"durationInMinutes"`
	CaloriesBurned    int       `json:"caloriesBurned"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func IsValidActivityType(activityType string) bool {
	return ValidActivityTypes[activityType]
}
