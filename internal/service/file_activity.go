package service

import (
	"context"
	"errors"
	"fitbyte/internal/model"
	"fitbyte/pkg/database"
	"time"
)

var activityTypes = map[string]float64{
	"Walking":    4,
	"Yoga":       4,
	"Stretching": 4,
	"Cycling":    8,
	"Swimming":   8,
	"Dancing":    8,
	"Hiking":     10,
	"Running":    10,
	"HIIT":       10,
	"JumpRope":   10,
}

// CreateActivity creates a new activity record in the database
func CreateActivity(req model.ActivityRequest, userID uint) (model.ActivityResponse, error) {
	db := database.GetDBPool()

	caloriesPerMinute, exists := activityTypes[req.ActivityType]
	if !exists {
		return model.ActivityResponse{}, errors.New("invalid activity type")
	}

	caloriesBurned := float64(req.DurationInMinutes) * caloriesPerMinute

	activityID := "act-" + time.Now().Format("20060102150405")
	now := time.Now()

	query := `INSERT INTO activities (id, "userId", "activityType", "doneAt", "durationInMinutes", "caloriesBurned", "createdAt", "updatedAt") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(context.Background(), query, activityID, userID, req.ActivityType, req.DoneAt, req.DurationInMinutes, int(caloriesBurned), now, now)
	if err != nil {
		return model.ActivityResponse{}, err
	}

	return model.ActivityResponse{
		ActivityID:        activityID,
		ActivityType:      req.ActivityType,
		DoneAt:            req.DoneAt,
		DurationInMinutes: req.DurationInMinutes,
		CaloriesBurned:    int(caloriesBurned),
		CreatedAt:         now,
		UpdatedAt:         now,
	}, nil
}
