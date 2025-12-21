package domain

import "time"

type Message struct {
	ID        string    `json:"id" bson:"_id"`
	BoardID   int64     `json:"board_id" bson:"board_id"`
	UserID    int64     `json:"user_id" bson:"user_id"`
	Username  string    `json:"username" bson:"username"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type MessageRequest struct {
	BoardID int64  `json:"board_id"`
	Content string `json:"content"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	BoardID   int64     `json:"board_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

