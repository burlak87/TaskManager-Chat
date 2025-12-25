package psql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/storage/psql/sqlc"
)

type StorageError struct {
	Message string
}

type Storage struct {
	queries *database.Queries
}

func NewStorage(queries *database.Queries) *Storage {
	return &Storage{queries: queries}
}

var (
	ErrTokenExpired = &StorageError{"token expired"}
)

func (e *StorageError) Error() string {
	return e.Message
}

func (s *Storage) InsertUser(user domain.User) (int64, error) {
	params := database.CreateUserParams{
		Username:     user.Username,
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Email:     	  user.Email,
		PasswordHash: user.PasswordHash,
		TwoFaEnabled: pgtype.Bool{Bool: user.TwoFAEnabled, Valid: true},
	}

	createdUser, err := s.queries.CreateUser(context.Background(), params)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "users_username_key") || strings.Contains(errStr, "duplicate key value violates unique constraint") {
			if strings.Contains(errStr, "username") {
				return 0, errors.New("пользователь с таким именем уже существует")
			} else if strings.Contains(errStr, "email") {
				return 0, errors.New("пользователь с таким email уже существует")
			}
		}
		return 0, err
	}
	return createdUser.ID, nil
}

func (s *Storage) SelectUser(email string) (domain.User, error) {
	user, err := s.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return domain.User{}, err
	}
	
	return domain.User{
		ID: user.ID,
		Username: user.Username,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
		TwoFAEnabled: user.TwoFaEnabled.Bool,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *Storage) SelectUserByID(id int64) (domain.User, error) {
	user, err := s.queries.GetUserByID(context.Background(), id)
	if err != nil {
		return domain.User{}, nil
	}
	
	return domain.User{
		ID: user.ID,
		Username: user.Username,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
		Email: user.Email,
		PasswordHash: user.PasswordHash,
		TwoFAEnabled: user.TwoFaEnabled.Bool,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *Storage) RefreshStore(userID int64, token string, expiresAt time.Time) error {
	params := database.CreateRefreshTokenParams{
		UserID: userID,
		Token: token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	}
	
	return s.queries.CreateRefreshToken(context.Background(), params)
}

func (s *Storage) RefreshGet(token string) (int64, error) {
	refreshToken, err := s.queries.GetRefreshToken(context.Background(), token)
	if err != nil {
		return 0, err
	}
	
	if time.Now().After(refreshToken.ExpiresAt.Time) {
		s.RefreshDelete(token)
		return 0, ErrTokenExpired
	}
	
	return refreshToken.UserID, nil
}

func (s *Storage) RefreshDelete(token string) error {
	return s.queries.DeleteRefreshToken(context.Background(), token)
}

func (s *Storage) LogAttempt(email string, result bool, attemptTime time.Time) error {
	params := database.CreateLoginAttemptParams{
		Email: email,
		Success: result,
		AttemptedAt: pgtype.Timestamptz{Time: attemptTime, Valid: true},
	}
	
	return s.queries.CreateLoginAttempt(context.Background(), params)
}

func (s *Storage) GetFailedLogAttempts(email string, windowStart time.Time) (int, error) {
	count, err := s.queries.GetRecentFailedAttempts(context.Background(), database.GetRecentFailedAttemptsParams{
		Email:       email,
		AttemptedAt: pgtype.Timestamptz{Time: windowStart, Valid: true},
	})
	if err != nil {
		return 0, err
	}
	
	return int(count), nil
}

func (s *Storage) UserBlocked(email string, windowStart time.Time) ([]map[string]interface{}, error) {
	blockedUntil, err := s.queries.GetBlockedStatus(context.Background(), email)
	if err != nil {
		return []map[string]interface{}{}, nil
	}
	
	var isBlocked bool
	if blockedUntil.Valid && blockedUntil.Time.After(time.Now()) {
		isBlocked = true
	}
	
	var blockedUntilStr string
	if blockedUntil.Valid {
		blockedUntilStr = blockedUntil.Time.Format(time.RFC3339)
	}

	return []map[string]interface{}{
		{
			"blocked_until": blockedUntilStr,
			"is_blocked": isBlocked,
		},
	}, nil
}

func (s *Storage) BlockUser(email, blockedUntil string) error {
	var blockedUntilTime pgtype.Timestamptz
	if blockedUntil != "" {
		t, err := time.Parse(time.RFC3339, blockedUntil)
		if err != nil {
			return err
		}
		blockedUntilTime = pgtype.Timestamptz{Time: t, Valid: true}
	} else {
		blockedUntilTime = pgtype.Timestamptz{Valid: false}
	}
	
	return s.queries.BlockUser(context.Background(), database.BlockUserParams{
		Email:        email,
		BlockedUntil: blockedUntilTime,
	})
}

func (s *Storage) RenovationTwoFAStatus(userID int64, enabled bool) error {
	return s.queries.UpdateTwoFAStatus(context.Background(), database.UpdateTwoFAStatusParams{
		ID:           userID,
		TwoFaEnabled: pgtype.Bool{Bool: enabled, Valid: true},
	})
}

func (s *Storage) ResetFailedAttempts(email string) error {
	return s.queries.ResetFailedAttempts(context.Background(), email)
}

func (s *Storage) UpdatePasswordHash(userID int64, passwordHash string) error {
	return s.queries.UpdatePasswordHash(context.Background(), database.UpdatePasswordHashParams{
		ID:           userID,
		PasswordHash: passwordHash,
	})
}

func (s *Storage) DeleteExpiredRefreshTokens() error {
	return s.queries.DeleteExpiredRefreshTokens(context.Background())
}

func (s *Storage) RefreshDeleteByUserID(userID int64) error {
	return s.queries.RefreshDeleteByUserID(context.Background(), userID)
}

func (s *Storage) InsertTwoFaCode(userID int64, code string, expiresAt time.Time) error {
	_, err := s.queries.CreateTwoFaCode(context.Background(), database.CreateTwoFaCodeParams{
		UserID:    userID,
		Code:      code,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	return err
}

func (s *Storage) SelectTwoFaCodeByUserID(userID int64) (domain.TwoFaCode, error) {
	code, err := s.queries.GetTwoFaCodeByUserID(context.Background(), userID)
	if err != nil {
		return domain.TwoFaCode{}, err
	}
	
	return domain.TwoFaCode{
		ID:        code.ID,
		UserID:    code.UserID,
		Code:      code.Code,
		ExpiresAt: code.ExpiresAt.Time,
		Attempts:  int(code.Attempts.Int32),
		IsUsed:    code.IsUsed.Bool,
		CreatedAt: code.CreatedAt.Time,
	}, nil
}

func (s *Storage) RenovationTwoFaCodeAttempts(codeID int64, attempts int) error {
	return s.queries.UpdateTwoFaCodeAttempts(context.Background(), database.UpdateTwoFaCodeAttemptsParams{
		ID:       codeID,
		Attempts: pgtype.Int4{Int32: int32(attempts), Valid: true},
	})
}

func (s *Storage) MarkTwoFaCodeUsed(codeID int64) error {
	return s.queries.MarkTwoFaCodeAsUsed(context.Background(), codeID)
}

func (s *Storage) SelectRecentCodeRequests(userID int64, since time.Time) (int, error) {
	count, err := s.queries.GetRecentCodeRequests(context.Background(), database.GetRecentCodeRequestsParams{
		UserID:    userID,
		CreatedAt: pgtype.Timestamptz{Time: since, Valid: true},
	})
	if err != nil {
		return 0, err
	}
	
	return int(count), nil
}

func (s *Storage) SelectRecentVerificationAttempts(userID int64, since time.Time) (int, error) {
	count, err := s.queries.GetRecentVerificationAttempts(context.Background(), database.GetRecentVerificationAttemptsParams{
		UserID:    userID,
		CreatedAt: pgtype.Timestamptz{Time: since, Valid: true},
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *Storage) CreateBoard(ctx context.Context, title string, description string, ownerID int64) (domain.Board, error) {
	params := database.CreateBoardParams{
		Title:       title,
		Description: description,
		OwnerID:     ownerID,
	}

	board, err := s.queries.CreateBoard(ctx, params)
	if err != nil {
		return domain.Board{}, err
	}

	return domain.Board{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		CreatedAt:   board.CreatedAt.Time,
		UpdatedAt:   board.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetBoardByID(ctx context.Context, id int64) (domain.Board, error) {
	board, err := s.queries.GetBoardByID(ctx, id)
	if err != nil {
		return domain.Board{}, err
	}

	return domain.Board{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		CreatedAt:   board.CreatedAt.Time,
		UpdatedAt:   board.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetBoardsByOwner(ctx context.Context, ownerID int64) ([]domain.Board, error) {
	boards, err := s.queries.GetBoardsByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Board, len(boards))
	for i, b := range boards {
		result[i] = domain.Board{
			ID:          b.ID,
			Title:       b.Title,
			Description: b.Description,
			OwnerID:     b.OwnerID,
			CreatedAt:   b.CreatedAt.Time,
			UpdatedAt:   b.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (s *Storage) UpdateBoard(ctx context.Context, id int64, title *string, description *string) (domain.Board, error) {
	params := database.UpdateBoardParams{
		ID:          id,
		Title:       title,
		Description: description,
	}

	board, err := s.queries.UpdateBoard(ctx, params)
	if err != nil {
		return domain.Board{}, err
	}

	return domain.Board{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
		OwnerID:     board.OwnerID,
		CreatedAt:   board.CreatedAt.Time,
		UpdatedAt:   board.UpdatedAt.Time,
	}, nil
}

func (s *Storage) DeleteBoard(ctx context.Context, id int64) error {
	return s.queries.DeleteBoard(ctx, id)
}

func (s *Storage) CreateColumn(ctx context.Context, boardID int64, title string, position int) (domain.Column, error) {
	params := database.CreateColumnParams{
		BoardID:  boardID,
		Title:    title,
		Position: int32(position),
	}

	column, err := s.queries.CreateColumn(ctx, params)
	if err != nil {
		return domain.Column{}, err
	}

	return domain.Column{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Title:     column.Title,
		Position:  int(column.Position),
		CreatedAt: column.CreatedAt.Time,
		UpdatedAt: column.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetColumnByID(ctx context.Context, id int64) (domain.Column, error) {
	column, err := s.queries.GetColumnByID(ctx, id)
	if err != nil {
		return domain.Column{}, err
	}

	return domain.Column{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Title:     column.Title,
		Position:  int(column.Position),
		CreatedAt: column.CreatedAt.Time,
		UpdatedAt: column.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetColumnsByBoardID(ctx context.Context, boardID int64) ([]domain.Column, error) {
	columns, err := s.queries.GetColumnsByBoardID(ctx, boardID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Column, len(columns))
	for i, c := range columns {
		result[i] = domain.Column{
			ID:        c.ID,
			BoardID:   c.BoardID,
			Title:     c.Title,
			Position:  int(c.Position),
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (s *Storage) UpdateColumn(ctx context.Context, id int64, title *string, position *int) (domain.Column, error) {
	var position32 *int32
	if position != nil {
		pos := int32(*position)
		position32 = &pos
	}

	params := database.UpdateColumnParams{
		ID:       id,
		Title:    title,
		Position: position32,
	}

	column, err := s.queries.UpdateColumn(ctx, params)
	if err != nil {
		return domain.Column{}, err
	}

	return domain.Column{
		ID:        column.ID,
		BoardID:   column.BoardID,
		Title:     column.Title,
		Position:  int(column.Position),
		CreatedAt: column.CreatedAt.Time,
		UpdatedAt: column.UpdatedAt.Time,
	}, nil
}

func (s *Storage) DeleteColumn(ctx context.Context, id int64) error {
	return s.queries.DeleteColumn(ctx, id)
}

func (s *Storage) CreateTask(ctx context.Context, boardID int64, columnID *int64, title, description string, position *int) (domain.Task, error) {
	var position32 *int32
	if position != nil {
		pos := int32(*position)
		position32 = &pos
	}

	params := database.CreateTaskParams{
		BoardID:     boardID,
		ColumnID:    columnID,
		Title:       title,
		Description: description,
		Position:    position32,
	}

	task, err := s.queries.CreateTask(ctx, params)
	if err != nil {
		return domain.Task{}, err
	}

	return domain.Task{
		ID:          task.ID,
		BoardID:     task.BoardID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Position:    int(task.Position.Int32),
		CreatedAt:   task.CreatedAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetTaskByID(ctx context.Context, id int64) (domain.Task, error) {
	task, err := s.queries.GetTaskByID(ctx, id)
	if err != nil {
		return domain.Task{}, err
	}

	return domain.Task{
		ID:          task.ID,
		BoardID:     task.BoardID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Position:    int(task.Position.Int32),
		CreatedAt:   task.CreatedAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
	}, nil
}

func (s *Storage) GetTasksByBoardID(ctx context.Context, boardID int64) ([]domain.Task, error) {
	tasks, err := s.queries.GetTasksByBoardID(ctx, boardID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Task, len(tasks))
	for i, t := range tasks {
		result[i] = domain.Task{
			ID:          t.ID,
			BoardID:     t.BoardID,
			ColumnID:    t.ColumnID,
			Title:       t.Title,
			Description: t.Description,
			Position:    int(t.Position.Int32),
			CreatedAt:   t.CreatedAt.Time,
			UpdatedAt:   t.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id int64, title, description *string, columnID *int64, position *int) (domain.Task, error) {
	var position32 *int32
	if position != nil {
		pos := int32(*position)
		position32 = &pos
	}

	params := database.UpdateTaskParams{
		ID:          id,
		Title:       title,
		Description: description,
		ColumnID:    columnID,
		Position:    position32,
	}

	task, err := s.queries.UpdateTask(ctx, params)
	if err != nil {
		return domain.Task{}, err
	}

	return domain.Task{
		ID:          task.ID,
		BoardID:     task.BoardID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Position:    int(task.Position.Int32),
		CreatedAt:   task.CreatedAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
	}, nil
}

func (s *Storage) DeleteTask(ctx context.Context, id int64) error {
	return s.queries.DeleteTask(ctx, id)
}