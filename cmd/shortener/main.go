package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

func main() {
	storage := storage.NewStorage()
	handlers := handlers.NewHandler(storage)

	appRouter := chi.NewRouter()
	appRouter.Get("/{id}", handlers.ShowHandler)
	appRouter.Post("/", handlers.CreateHandler)
	appRouter.Post("/api/shorten", handlers.ApiCreateHandler)

	http.ListenAndServe(":8080", appRouter)
}
