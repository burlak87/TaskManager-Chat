package database

import (
	"context"
)

func (q *Queries) GetDB() DBTX {
	return q.db
}

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    task_id,
    title,
    message,
    type,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, user_id, task_id, title, message, is_read, type, created_at, expires_at
`

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRow(ctx, createNotification,
		arg.UserID,
		arg.TaskID,
		arg.Title,
		arg.Message,
		arg.Type,
		arg.ExpiresAt,
	)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TaskID,
		&i.Title,
		&i.Message,
		&i.IsRead,
		&i.Type,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getNotificationsByUserID = `-- name: GetNotificationsByUserID :many
SELECT id, user_id, task_id, title, message, is_read, type, created_at, expires_at FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) GetNotificationsByUserID(ctx context.Context, userID int64) ([]Notification, error) {
	rows, err := q.db.Query(ctx, getNotificationsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.TaskID,
			&i.Title,
			&i.Message,
			&i.IsRead,
			&i.Type,
			&i.CreatedAt,
			&i.ExpiresAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markNotificationAsRead = `-- name: MarkNotificationAsRead :exec
UPDATE notifications
SET is_read = TRUE
WHERE id = $1 AND user_id = $2
`

func (q *Queries) MarkNotificationAsRead(ctx context.Context, arg MarkNotificationAsReadParams) error {
	_, err := q.db.Exec(ctx, markNotificationAsRead, arg.ID, arg.UserID)
	return err
}

const getUnreadNotificationsCount = `-- name: GetUnreadNotificationsCount :one
SELECT COUNT(*) FROM notifications
WHERE user_id = $1 AND is_read = FALSE
`

func (q *Queries) GetUnreadNotificationsCount(ctx context.Context, arg GetUnreadNotificationsCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getUnreadNotificationsCount, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
