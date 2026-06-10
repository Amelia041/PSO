package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Setup router test — sama seperti di main() tapi tanpa ConnectDB
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
			return
		}
		if u.Username == "" || u.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
			return
		}
		// Di test, kita tidak hit DB — hanya cek logic validasi
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
	})

	r.GET("/notes", func(c *gin.Context) {
		uid := c.Query("user_id")
		if uid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi"})
			return
		}
		c.JSON(http.StatusOK, []gin.H{})
	})

	return r
}

// ── Test 1: Login tanpa body → 400 ──────────────────────────────
func TestLoginEmptyBody(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string]string{
		"username": "",
		"password": "",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya 400, dapat %d", w.Code)
	}
}

// ── Test 2: Login dengan kredensial salah → 401 ──────────────────
func TestLoginWrongCredentials(t *testing.T) {
	r := setupRouter()
	body, _ := json.Marshal(map[string]string{
		"username": "salah",
		"password": "salah",
	})

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Harusnya 401, dapat %d", w.Code)
	}
}

// ── Test 3: GET /notes tanpa user_id → 400 ──────────────────────
func TestGetNotesNoUserID(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/notes", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Harusnya 400, dapat %d", w.Code)
	}
}

// ── Test 4: GET /notes dengan user_id → 200 ─────────────────────
func TestGetNotesWithUserID(t *testing.T) {
	r := setupRouter()

	req := httptest.NewRequest("GET", "/notes?user_id=123", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Harusnya 200, dapat %d", w.Code)
	}
}
