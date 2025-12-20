package websocket

import (
	"log"
	"sync"

	"github.com/your-team/taskmanager-chat/backend/internal/models"
)

type Hub struct {
	rooms       map[string]map[*Client]bool
	register    chan *Client
	unregister  chan *Client
	broadcast   chan *BroadcastMessage
	messageChan chan *models.Message
	mu          sync.RWMutex
}

type BroadcastMessage struct {
	RoomID  string
	Message interface{}
	Exclude *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *BroadcastMessage, 256),
		messageChan: make(chan *models.Message, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToRoom(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[client.RoomID] == nil {
		h.rooms[client.RoomID] = make(map[*Client]bool)
	}
	h.rooms[client.RoomID][client] = true

	log.Printf("Client %s joined room %s. Total clients in room: %d", client.ID, client.RoomID, len(h.rooms[client.RoomID]))
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[client.RoomID]; ok {
		if _, exists := room[client]; exists {
			delete(room, client)
			close(client.send)
			log.Printf("Client %s left room %s. Remaining clients in room: %d", client.ID, client.RoomID, len(room))

			if len(room) == 0 {
				delete(h.rooms, client.RoomID)
			}
		}
	}
}

func (h *Hub) broadcastToRoom(message *BroadcastMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.rooms[message.RoomID]
	if !exists {
		return
	}

	for client := range room {
		if message.Exclude != nil && client == message.Exclude {
			continue
		}

		select {
		case client.send <- message.Message:
		default:
			close(client.send)
			delete(room, client)
		}
	}
}

func (h *Hub) GetRoomClientsCount(roomID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if room, exists := h.rooms[roomID]; exists {
		return len(room)
	}
	return 0
}

func (h *Hub) BroadcastMessage(roomID string, message interface{}, exclude *Client) {
	select {
	case h.broadcast <- &BroadcastMessage{
		RoomID:  roomID,
		Message: message,
		Exclude: exclude,
	}:
	default:
		log.Printf("Broadcast channel is full, dropping message for room %s", roomID)
	}
}

func (h *Hub) SendToRoom(roomID string, message *models.OutgoingMessage) {
	h.BroadcastMessage(roomID, message, nil)
}

func (h *Hub) GetMessageChan() <-chan *models.Message {
	return h.messageChan
}

