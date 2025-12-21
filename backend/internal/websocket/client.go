package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/your-team/taskmanager-chat/backend/internal/models"
	"github.com/your-team/taskmanager-chat/backend/pkg/config"
)

const (
	pongWait  = 60 * time.Second
	pingPeriod = 54 * time.Second
	writeWait  = 10 * time.Second
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan interface{}
	ID       string
	Username string
	RoomID   string
	config   config.WebSocketConfig
}

func NewClient(hub *Hub, conn *websocket.Conn, userID, username, roomID string, cfg config.WebSocketConfig) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan interface{}, 256),
		ID:       userID,
		Username: username,
		RoomID:   roomID,
		config:   cfg,
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		c.handleMessage(messageBytes)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(messageBytes []byte) {
	var incomingMsg models.IncomingMessage
	if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		c.sendError("Invalid message format")
		return
	}

	switch incomingMsg.Type {
	case "message":
		c.handleChatMessage(incomingMsg.Payload)
	case "ping":
		c.handlePing(incomingMsg.Payload)
	case "join":
		c.handleJoin(incomingMsg.Payload)
	case "leave":
		c.handleLeave(incomingMsg.Payload)
	default:
		c.sendError("Unknown message type: " + incomingMsg.Type)
	}
}

func (c *Client) handleChatMessage(payload map[string]interface{}) {
	content, ok := payload["content"].(string)
	if !ok || content == "" {
		c.sendError("Content is required")
		return
	}

	boardID, ok := payload["board_id"].(string)
	if !ok {
		boardID = c.RoomID
	}

	var mentions []string
	if mentionsRaw, ok := payload["mentions"].([]interface{}); ok {
		for _, m := range mentionsRaw {
			if mention, ok := m.(string); ok {
				mentions = append(mentions, mention)
			}
		}
	}

	message := &models.Message{
		BoardID:  boardID,
		UserID:   c.ID,
		Username: c.Username,
		Content:  content,
		Mentions: mentions,
	}

	c.hub.messageChan <- message
}

func (c *Client) handlePing(payload map[string]interface{}) {
	timestamp := time.Now().Unix()
	pong := models.OutgoingMessage{
		Type: "pong",
		Payload: models.PongPayload{
			Timestamp: timestamp,
		},
	}
	c.send <- pong
}

func (c *Client) handleJoin(payload map[string]interface{}) {
	boardID, ok := payload["board_id"].(string)
	if !ok {
		c.sendError("board_id is required")
		return
	}

	if c.RoomID != "" && c.RoomID != boardID {
		c.hub.unregister <- c
	}

	c.RoomID = boardID
	c.hub.register <- c

	response := models.OutgoingMessage{
		Type: "joined",
		Payload: models.JoinRoomPayload{
			BoardID: boardID,
		},
	}
	c.send <- response
}

func (c *Client) handleLeave(payload map[string]interface{}) {
	boardID, ok := payload["board_id"].(string)
	if !ok {
		boardID = c.RoomID
	}

	if c.RoomID == boardID {
		c.hub.unregister <- c
		c.RoomID = ""
	}

	response := models.OutgoingMessage{
		Type: "left",
		Payload: models.LeaveRoomPayload{
			BoardID: boardID,
		},
	}
	c.send <- response
}

func (c *Client) sendError(message string) {
	errorMsg := models.OutgoingMessage{
		Type: "error",
		Payload: models.ErrorPayload{
			Message: message,
		},
	}
	select {
	case c.send <- errorMsg:
	default:
		log.Printf("Error channel is full, cannot send error to client %s", c.ID)
	}
}

func ServeWebSocket(hub *Hub, c *gin.Context, userID, username string, cfg config.WebSocketConfig) {
	boardID := c.Query("board_id")
	if boardID == "" {
		c.JSON(400, gin.H{"error": "board_id is required"})
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:    cfg.ReadBufferSize,
		WriteBufferSize:   cfg.WriteBufferSize,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: cfg.EnableCompression,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := NewClient(hub, conn, userID, username, boardID, cfg)
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

