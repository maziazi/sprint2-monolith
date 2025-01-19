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

	_, err = db.Exec(context.Background(), `INSERT INTO profiles (email, "userId") VALUES ($1, $2)`,
		email, userID)
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

func GetUser(id uint) (*model.UserProfile, error) {
	db := database.GetDBPool()

	var user model.UserProfile
	err := db.QueryRow(context.Background(),
		`SELECT id, email, "weightUnit", "heightUnit", weight, height, name, "imageUri", preference
         FROM profiles 
         WHERE "userId" = $1`, id).
		Scan(&user.Id, &user.Email, &user.WeightUnit, &user.HeightUnit, &user.Weight, &user.Height, &user.Name, &user.ImageUri, &user.Preference)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %v", err)
		}
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return &user, nil
}

func PatchUser(id uint, updatedProfile *model.UserProfile) (**model.UserProfile, error) {
	db := database.GetDBPool()

	// Mendapatkan data user lama
	var oldUser model.UserProfile
	err := db.QueryRow(context.Background(), `SELECT id, email, "weightUnit", "heightUnit", weight, height, name, "imageUri", preference FROM profiles WHERE "userId" = $1`, id).
		Scan(&oldUser.Id, &oldUser.Email, &oldUser.WeightUnit, &oldUser.HeightUnit, &oldUser.Weight, &oldUser.Height, &oldUser.Name, &oldUser.ImageUri, &oldUser.Preference)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	// Menyusun query update dinamis
	updateQuery := `UPDATE profiles SET`
	params := []interface{}{}
	paramIndex := 1

	// Cek dan tambahkan field yang akan diperbarui

	// Periksa WeightUnit
	if updatedProfile.WeightUnit != nil {
		updateQuery += fmt.Sprintf(` "weightUnit" = $%d,`, paramIndex)
		params = append(params, updatedProfile.WeightUnit)
		paramIndex++
	}

	// Periksa HeightUnit
	if updatedProfile.HeightUnit != nil {
		updateQuery += fmt.Sprintf(` "heightUnit" = $%d,`, paramIndex)
		params = append(params, updatedProfile.HeightUnit)
		paramIndex++
	}

	// Periksa Weight (perbaikan: pointer dibandingkan dengan nilai langsung)
	if updatedProfile.Weight != nil {
		updateQuery += fmt.Sprintf(` weight = $%d,`, paramIndex)
		params = append(params, updatedProfile.Weight)
		paramIndex++
	}

	// Periksa Height (perbaikan: pointer dibandingkan dengan nilai langsung)
	if updatedProfile.Height != nil {
		updateQuery += fmt.Sprintf(` height = $%d,`, paramIndex)
		params = append(params, updatedProfile.Height)
		paramIndex++
	}

	// Periksa Name
	if updatedProfile.Name != nil {
		updateQuery += fmt.Sprintf(` name = $%d,`, paramIndex)
		params = append(params, updatedProfile.Name)
		paramIndex++
	}

	// Periksa ImageUri
	if updatedProfile.ImageUri != nil {
		updateQuery += fmt.Sprintf(` "imageUri" = $%d,`, paramIndex)
		params = append(params, updatedProfile.ImageUri)
		paramIndex++
	}

	// Periksa Preference
	if updatedProfile.Preference != nil {
		updateQuery += fmt.Sprintf(` preference = $%d,`, paramIndex)
		params = append(params, updatedProfile.Preference)
		paramIndex++
	}

	// Hapus koma terakhir jika ada perubahan
	if len(params) > 0 {
		updateQuery = updateQuery[:len(updateQuery)-1]
	}

	// Tambahkan kondisi WHERE untuk identifikasi user
	updateQuery += fmt.Sprintf(` WHERE "userId" = $%d RETURNING id`, paramIndex)
	params = append(params, id)

	// Jalankan query update
	var updatedId uint
	err = db.QueryRow(context.Background(), updateQuery, params...).Scan(&updatedId)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Mengambil data yang sudah diperbarui
	updatedProfile.Id = updatedId
	return &updatedProfile, nil

}
