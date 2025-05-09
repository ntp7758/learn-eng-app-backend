package repository

import (
	"fmt"
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/pkg/db"
	"strings"

	"gorm.io/gorm"
)

type WordRepository interface {
	Get(id uint) (words *domain.Word, err error)
	GetByWord(w string) (word *domain.Word, err error)
	GetByWordAndPartsOfSpeech(w string, pos string) (word *domain.Word, err error)
	GetAll() (words []domain.Word, err error)
	Add(word domain.Word) error
	Update(word domain.Word) error
	UpdateAccuracy(word domain.Word) error
	BatchUpdateScores(records []domain.Word, batchSize int) error
	Delete(id uint) error
	PermanentDelete(id uint) error
}

type wordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db db.PostgreSQLDBClient) (WordRepository, error) {
	dbClient := db.GetDB()
	err := dbClient.AutoMigrate(domain.Word{})
	if err != nil {
		return nil, err
	}
	return &wordRepository{db: dbClient}, nil
}

func (r *wordRepository) Add(word domain.Word) error {
	tx := r.db.Create(&word)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *wordRepository) GetAll() (words []domain.Word, err error) {
	tx := r.db.Preload("Meanings").Find(&words)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return words, nil
}

func (r *wordRepository) Get(id uint) (word *domain.Word, err error) {
	tx := r.db.Preload("Meanings").Preload("GuessAccuracy").First(&word, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return word, nil
}

func (r *wordRepository) GetByWord(w string) (word *domain.Word, err error) {
	tx := r.db.Where("word = ?", w).Preload("Meanings").Take(&word)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return word, nil
}

func (r *wordRepository) GetByWordAndPartsOfSpeech(w string, pos string) (word *domain.Word, err error) {
	tx := r.db.Where("word = ?", w).Where("parts_of_speech = ?", pos).Preload("Meanings").Take(&word)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return word, nil
}

func (r *wordRepository) Update(word domain.Word) error {
	tx := r.db.Save(&word)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *wordRepository) UpdateAccuracy(word domain.Word) error {
	tx := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&word)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *wordRepository) BatchUpdateScores(records []domain.Word, batchSize int) error {
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}

		batch := records[i:end]

		var ids []string
		var cases []string

		for _, rec := range batch {
			ids = append(ids, fmt.Sprintf("%d", rec.ID))
			cases = append(cases, fmt.Sprintf("WHEN %d THEN %f", rec.ID, rec.Score))
		}

		sql := fmt.Sprintf(`
            UPDATE words
            SET score = CASE id
                %s
            END
            WHERE id IN (%s)
        `, strings.Join(cases, " "), strings.Join(ids, ","))

		if err := r.db.Exec(sql).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *wordRepository) Delete(id uint) error {
	tx := r.db.Delete(&domain.Word{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *wordRepository) PermanentDelete(id uint) error {
	tx := r.db.Unscoped().Delete(&domain.Word{}, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
