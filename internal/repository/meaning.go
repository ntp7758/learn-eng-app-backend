package repository

import (
	"errors"
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/pkg/db"

	"gorm.io/gorm"
)

type MeaningRepository interface {
	Get(id uint) (Meanings *domain.Meaning, err error)
	GetByText(w string) (Meaning *domain.Meaning, err error)
	Add(Meaning domain.Meaning) error
	Update(Meaning domain.Meaning) error
	GetOrCreateMeaning(text string) (meaning *domain.Meaning, err error)
}

type meaningRepository struct {
	db *gorm.DB
}

func NewMeaningRepository(db db.PostgreSQLDBClient) (MeaningRepository, error) {
	dbClient := db.GetDB()
	err := dbClient.AutoMigrate(domain.Meaning{})
	if err != nil {
		return nil, err
	}
	return &meaningRepository{db: dbClient}, nil
}

func (r *meaningRepository) Add(Meaning domain.Meaning) error {
	tx := r.db.Create(&Meaning)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *meaningRepository) Get(id uint) (Meaning *domain.Meaning, err error) {
	tx := r.db.First(Meaning, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return Meaning, nil
}

func (r *meaningRepository) GetByText(w string) (Meaning *domain.Meaning, err error) {
	tx := r.db.Where("text = ?", w).Take(&Meaning)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return Meaning, nil
}

func (r *meaningRepository) Update(Meaning domain.Meaning) error {
	tx := r.db.Save(&Meaning)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *meaningRepository) GetOrCreateMeaning(text string) (meaning *domain.Meaning, err error) {
	err = r.db.Where("text = ?", text).Take(&meaning).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			meaning = &domain.Meaning{Text: text}
			if err := r.db.Create(&meaning).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return meaning, nil
}
