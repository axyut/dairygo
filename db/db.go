package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongo() (*Mongo, error) {
	var dbName, uri string = "dairyDB", "mongodb://localhost:27017/"
	if err := godotenv.Load(); err != nil {
		log.Println("Set your 'MONGODB_URI' environment variable. " + "No .env file found\nUsing the default 'mongodb://localhost:27017'")
		uri = uri + dbName
	} else {
		dbName = os.Getenv("DB_NAME")
		uri = os.Getenv("MONGODB_URI")
		if uri == "" || dbName == "" {
			uri = "mongodb://localhost:27017/"
			dbName = "dairyDB"
		}
		uri = uri + dbName

	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// Check the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB! ", dbName)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	database := client.Database(dbName)
	return &Mongo{db: database}, nil
}

func (m *Mongo) GetDB() *mongo.Database {
	return m.db
}

func (m *Mongo) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
