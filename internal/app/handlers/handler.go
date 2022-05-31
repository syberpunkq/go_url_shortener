package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

func MyHandler(w http.ResponseWriter, r *http.Request) {

	// GET /{id}
	if r.Method == http.MethodGet {
		key := strings.Trim(r.URL.Path[1:], "/")

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

		// POST /
	} else if r.Method == http.MethodPost {
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
	} else {
		http.Error(w, "Bad Request", 400)
		return
	}
}
