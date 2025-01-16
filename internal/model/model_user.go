package model

import "time"

type ModelUser struct {
	Id         uint      `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	Name       string    `json:"name"`
	Preference string    `json:"preference"`
	weightUnit string    `json:"weightUnit"`
	heightUnit string    `json:"heightInit"`
	weight     int       `json:"weight"`
	height     int       `json:"height"`
	imageUri   string    `json:"imageUri"`
}
