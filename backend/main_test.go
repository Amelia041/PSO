package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLoginEmptyBody(t *testing.T) {
	r := SetupRouter() // ← pakai fungsi asli dari main.go
	body, _ := json.Marshal(map[string]string{"username": "", "password": ""})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya 400, dapat %d", w.Code)
	}
}

func TestLoginWrongCredentials(t *testing.T) {
	r := SetupRouter()
	body, _ := json.Marshal(map[string]string{"username": "salah", "password": "salah"})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// DB nil → 503, DB ada tapi salah kredensial → 401
	if w.Code != http.StatusServiceUnavailable && w.Code != http.StatusUnauthorized {
		t.Errorf("Harusnya 503 atau 401, dapat %d", w.Code)
	}
}

func TestGetNotesWithUserID(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("GET", "/notes?user_id=123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// DB nil → 503, DB ada → 200
	if w.Code != http.StatusServiceUnavailable && w.Code != http.StatusOK {
		t.Errorf("Harusnya 503 atau 200, dapat %d", w.Code)
	}
}

func TestGetNotesWithCategory(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("GET", "/notes?user_id=123&category=PSO", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable && w.Code != http.StatusOK {
		t.Errorf("Harusnya 503 atau 200, dapat %d", w.Code)
	}
}

func TestPostNotesValid(t *testing.T) {
	r := SetupRouter()
	body, _ := json.Marshal(map[string]string{
		"user_id":  "123",
		"title":    "Catatan Test",
		"category": "PSO",
		"content":  "Isi catatan",
		"date":     "2026-06-10",
	})

	req := httptest.NewRequest("POST", "/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// DB nil → 503, DB ada → 200
	if w.Code != http.StatusServiceUnavailable && w.Code != http.StatusOK {
		t.Errorf("Harusnya 503 atau 200, dapat %d", w.Code)
	}
}

func TestPostNotesMissingTitle(t *testing.T) {
	r := SetupRouter()
	body, _ := json.Marshal(map[string]string{
		"user_id": "123",
		"title":   "", // kosong
	})

	req := httptest.NewRequest("POST", "/notes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya 400, dapat %d", w.Code)
	}
}
