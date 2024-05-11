package internal

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer()

	t.Run("should return index response", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
	t.Run("should return callback response", func(t *testing.T) {
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api-2/complete", nil)

		server.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
