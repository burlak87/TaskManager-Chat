package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/websocket"
	"github.com/your-team/taskmanager-chat/backend/pkg/config"
)

type WebSocketHandler struct {
	hub *websocket.Hub
	cfg config.WebSocketConfig
}

func NewWebSocketHandler(hub *websocket.Hub, cfg config.WebSocketConfig) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
		cfg: cfg,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "user not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "username not found"})
		return
	}

	websocket.ServeWebSocket(h.hub, c, userID.(string), username.(string), h.cfg)
}

