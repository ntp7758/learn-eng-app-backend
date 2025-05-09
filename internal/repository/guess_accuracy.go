package repository

import (
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/pkg/db"

	"gorm.io/gorm"
)

type GuessAccuracyRepository interface {
}

type guessAccuracyRepository struct {
	db *gorm.DB
}

func NewGuessAccuracyRepository(db db.PostgreSQLDBClient) (GuessAccuracyRepository, error) {
	dbClient := db.GetDB()
	err := dbClient.AutoMigrate(domain.GuessAccuracy{})
	if err != nil {
		return nil, err
	}
	return &guessAccuracyRepository{db: dbClient}, nil
}
