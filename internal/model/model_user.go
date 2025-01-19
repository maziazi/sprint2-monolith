package model

type User struct {
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"createdAt"`
}
type UserProfile struct {
	Id         uint    `json:"id"`
	Email      string  `json:"email"`
	WeightUnit *string `json:"weightUnit"`
	HeightUnit *string `json:"heightUnit"`
	Weight     *uint64 `json:"weight"`
	Height     *int64  `json:"height"`
	Name       *string `json:"name"`
	ImageUri   *string `json:"imageUri"`
	Preference *string `json:"preference"`
}
