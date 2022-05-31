package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

// GET /{id}
func ShowHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "id")

	// If value persists in dictionary - redirect
	value, ok := storage.FindKey(key)
	if ok {
		w.Header().Set("Location", value)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return

		// Else error
	} else {
		http.Error(w, "No such url", 404)
		return
	}
}

// POST /
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	url := string(body)

	// If not exists - create new key-value pair
	index, ok := storage.FindVal(url)
	if !ok {
		index = storage.Add(url)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://localhost:8080/%s", index)
}
