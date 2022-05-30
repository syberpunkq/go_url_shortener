package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

var dict = make(map[string]string)
var currentIndex = 1

func handler(w http.ResponseWriter, r *http.Request) {

	// GET /{id}
	if r.Method == http.MethodGet {
		key := strings.Trim(r.URL.Path[1:], "/")

		// If value persists in dictionary - redirect
		if value, ok := dict[key]; ok {
			w.Header().Set("Location", value)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
			// Else 404
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

		// Check if url already in dict
		exists := false
		index := "0"
		for k, v := range dict {
			if v == url {
				exists = true
				index = k
				break
			}
		}

		// If not exists - create new key-value pair
		if !exists {
			index = fmt.Sprint(currentIndex)
			dict[index] = url

			currentIndex++
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://localhost:8080/%s", index)
	} else {
		http.Error(w, "Bad Request", 400)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
