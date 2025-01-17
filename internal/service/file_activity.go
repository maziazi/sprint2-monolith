package service

import (
	"fitbyte/internal/model"
	"time"
)

func CreateActivity(req model.ActivityRequest) (model.ActivityResponse, error) {
	// Simulate ID generation and calories calculation
	activityID := "act-" + time.Now().Format("20060102150405")
	caloriesBurned := req.DurationInMinutes * 5 // Example calculation
	now := time.Now()

	return model.ActivityResponse{
		ActivityID:        activityID,
		ActivityType:      req.ActivityType,
		DoneAt:            req.DoneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}
