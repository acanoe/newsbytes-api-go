package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestRegisterRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	t.Run("Register user OK", func(t *testing.T) {
		payload := strings.NewReader(`{"username": "johndoe", "password": "password"}`)

		req, _ := http.NewRequest("POST", "/auth/register", payload)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// assert.Equal(t, 200, w.Code)
		// assert.Equal(t, `{"message":"user registered"}`, w.Body.String())
	})

}

func TestLoginRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	t.Run("Login user OK", func(t *testing.T) {
		payload := strings.NewReader(`{"username": "johndoe", "password": "password"}`)

		req, _ := http.NewRequest("POST", "/auth/login", payload)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// assert.Equal(t, 200, w.Code)
	})

}
