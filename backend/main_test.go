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

	// Tanpa DB akan 500, tapi bukan 400
	if w.Code == http.StatusBadRequest {
		t.Errorf("Tidak seharusnya 400 untuk input valid")
	}
}

func TestGetNotesNoUserID(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("GET", "/notes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya 400, dapat %d", w.Code)
	}
}

func TestGetNotesWithUserID(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("GET", "/notes?user_id=123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Tanpa DB akan 500, tapi validasi user_id sudah lewat (bukan 400)
	if w.Code == http.StatusBadRequest {
		t.Errorf("Tidak seharusnya 400 untuk user_id yang valid")
	}
}

func TestGetNotesWithCategory(t *testing.T) {
	r := SetupRouter()

	req := httptest.NewRequest("GET", "/notes?user_id=123&category=PSO", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code == http.StatusBadRequest {
		t.Errorf("Tidak seharusnya 400 untuk input valid")
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

	if w.Code == http.StatusBadRequest {
		t.Errorf("Tidak seharusnya 400 untuk input valid")
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
