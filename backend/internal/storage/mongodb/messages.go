package mongodb

import (
	"context"
	"time"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageStorage struct {
	collection *mongo.Collection
}

func NewMessageStorage(client *mongo.Client, dbName string) *MessageStorage {
	return &MessageStorage{
		collection: client.Database(dbName).Collection("messages"),
	}
}

func (s *MessageStorage) SaveMessage(ctx context.Context, msg domain.Message) error {
	msg.ID = primitive.NewObjectID().Hex()
	msg.CreatedAt = time.Now()
	
	_, err := s.collection.InsertOne(ctx, msg)
	return err
}

func (s *MessageStorage) GetMessagesByBoardID(ctx context.Context, boardID int64, limit int64) ([]domain.Message, error) {
	filter := bson.M{"board_id": boardID}
	opts := options.Find().SetSort(bson.D{bson.E{Key: "created_at", Value: -1}}).SetLimit(limit)
	
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var messages []domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, nil
}

