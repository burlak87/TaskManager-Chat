package service

import (
	"context"
	"time"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/storage/psql"
	"github.com/your-team/taskmanager-chat/backend/pkg/logging"
)

type NotificationService struct {
	storage *psql.Storage
	logger  *logging.Logger
}

func NewNotificationService(storage *psql.Storage, logger *logging.Logger) *NotificationService {
	return &NotificationService{
		storage: storage,
		logger:  logger,
	}
}

func (s *NotificationService) Create(n domain.Notification) (domain.Notification, error) {
	return s.storage.CreateNotification(n)
}

func (s *NotificationService) GetByUser(userID int64) ([]domain.Notification, error) {
	return s.storage.GetNotifications(userID)
}

func (s *NotificationService) MarkAsRead(id, userID int64) error {
	return s.storage.MarkAsRead(id, userID)
}

func (s *NotificationService) GetUnreadCount(userID int64) (int64, error) {
	return s.storage.GetUnreadCount(userID)
}

func (s *NotificationService) StartDeadlineChecker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkDeadlines()
		}
	}
}

func (s *NotificationService) checkDeadlines() {
	tasks, err := s.storage.GetTasksWithUpcomingDeadlines(24 * time.Hour)
	if err != nil {
		s.logger.Errorf("Failed to check deadlines: %v", err)
		return
	}

	for _, task := range tasks {
		_, err := s.storage.CreateNotification(domain.Notification{
			UserID:  task.UserID,
			TaskID:  task.ID,
			Title:   "Task Deadline Approaching",
			Message: "Task '" + task.Title + "' is due soon.",
			Type:    "deadline",
		})
		if err != nil {
			continue
		}
		s.logger.Infof("Created deadline notification for task %d user %d", task.ID, task.UserID)
	}
}
