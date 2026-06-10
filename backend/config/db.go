package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	// v2: Connect tidak butuh ctx sebagai parameter
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Gagal connect ke MongoDB: %v", err)
	}

	// Ping tetap pakai ctx untuk timeout check
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB tidak merespon: %v", err)
	}

	DB = client.Database("studentsync")
	log.Println("Cosmos MongoDB Connected")
}
