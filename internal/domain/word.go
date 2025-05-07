package domain

import "gorm.io/gorm"

var (
	WordPartsOfSpeechType map[string]bool = map[string]bool{
		"Noun":         true, // คำนาม
		"Verb-1":       true, // กริยาช่อง 1
		"Verb-2":       true, // กริยาช่อง 2
		"Verb-3":       true, // กริยาช่อง 3
		"Adjective":    true, // คำคุณศัพท์
		"Adverb":       true, // คำวิเศษณ์
		"Pronoun":      true, // สรรพนาม
		"Preposition":  true, // คำบุพบท
		"Conjunction":  true, // คำสันธาน
		"Interjection": true, // คำอุทาน
	}
)

type Word struct {
	gorm.Model
	Word          string
	Meanings      []*Meaning `gorm:"many2many:word_meanings;"` // m-m
	Category      string
	Score         float32
	PartsOfSpeech string
	GuessAccuracy WordGuessAccuracy // 1-1
}

type Meaning struct {
	gorm.Model
	Text string `gorm:"unique"`
}

type WordGuessAccuracy struct {
	gorm.Model
	WordID  string
	Total   int
	Correct int
	Wrong   int
}

type WordRequest struct {
	Word          string   `json:"word"`
	Meanings      []string `json:"meanings"`
	Category      string   `json:"category"`
	PartsOfSpeech string   `json:"partsOfSpeech"`
}

type WordResponse struct {
	Word          string                    `json:"word"`
	Meanings      []string                  `json:"meanings"`
	Category      string                    `json:"category"`
	Score         float32                   `json:"score"`
	PartsOfSpeech string                    `json:"partsOfSpeech"`
	GuessAccuracy WordGuessAccuracyResponse `json:"guessAccuracy"`
}

type WordGuessAccuracyResponse struct {
	Total   int `json:"total"`
	Correct int `json:"correct"`
	Wrong   int `json:"wrong"`
}
