package main

import (
	"net/http"

	"github.com/syberpunkq/go_url_shortener/internal/app/handlers"
)

func main() {
	http.HandleFunc("/", handlers.MyHandler)
	http.ListenAndServe(":8080", nil)
}
