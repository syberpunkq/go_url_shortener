package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syberpunkq/go_url_shortener/internal/app/storage"
)

var baseURL = "http://localhost:8080"

func NewRouter() chi.Router {
	storage := storage.NewStorage()
	handler := NewHandler(storage, baseURL)
	r := chi.NewRouter()
	r.Get("/{id}", handler.ShowHandler)
	r.Post("/", handler.CreateHandler)
	r.Post("/api/shorten", handler.APICreateHandler)
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, json bool) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if json {
		req.Header.Set("Content-Type", "application/json")
	}
	require.NoError(t, err)

	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	statusCode := resp.StatusCode
	return statusCode, string(respBody)
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	statusCode, body := testRequest(t, ts, "POST", "/", "http://ya.ru", false)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "http://localhost:8080/1", body)

	statusCode, body = testRequest(t, ts, "GET", "/1", "", false)
	assert.Equal(t, http.StatusTemporaryRedirect, statusCode)
	// assert.Equal(t, resp.Header.Get("Location"), "http://ya.ru")
	assert.Equal(t, "", body)

	statusCode, body = testRequest(t, ts, "GET", "/2", "", false)
	assert.Equal(t, http.StatusNotFound, statusCode)
	assert.Equal(t, "No such url\n", body)

	statusCode, body = testRequest(t, ts, "POST", "/api/shorten", "{\"url\":\"http://ya.ru\"}", true)
	assert.Equal(t, http.StatusCreated, statusCode)
	assert.Equal(t, "{\"result\":\"http://localhost:8080/1\"}", body)

	statusCode, _ = testRequest(t, ts, "POST", "/api/shorten", "{\"url\":", true)
	assert.Equal(t, http.StatusBadRequest, statusCode)
}
