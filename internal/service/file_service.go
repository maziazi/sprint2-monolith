package service

import (
	"context"
	"database/sql"
	"errors"
	"fitbyte/internal/model"
	"fitbyte/pkg/database"
	"log"
)

func AddFile(fileURL string) (*model.File, error) {
	db := database.GetDBPool()
	var file model.File
	err := db.QueryRow(context.Background(), "INSERT INTO file (uri) VALUES ($1) RETURNING id, uri", fileURL).Scan(&file.ID, &file.URI)
	if err != nil {
		log.Printf("Error inserting file into database: %v", err)
		return nil, err
	}

	log.Printf("File stored in database: ID = %d, URI = %s", file.ID, file.URI)
	return &file, nil
}

func GetFileByID(id int) (*model.File, error) {
	db := database.GetDBPool()

	var file model.File
	err := db.QueryRow(context.Background(), "SELECT id, uri FROM file WHERE id = $1", id).Scan(&file.ID, &file.URI)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return &file, nil
}
