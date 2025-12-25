package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Notification struct {
	ID        int64              `json:"id"`
	UserID    int64              `json:"user_id"`
	TaskID    pgtype.Int8        `json:"task_id"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	IsRead    pgtype.Bool        `json:"is_read"`
	Type      string             `json:"type"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

type CreateNotificationParams struct {
	UserID    int64              `json:"user_id"`
	TaskID    pgtype.Int8        `json:"task_id"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	Type      string             `json:"type"`
	ExpiresAt pgtype.Timestamptz `json:"expires_at"`
}

type MarkNotificationAsReadParams struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

type GetUnreadNotificationsCountParams struct {
	UserID int64 `json:"user_id"`
}
