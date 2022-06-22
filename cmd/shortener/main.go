package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syberpunkq/go_url_shortener/config"
	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

func main() {

	var stor *storage.Storage
	config := config.NewConfig()

	if config.FileStoragePath != "" {
		stor = storage.FileStorage(config.FileStoragePath)
	} else {
		stor = storage.NewStorage()
	}

	handlers := handlers.NewHandler(stor, config.BaseURL)

	appRouter := chi.NewRouter()
	appRouter.Get("/{id}", handlers.ShowHandler)
	appRouter.Post("/", handlers.CreateHandler)
	appRouter.Post("/api/shorten", handlers.ApiCreateHandler)

	http.ListenAndServe(config.ServerAddress, appRouter)
}
