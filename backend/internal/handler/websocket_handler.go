package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
)

type WebSocketHandler struct {
	handler *websocket.Handler
}

func NewWebSocketHandler(handler *websocket.Handler) *WebSocketHandler {
	return &WebSocketHandler{
		handler: handler,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	h.handler.HandleWebSocket(c)
}

