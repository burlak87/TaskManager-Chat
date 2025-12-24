package websocket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   int64
	username string
	boardID  int64
}

type Hub struct {
	clients    map[*Client]bool
	rooms      map[int64]map[*Client]bool
	broadcast  chan domain.MessageResponse
	register   chan *Client
	unregister chan *Client
	storage    MessageStorage
	logger     *logrus.Logger
	mu         sync.RWMutex
}

type MessageStorage interface {
	SaveMessage(ctx context.Context, msg domain.Message) error
}

func NewHub(storage MessageStorage, logger *logrus.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[int64]map[*Client]bool),
		broadcast:  make(chan domain.MessageResponse, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		storage:    storage,
		logger:     logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.rooms[client.boardID] == nil {
				h.rooms[client.boardID] = make(map[*Client]bool)
			}
			h.rooms[client.boardID][client] = true
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Infof("Client registered: userID=%d, boardID=%d", client.userID, client.boardID)

		case client := <-h.unregister:
			h.mu.Lock()
			if room, ok := h.rooms[client.boardID]; ok {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.send)
					if len(room) == 0 {
						delete(h.rooms, client.boardID)
					}
				}
			}
			delete(h.clients, client)
			h.mu.Unlock()
			h.logger.Infof("Client unregistered: userID=%d, boardID=%d", client.userID, client.boardID)

		case message := <-h.broadcast:
			h.mu.RLock()
			room, exists := h.rooms[message.BoardID]
			if exists {
				for client := range room {
					select {
					case client.send <- h.messageToJSON(message):
					default:
						close(client.send)
						delete(room, client)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
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
		var msg domain.MessageRequest
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.logger.Errorf("WebSocket error: %v", err)
			}
			break
		}

		if msg.BoardID != c.boardID {
			continue
		}

		message := domain.Message{
			BoardID:  msg.BoardID,
			UserID:   c.userID,
			Username: c.username,
			Content:  msg.Content,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = c.hub.storage.SaveMessage(ctx, message)
		cancel()

		if err != nil {
			c.hub.logger.Errorf("Failed to save message: %v", err)
			continue
		}

		response := domain.MessageResponse{
			ID:        message.ID,
			BoardID:   message.BoardID,
			UserID:    message.UserID,
			Username:  message.Username,
			Content:   message.Content,
			CreatedAt: message.CreatedAt,
		}

		c.hub.broadcast <- response
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
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

func (h *Hub) messageToJSON(msg domain.MessageResponse) []byte {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		h.logger.Errorf("Failed to marshal message: %v", err)
		return []byte{}
	}
	return jsonData
}