package handler

import (
	"learn-eng-app-backend/internal/domain"
	"learn-eng-app-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type WordHandler interface {
	GetAllWord(c *fiber.Ctx) error
	AddWord(c *fiber.Ctx) error
	GetRandomWord(c *fiber.Ctx) error
	UpdateWordAccuracy(c *fiber.Ctx) error
}

type wordHandler struct {
	wordUsecase usecase.WordUsecase
}

func NewWordHandler(wordUsecase usecase.WordUsecase) WordHandler {
	return &wordHandler{wordUsecase: wordUsecase}
}

func (h *wordHandler) GetAllWord(c *fiber.Ctx) error {
	words, err := h.wordUsecase.GetAllWord()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    words,
	})
}

func (h *wordHandler) AddWord(c *fiber.Ctx) error {
	var req domain.WordRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.wordUsecase.AddWord(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}

func (h *wordHandler) GetRandomWord(c *fiber.Ctx) error {
	word, err := h.wordUsecase.GetRandomWord()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    word,
	})
}

func (h *wordHandler) UpdateWordAccuracy(c *fiber.Ctx) error {
	var req domain.UpdateWordQuizzRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.wordUsecase.UpdateWordAccuracy(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
	})
}
