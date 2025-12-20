package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BoardID   string             `bson:"board_id" json:"board_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Content   string             `bson:"content" json:"content"`
	Mentions  []string           `bson:"mentions,omitempty" json:"mentions,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type MessageRequest struct {
	BoardID  string   `json:"board_id" binding:"required"`
	Content  string   `json:"content" binding:"required"`
	Mentions []string `json:"mentions,omitempty"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Mentions  []string  `json:"mentions,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) ToResponse() *MessageResponse {
	return &MessageResponse{
		ID:        m.ID.Hex(),
		BoardID:   m.BoardID,
		UserID:    m.UserID,
		Username:  m.Username,
		Content:   m.Content,
		Mentions:  m.Mentions,
		CreatedAt: m.CreatedAt,
	}
}

