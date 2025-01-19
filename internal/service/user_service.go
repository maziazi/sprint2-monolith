package service

import (
	"context"
	"errors"
	"fitbyte/internal/model"
	"fitbyte/pkg/database"
	"fmt"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	//jwtSecret             = []byte(getJWTSecret())
)

//func getJWTSecret() string {
//	secret := config.LoadEnv().JWTSecret
//	if secret == "" {
//		fmt.Println("⚠️  WARNING: JWT_SECRET tidak terbaca, gunakan default untuk debugging!")
//		secret = "default-secret-key"
//	}
//	return secret
//}

func RegisterUser(email, password string) (*model.User, error) {
	db := database.GetDBPool()

	// Check if email exists
	var existingUser model.User
	err := db.QueryRow(context.Background(), "SELECT email FROM users WHERE email = $1", email).Scan(&existingUser.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Insert user into database and get the generated ID
	var userID uint
	err = db.QueryRow(context.Background(),
		`INSERT INTO users (email, password, "createdAt") 
         VALUES ($1, $2, $3) 
         RETURNING id`,
		email, string(hashedPassword), time.Now(),
	).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %v", err)
	}

	return &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Id:       userID,
	}, nil
}
func Authenticate(email, password string) (*model.User, error) {
	db := database.GetDBPool()
	var user model.User

	// Retrieve user by email
	err := db.QueryRow(context.Background(), "SELECT id, email, password FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Email, &user.Password)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailNotFound
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}

//func GenerateToken(email string, userId uint) (string, error) {
//	claims := jwt.MapClaims{
//		"email":  email,
//		"userID": userId,
//		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 1 hari
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(jwtSecret)
//}
