package websocket

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const maxMessageSize = 512 * 1024

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub    *Hub
	logger *logrus.Logger
}

func NewHandler(hub *Hub, logger *logrus.Logger) *Handler {
	return &Handler{
		hub:    hub,
		logger: logger,
	}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDInt64, ok := userID.(int64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	username := "User" + strconv.FormatInt(userIDInt64, 10)
	if u, exists := c.Get("username"); exists {
		if uStr, ok := u.(string); ok && uStr != "" {
			username = uStr
		}
	}

	boardIDStr := c.Query("board_id")
	if boardIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}

	conn.SetReadLimit(maxMessageSize)

	client := &Client{
		hub:      h.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userIDInt64,
		username: username,
		boardID:  boardID,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

