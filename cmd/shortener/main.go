package main

import (
	"net/http"

	"github.com/syberpunkq/go_url_shortener/internal/app/router"
)

func main() {
	appRouter := router.New()
	http.ListenAndServe(":8080", appRouter)
}
