package main

import (
	"learn-eng-app-backend/internal/handler"
	"learn-eng-app-backend/internal/repository"
	"learn-eng-app-backend/internal/usecase"
	"learn-eng-app-backend/pkg/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	dbClient, err := db.NewPostgreSQLDBConnection()
	if err != nil {
		panic(err)
	}

	dbClient.SetLogger()

	wordRepo, err := repository.NewWordRepository(dbClient)
	if err != nil {
		panic(err)
	}
	meaningRepo, err := repository.NewMeaningRepository(dbClient)
	if err != nil {
		panic(err)
	}

	wordUsecase := usecase.NewWordUsecase(wordRepo, meaningRepo)

	wordHandler := handler.NewWordHandler(wordUsecase)

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/get-all", wordHandler.GetAllWord)
	app.Post("/add-word-or-meaning", wordHandler.AddWord)

	app.Listen(":8080")
}
