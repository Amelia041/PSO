package main

import (
	"context"
	"log"
	"mahasiswa_app/backend/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	config.ConnectDB()
	r := gin.Default()

	// ── CORS Middleware ──────────────────────────────────────────
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204
			return
		}
		c.Next()
	})

	// ── POST /login ──────────────────────────────────────────────
	r.POST("/login", func(c *gin.Context) {
		var u struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		// ShouldBindJSON: kalau gagal parse, balas 400 — jangan lanjut
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
			return
		}
		if u.Username == "" || u.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var result bson.M
		err := config.DB.Collection("users").FindOne(ctx, bson.M{
			"username": u.Username,
			"password": u.Password, // TODO: ganti bcrypt.CompareHashAndPassword
		}).Decode(&result)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       result["_id"],
			"username": u.Username,
		})
	})

	// ── POST /notes ───────────────────────────────────────────────
	r.POST("/notes", func(c *gin.Context) {
		var n struct {
			UserID   string `json:"user_id"`
			Title    string `json:"title"`
			Content  string `json:"content"`
			Category string `json:"category"`
			Date     string `json:"date"`
		}
		if err := c.ShouldBindJSON(&n); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request tidak valid"})
			return
		}
		if n.Title == "" || n.UserID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id dan title wajib diisi"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := config.DB.Collection("notes").InsertOne(ctx, bson.M{
			"user_id":    n.UserID,
			"title":      n.Title,
			"content":    n.Content,
			"category":   n.Category,
			"created_at": n.Date,
		})
		if err != nil {
			log.Printf("InsertOne error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "saved"})
	})

	// ── GET /notes ────────────────────────────────────────────────
	r.GET("/notes", func(c *gin.Context) {
		uid := c.Query("user_id")
		cat := c.Query("category")

		if uid == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id wajib diisi"})
			return
		}

		// Filter pakai bson.M — TIDAK pakai string concat (aman dari injection)
		filter := bson.M{"user_id": uid}
		if cat != "" {
			filter["category"] = cat
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		opts := options.Find().SetSort(bson.M{"created_at": -1})
		cursor, err := config.DB.Collection("notes").Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
			return
		}
		defer cursor.Close(ctx)

		// Inisialisasi slice kosong — supaya response JSON [] bukan null
		notes := []gin.H{}
		for cursor.Next(ctx) {
			var doc bson.M
			if err := cursor.Decode(&doc); err != nil {
				continue
			}
			notes = append(notes, gin.H{
				"title":    doc["title"],
				"category": doc["category"],
				"date":     doc["created_at"],
			})
		}

		c.JSON(http.StatusOK, notes)
	})

	log.Println("Server berjalan di :8080")
	r.Run(":8080")
}
