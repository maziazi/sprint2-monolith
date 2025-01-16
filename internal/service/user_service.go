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
)

func RegisterUser(email, password string) (*model.ModelUser, error) {
	db := database.GetDBPool()

	var existingUser model.ModelUser
	err := db.QueryRow(context.Background(), "SELECT email FROM users WHERE email = $1", email).Scan(&existingUser.Email)

	// Jika email sudah ada, kembalikan error 409 Conflict
	if err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) { // Tangani error database lain
		return nil, fmt.Errorf("Database Error: %v", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %v", err)
	}

	// Insert user baru ke database
	_, err = db.Exec(context.Background(), "INSERT INTO users (email, password, created_at) VALUES ($1, $2, $3)",
		email, string(hashedPassword), time.Now())

	if err != nil {
		return nil, fmt.Errorf("Failed to register user: %v", err)
	}

	// Return user yang berhasil dibuat
	return &model.ModelUser{Email: email, Password: string(hashedPassword), CreatedAt: time.Now()}, nil
}

func Authenticate(email, password string) (*model.ModelUser, error) {
	db := database.GetDBPool()
	var user model.ModelUser

	// Cari user berdasarkan email
	err := db.QueryRow(context.Background(), "SELECT id, email, password FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Email, &user.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrEmailNotFound
	} else if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Update timestamp login terakhir
	_, err = db.Exec(context.Background(), "UPDATE users SET last_login = $1 WHERE email = $2", time.Now(), user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to update last login timestamp: %v", err)
	}

	return &user, nil
}
