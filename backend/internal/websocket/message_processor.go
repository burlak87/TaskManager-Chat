package websocket

import (
	"context"
	"log"

	"github.com/your-team/taskmanager-chat/backend/internal/models"
	"github.com/your-team/taskmanager-chat/backend/internal/repository"
)

type MessageProcessor struct {
	hub          *Hub
	messageRepo  *repository.MessageRepository
	saveCallback func(*models.Message) (*models.MessageResponse, error)
}

func NewMessageProcessor(hub *Hub, messageRepo *repository.MessageRepository) *MessageProcessor {
	return &MessageProcessor{
		hub:         hub,
		messageRepo: messageRepo,
	}
}

func (mp *MessageProcessor) SetSaveCallback(callback func(*models.Message) (*models.MessageResponse, error)) {
	mp.saveCallback = callback
}

func (mp *MessageProcessor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-mp.hub.GetMessageChan():
			mp.processMessage(ctx, message)
		}
	}
}

func (mp *MessageProcessor) processMessage(ctx context.Context, message *models.Message) {
	var response *models.MessageResponse
	var err error

	if mp.saveCallback != nil {
		response, err = mp.saveCallback(message)
	} else {
		savedMessage, saveErr := mp.messageRepo.Create(ctx, message)
		if saveErr != nil {
			log.Printf("Error saving message: %v", saveErr)
			return
		}
		response = savedMessage.ToResponse()
	}

	if err != nil {
		log.Printf("Error processing message: %v", err)
		return
	}

	outgoingMsg := &models.OutgoingMessage{
		Type:    "message",
		Payload: response,
	}
	mp.hub.SendToRoom(message.BoardID, outgoingMsg)
}

