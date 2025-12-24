package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	if boardIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	messages, err := h.chatService.GetMessages(c.Request.Context(), boardIDStr, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) GetMessageByID(c *gin.Context) {
	messageID := c.Param("message_id")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message_id is required"})
		return
	}

	message, err := h.chatService.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *ChatHandler) CreateMessage(c *gin.Context) {
	var req domain.MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found"})
		return
	}

	userIDInt64 := userID.(int64)

	message := &domain.Message{
		BoardID:  req.BoardID,
		UserID:   userIDInt64,
		Username: username.(string),
		Content:  req.Content,
	}

	response, err := h.chatService.SaveMessage(c.Request.Context(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ChatHandler) GetMessagesCount(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	if boardIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "board_id is required"})
		return
	}

	count, err := h.chatService.CountMessages(c.Request.Context(), boardIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

