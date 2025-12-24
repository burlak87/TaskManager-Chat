package service

import (
	"context"
	"time"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
)

type ChatService struct {
	hub      *websocket.Hub
	storage  websocket.MessageStorage
}

func NewChatService(hub *websocket.Hub, storage websocket.MessageStorage) *ChatService {
	return &ChatService{
		hub:      hub,
		storage:  storage,
	}
}

func (s *ChatService) SaveMessage(ctx context.Context, message *domain.Message) (*domain.MessageResponse, error) {
	message.CreatedAt = time.Now()
	err := s.storage.SaveMessage(ctx, *message)
	if err != nil {
		return nil, err
	}

	response := &domain.MessageResponse{
		ID:        message.ID,
		BoardID:   message.BoardID,
		UserID:    message.UserID,
		Username:  message.Username,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	return response, nil
}

func (s *ChatService) GetMessages(ctx context.Context, boardID string, limit, offset int64) ([]*domain.MessageResponse, error) {
	return nil, nil
}

func (s *ChatService) GetMessageByID(ctx context.Context, id string) (*domain.MessageResponse, error) {
	return nil, nil
}

func (s *ChatService) CountMessages(ctx context.Context, boardID string) (int64, error) {
	return 0, nil
}

