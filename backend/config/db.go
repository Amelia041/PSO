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

	uri := os.Getenv("mongodb://pso-17-ci-cd-server:18vPaEDSg89dAYMpG6AfiqjxfyARKod5gLhCdYn8laArcGJMVk9RqLXLo9Wr50bcGoLGgRYlTGsDACDb4hxlFg==@pso-17-ci-cd-server.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@pso-17-ci-cd-server@")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("studentsync")

	log.Println("Cosmos MongoDB Connected")
}
