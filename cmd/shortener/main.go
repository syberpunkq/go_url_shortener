package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

func main() {
	// os.Setenv("FILE_STORAGE_PATH", "storage.txt")

	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		serverAddress = ":8080"
	}
	baseURL, ok := os.LookupEnv("BASE_URL")
	if !ok {
		baseURL = "http://localhost:8080"
	}

	var stor *storage.Storage
	fileStoragePath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if !ok {
		stor = storage.NewStorage()
	} else {
		stor = storage.FileStorage(fileStoragePath)
	}
	// storage := storage.NewStorage()
	handlers := handlers.NewHandler(stor, baseURL)

	appRouter := chi.NewRouter()
	appRouter.Get("/{id}", handlers.ShowHandler)
	appRouter.Post("/", handlers.CreateHandler)
	appRouter.Post("/api/shorten", handlers.ApiCreateHandler)

	http.ListenAndServe(serverAddress, appRouter)
}
