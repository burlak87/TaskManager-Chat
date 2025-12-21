package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

type Config struct {
	URI      string
	Database string
	Username string
	Password string
}

func NewMongoDB(ctx context.Context, cfg Config) (*MongoDB, error) {
	uri := cfg.URI
	if cfg.Username != "" && cfg.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s", cfg.Username, cfg.Password, cfg.URI)
	} else if cfg.URI == "" {
		uri = "mongodb://localhost:27017"
	}

	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.Database)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

