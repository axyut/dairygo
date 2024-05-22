package db

import (
	"context"
	"fmt"
	"log"
	"os"

	// "time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}
type Collection *mongo.Collection
type Collections struct {
	Users        Collection
	Goods        Collection
	Audiences    Collection
	Transactions Collection
}

var Ctx = context.TODO()

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

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// Check the connection
	if err := client.Ping(Ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB! ", dbName)
	defer func() {

	}()
	database := client.Database(dbName)

	return &Mongo{Client: client, DB: database}, nil
}

func (m *Mongo) Close(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	// cancel()
	return nil
}

func GetCollections(ctx context.Context, m *mongo.Database) *Collections {
	return &Collections{
		Users:        m.Collection("users"),
		Goods:        m.Collection("goods"),
		Audiences:    m.Collection("audiences"),
		Transactions: m.Collection("transactions"),
	}
}

func GetUserColl(ctx context.Context, m *mongo.Database) Collection {
	return GetCollections(ctx, m).Users
}

func GetAudienceColl(ctx context.Context, m *mongo.Database) Collection {
	return GetCollections(ctx, m).Audiences
}

func GetGoodsColl(ctx context.Context, m *mongo.Database) Collection {
	return GetCollections(ctx, m).Goods
}

func GetTransactionColl(ctx context.Context, m *mongo.Database) Collection {
	return GetCollections(ctx, m).Transactions
}
