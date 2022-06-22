package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

type Storaging interface {
	FindKey(key string) (string, bool)
	FindVal(val string) (string, bool)
	Add(val string) string
}

type Handler struct {
	BaseURL string
	Storaging
}

type Data struct {
	Result string `json:"result"`
}

func NewHandler(s *storage.Storage, URL string) *Handler {
	return &Handler{
		Storaging: storage.NewStorage(),
		BaseURL:   URL,
	}
}

// GET /{id}
func (h Handler) ShowHandler(w http.ResponseWriter, r *http.Request) {

	key := chi.URLParam(r, "id")
	// If value persists in dictionary - redirect
	value, ok := h.Storaging.FindKey(key)
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
func (h Handler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	url := string(body)

	// If not exists - create new key-value pair
	index, ok := h.Storaging.FindVal(url)
	if !ok {
		index = h.Storaging.Add(url)
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", h.BaseURL, index)
}

// POST /api/shorten
func (h Handler) ApiCreateHandler(w http.ResponseWriter, r *http.Request) {
	//recieves body {"url":"<some_url>"}
	//returnes body {"result":"<shorten_url>"}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	url := data["url"].(string)

	// If not exists - create new key-value pair
	index, ok := h.Storaging.FindVal(url)
	if !ok {
		index = h.Storaging.Add(url)
	}

	result := Data{Result: fmt.Sprintf("%s/%s", h.BaseURL, index)}
	json, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(json))
}
