package psql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	database "github.com/your-team/taskmanager-chat/backend/internal/storage/psql/sqlc"
)

func (s *Storage) CreateNotification(n domain.Notification) (domain.Notification, error) {
	var taskID pgtype.Int8
	if n.TaskID != 0 {
		taskID = pgtype.Int8{Int64: n.TaskID, Valid: true}
	} else {
		taskID = pgtype.Int8{Valid: false}
	}

	var expiresAt pgtype.Timestamptz
	if !n.ExpiresAt.IsZero() {
		expiresAt = pgtype.Timestamptz{Time: n.ExpiresAt, Valid: true}
	} else {
		expiresAt = pgtype.Timestamptz{Valid: false}
	}

	res, err := s.queries.CreateNotification(context.Background(), database.CreateNotificationParams{
		UserID:    n.UserID,
		TaskID:    taskID,
		Title:     n.Title,
		Message:   n.Message,
		Type:      n.Type,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return domain.Notification{}, err
	}

	return domain.Notification{
		ID:        res.ID,
		UserID:    res.UserID,
		TaskID:    res.TaskID.Int64,
		Title:     res.Title,
		Message:   res.Message,
		IsRead:    res.IsRead.Bool,
		Type:      res.Type,
		CreatedAt: res.CreatedAt.Time,
		ExpiresAt: res.ExpiresAt.Time,
	}, nil
}

func (s *Storage) GetNotifications(userID int64) ([]domain.Notification, error) {
	rows, err := s.queries.GetNotificationsByUserID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	var notifications []domain.Notification
	for _, row := range rows {
		notifications = append(notifications, domain.Notification{
			ID:        row.ID,
			UserID:    row.UserID,
			TaskID:    row.TaskID.Int64,
			Title:     row.Title,
			Message:   row.Message,
			IsRead:    row.IsRead.Bool,
			Type:      row.Type,
			CreatedAt: row.CreatedAt.Time,
			ExpiresAt: row.ExpiresAt.Time,
		})
	}
	return notifications, nil
}

func (s *Storage) MarkAsRead(id, userID int64) error {
	return s.queries.MarkNotificationAsRead(context.Background(), database.MarkNotificationAsReadParams{
		ID:     id,
		UserID: userID,
	})
}

func (s *Storage) GetUnreadCount(userID int64) (int64, error) {
	return s.queries.GetUnreadNotificationsCount(context.Background(), database.GetUnreadNotificationsCountParams{
		UserID: userID,
	})
}

func (s *Storage) GetTasksWithUpcomingDeadlines(window time.Duration) ([]domain.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, deadline, created_at, updated_at
		FROM tasks
		WHERE deadline BETWEEN NOW() AND NOW() + $1
	`
	rows, err := s.queries.GetDB().Query(context.Background(), query, window)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.Deadline,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
