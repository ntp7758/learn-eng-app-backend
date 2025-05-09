package usecase

import (
	"errors"
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/internal/repository"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WordUsecase interface {
	GetAllWord() ([]domain.WordResponse, error)
	AddWord(req domain.WordRequest) error
	GetRandomWord() (*domain.WordResponse, error)
	UpdateWordAccuracy(req domain.UpdateWordQuizzRequest) error
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
			Score:         1,
			PartsOfSpeech: req.PartsOfSpeech,
			GuessAccuracy: domain.GuessAccuracy{
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
				checkUpdate = true
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

func (u *wordUsecase) GetRandomWord() (*domain.WordResponse, error) {
	words, err := u.wordRepo.GetAll()
	if err != nil {
		return nil, err
	}

	tmp := []domain.Word{}
	for _, v := range words {
		if len(tmp) == 0 {
			tmp = append(tmp, v)
		} else {
			if tmp[0].Score < v.Score {
				tmp = []domain.Word{}
				tmp = append(tmp, v)
			} else if tmp[0].Score == v.Score {
				tmp = append(tmp, v)
			}
		}
	}

	id := uint(0)
	res := domain.WordResponse{}
	if len(tmp) > 1 {
		rand.Seed(time.Now().UnixNano())
		w := tmp[rand.Intn(len(tmp))]
		id = w.ID
		res = domain.WordResponse{
			Word:          w.Word,
			Category:      w.Category,
			Score:         w.Score,
			PartsOfSpeech: w.PartsOfSpeech,
			GuessAccuracy: domain.WordGuessAccuracyResponse{
				Total:   w.GuessAccuracy.Total,
				Correct: w.GuessAccuracy.Correct,
				Wrong:   w.GuessAccuracy.Wrong,
			},
		}
		for _, v := range w.Meanings {
			res.Meanings = append(res.Meanings, v.Text)
		}
	}

	for i, v := range words {
		if v.ID == id {
			continue
		}

		words[i].Score++
	}

	err = u.wordRepo.BatchUpdateScores(words, 1000)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u *wordUsecase) UpdateWordAccuracy(req domain.UpdateWordQuizzRequest) error {

	if strings.TrimSpace(req.WordID) == "" {
		return errors.New("have empty input")
	}

	wid, err := strconv.Atoi(req.WordID)
	if err != nil {
		return err
	}

	word, err := u.wordRepo.Get(uint(wid))
	if err != nil {
		return err
	}

	if req.Correct {
		word.GuessAccuracy.Total++
		word.GuessAccuracy.Correct++
		word.Score = 1
	} else {
		word.GuessAccuracy.Total++
		word.GuessAccuracy.Wrong++
		word.Score += float32(1 + word.GuessAccuracy.Wrong)
	}

	err = u.wordRepo.UpdateAccuracy(*word)
	if err != nil {
		return err
	}

	return nil
}
