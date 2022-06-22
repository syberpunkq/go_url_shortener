package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

func main() {
	serverAddress, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		serverAddress = ":8080"
	}
	baseURL, ok := os.LookupEnv("BASE_URL")
	if !ok {
		baseURL = "http://localhost:8080"
	}

	storage := storage.NewStorage()
	handlers := handlers.NewHandler(storage, baseURL)

	appRouter := chi.NewRouter()
	appRouter.Get("/{id}", handlers.ShowHandler)
	appRouter.Post("/", handlers.CreateHandler)
	appRouter.Post("/api/shorten", handlers.ApiCreateHandler)

	http.ListenAndServe(serverAddress, appRouter)
}
