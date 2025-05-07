package usecase

import (
	"errors"
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/internal/repository"
	"strings"

	"gorm.io/gorm"
)

type WordUsecase interface {
	GetAllWord() ([]domain.WordResponse, error)
	AddWord(req domain.WordRequest) error
}

type wordUsecase struct {
	wordRepo    repository.WordRepository
	meaningRepo repository.MeaningRepository
}

func NewWordUsecase(wordRepo repository.WordRepository, meaningRepo repository.MeaningRepository) WordUsecase {
	return &wordUsecase{wordRepo: wordRepo, meaningRepo: meaningRepo}
}

func (u *wordUsecase) GetAllWord() ([]domain.WordResponse, error) {
	words, err := u.wordRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := []domain.WordResponse{}
	for _, v := range words {
		word := domain.WordResponse{
			Word:          v.Word,
			Category:      v.Category,
			Score:         v.Score,
			PartsOfSpeech: v.PartsOfSpeech,
			GuessAccuracy: domain.WordGuessAccuracyResponse{
				Total:   v.GuessAccuracy.Total,
				Correct: v.GuessAccuracy.Correct,
				Wrong:   v.GuessAccuracy.Wrong,
			},
		}
		for _, w := range v.Meanings {
			word.Meanings = append(word.Meanings, w.Text)
		}
		res = append(res, word)
	}

	return res, nil
}

func (u *wordUsecase) AddWord(req domain.WordRequest) error {

	if strings.TrimSpace(req.Word) == "" || len(req.Meanings) == 0 || strings.TrimSpace(req.PartsOfSpeech) == "" {
		return errors.New("have empty input")
	}

	if !domain.WordPartsOfSpeechType[req.PartsOfSpeech] {
		return errors.New("parts of speech invalid")
	}

	w, err := u.wordRepo.GetByWordAndPartsOfSpeech(req.Word, req.PartsOfSpeech)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}

		word := domain.Word{
			Word:          req.Word,
			Category:      req.Category,
			Score:         0,
			PartsOfSpeech: req.PartsOfSpeech,
			GuessAccuracy: domain.WordGuessAccuracy{
				Total:   0,
				Correct: 0,
				Wrong:   0,
			},
		}

		for _, v := range req.Meanings {
			m, err := u.meaningRepo.GetOrCreateMeaning(v)
			if err != nil {
				return err
			}
			word.Meanings = append(word.Meanings, m)
		}

		err = u.wordRepo.Add(word)
		if err != nil {
			return err
		}
	} else {
		checkUpdate := false
		for _, reqM := range req.Meanings {
			for _, v := range w.Meanings {
				if v.Text == reqM {
					continue
				}

				m, err := u.meaningRepo.GetOrCreateMeaning(reqM)
				if err != nil {
					return err
				}
				w.Meanings = append(w.Meanings, m)
			}
		}

		if checkUpdate {
			err = u.wordRepo.Update(*w)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
