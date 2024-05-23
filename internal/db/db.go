package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/axyut/dairygo/internal/config"
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

var Ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

func NewMongo(conf config.Config, logger *slog.Logger) (*Mongo, error) {
	var dbName, uri string = conf.DB_NAME, conf.DB_URI

	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	// Check the connection
	if err := client.Ping(Ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	logger.Info("Connected to MongoDB", slog.String("DB_NAME", dbName))
	database := client.Database(dbName)

	return &Mongo{Client: client, DB: database}, nil
}

func (m *Mongo) Close(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}
	cancel()
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
