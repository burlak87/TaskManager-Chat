package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	TwoFAEnabled bool      `json:"two_fa_enabled"`
}

type Board struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OwnerID     int64     `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BoardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Column struct {
	ID        int64     `json:"id"`
	BoardID   int64     `json:"board_id"`
	Title     string    `json:"title"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ColumnRequest struct {
	Title    string `json:"title"`
	Position int    `json:"position"`
}

type Task struct {
	ID          int64     `json:"id"`
	BoardID     int64     `json:"board_id"`
	ColumnID    int64     `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ColumnID    *int64 `json:"column_id,omitempty"`
	Position    *int   `json:"position,omitempty"`
}

type TaskUpdateRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	ColumnID    *int64  `json:"column_id,omitempty"`
	Position    *int    `json:"position,omitempty"`
}

type TwoFaCodes struct {
	RequiresTwoFa bool   `json:"requires_two_fa"`
	TempToken     string `json:"temp_token"`
}

type TwoFaCode struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
	Attempts  int       `json:"attempts"`
	IsUsed    bool      `json:"is_used"`
	CreatedAt time.Time `json:"created_at"`
}

type Code struct {
	TempToken string `json:"temp_token"`
	Code      string `json:"code"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TwoFaToggleRequest struct {
	Password string `json:"password"`
}