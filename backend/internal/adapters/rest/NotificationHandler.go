package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/service"
	"github.com/your-team/taskmanager-chat/backend/pkg/logging"
)

type NotificationHandler struct {
	service *service.NotificationService
	logger  *logging.Logger
}

func NewNotificationHandler(service *service.NotificationService, logger *logging.Logger) *NotificationHandler {
	return &NotificationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *NotificationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/notifications", h.GetNotifications)
	rg.POST("/notifications/:id/read", h.MarkAsRead)
	rg.GET("/notifications/unread-count", h.GetUnreadCount)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := userID.(int64)
	if !ok {
		if uidFloat, ok := userID.(float64); ok {
			uid = int64(uidFloat)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user id type mismatch"})
			return
		}
	}

	notifications, err := h.service.GetByUser(uid)
	if err != nil {
		h.logger.Errorf("Failed to get notifications for user %d: %v", uid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notifications"})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, _ := userID.(int64)

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.service.MarkAsRead(id, uid); err != nil {
		h.logger.Errorf("Failed to mark notification %d as read: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, _ := userID.(int64)

	count, err := h.service.GetUnreadCount(uid)
	if err != nil {
		h.logger.Errorf("Failed to get unread count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
