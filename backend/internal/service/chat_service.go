package service

import (
	"context"

	"github.com/your-team/taskmanager-chat/backend/internal/models"
	"github.com/your-team/taskmanager-chat/backend/internal/repository"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
)

type ChatService struct {
	messageRepo *repository.MessageRepository
	hub         *websocket.Hub
}

func NewChatService(messageRepo *repository.MessageRepository, hub *websocket.Hub) *ChatService {
	return &ChatService{
		messageRepo: messageRepo,
		hub:         hub,
	}
}

func (s *ChatService) SaveMessage(ctx context.Context, message *models.Message) (*models.MessageResponse, error) {
	savedMessage, err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	outgoingMsg := &models.OutgoingMessage{
		Type:    "message",
		Payload: savedMessage.ToResponse(),
	}
	s.hub.SendToRoom(message.BoardID, outgoingMsg)

	return savedMessage.ToResponse(), nil
}

func (s *ChatService) GetMessages(ctx context.Context, boardID string, limit, offset int64) ([]*models.MessageResponse, error) {
	messages, err := s.messageRepo.GetByBoardID(ctx, boardID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*models.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = msg.ToResponse()
	}

	return responses, nil
}

func (s *ChatService) GetMessageByID(ctx context.Context, id string) (*models.MessageResponse, error) {
	message, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return message.ToResponse(), nil
}

func (s *ChatService) CountMessages(ctx context.Context, boardID string) (int64, error) {
	return s.messageRepo.CountByBoardID(ctx, boardID)
}

