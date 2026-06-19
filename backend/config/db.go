package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// GetMongoURI — dipisah supaya bisa dites tanpa koneksi
func GetMongoURI() string {
	return os.Getenv("MONGODB_URI")
}

// ConnectDB — koneksi ke CosmosDB
func ConnectDB() {
	uri := GetMongoURI()
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		cancel()
		log.Fatalf("Gagal connect ke MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		cancel()
		log.Fatalf("MongoDB tidak merespon: %v", err)
	}
	cancel()

	DB = client.Database("studentsync")
	log.Println("Cosmos MongoDB Connected")
}
