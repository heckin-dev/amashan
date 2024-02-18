package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck_GetHealthcheck(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   int
	}{
		{
			name:   "GET /health",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "POST /health",
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/health", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(NewHealthcheck().GetHealthcheck)

			handler.ServeHTTP(rr, req)

			got := rr.Code
			want := http.StatusOK

			assert.Truef(t, got == want, "got status %v, wanted %v", got, want)
		})
	}
}
