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

func NewRouter() chi.Router {
	storage := storage.NewStorage()
	handler := NewHandler(storage)
	r := chi.NewRouter()
	r.Get("/{id}", handler.ShowHandler)
	r.Post("/", handler.CreateHandler)
	return r
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	require.NoError(t, err)

	client := http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "http://ya.ru")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/1", body)

	resp, body = testRequest(t, ts, "GET", "/1", "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Location"), "http://ya.ru")
	assert.Equal(t, "", body)

	resp, body = testRequest(t, ts, "GET", "/2", "")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "No such url\n", body)
}
