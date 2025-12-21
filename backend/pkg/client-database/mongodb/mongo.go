package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/your-team/taskmanager-chat/backend/pkg/config"
	repeatable "github.com/your-team/taskmanager-chat/backend/pkg/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, maxAttempts int, cfg config.MongoConfig) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	
	var client *mongo.Client
	var err error
	
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		
		clientOptions := options.Client().ApplyURI(uri)
		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			return err
		}
		
		err = client.Ping(ctx, nil)
		if err != nil {
			return err
		}
		
		return nil
	}, maxAttempts, 5*time.Second)
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	
	return client, nil
}

