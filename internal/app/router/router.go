package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
)

func New() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", handlers.ShowHandler)
	r.Post("/", handlers.CreateHandler)
	return r
}
