package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyHandler(t *testing.T) {
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		method  string
		request string
		want    want
	}{
		{
			name:    "Post create short url",
			method:  "POST",
			request: "/",
			want:    want{statusCode: 201, body: "http://localhost:8080/1"},
		},
		{
			name:    "Get recieve long url by short",
			method:  "GET",
			request: "http://localhost:8080/1",
			want:    want{statusCode: 307, body: ""},
		},
		{
			name:    "Get invalid short url",
			method:  "GET",
			request: "http://localhost:8080/2",
			want:    want{statusCode: 404, body: "No such url"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			method := http.MethodGet
			if tt.method == "POST" {
				method = http.MethodPost
			}
			request := httptest.NewRequest(method, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(MyHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}