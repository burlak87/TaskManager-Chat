package repository

import (
	"context"
	"time"

	"github.com/your-team/taskmanager-chat/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *MessageRepository) GetByBoardID(ctx context.Context, boardID string, limit, offset int64) ([]*models.Message, error) {
	filter := bson.M{"board_id": boardID}
	sort := bson.D{{Key: "created_at", Value: -1}}
	opts := options.Find().
		SetSort(sort).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id string) (*models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var message models.Message
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessageRepository) CountByBoardID(ctx context.Context, boardID string) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"board_id": boardID})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *MessageRepository) DeleteByBoardID(ctx context.Context, boardID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"board_id": boardID})
	return err
}
