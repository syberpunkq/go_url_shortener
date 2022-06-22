package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

type Storaging interface {
	FindKey(key string) (string, bool)
	FindVal(val string) (string, bool)
	Add(val string) (string, error)
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
		Storaging: s,
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
		index, err = h.Storaging.Add(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", h.BaseURL, index)
}

// POST /api/shorten
func (h Handler) APICreateHandler(w http.ResponseWriter, r *http.Request) {
	//recieves body {"url":"<some_url>"}
	//returnes body {"result":"<shorten_url>"}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := data["url"].(string)

	// If not exists - create new key-value pair
	index, ok := h.Storaging.FindVal(url)
	if !ok {
		index, err = h.Storaging.Add(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	result := Data{Result: fmt.Sprintf("%s/%s", h.BaseURL, index)}
	json, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(json))
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (h Handler) GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// If request was compressed
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				io.WriteString(w, err.Error())
				next.ServeHTTP(w, r)
				return
			}
			defer reader.Close()
			r.Body = reader
		}

		// If doesnt accept gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
