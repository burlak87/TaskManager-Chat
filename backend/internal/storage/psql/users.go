package psql

import (
	"context"
	"time"
	
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/storage/psql/sqlc"
)

type Storage struct {
	queries *sqlc.Queries
}

func NewStorage(queries *sqlc.Queries) *Storage {
	return &Storage{queries: queries}
}

func (s *Storage) InsertUsers(user domain.User) (int64, error) {
	params := sqlc.CreateUserParams{
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Email:     	  user.Email,
		PasswordHash: user.PasswordHash,
	}
	
	user, err := s.queries.CreateUser(context.Background(), params)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (s *Storage) SelectUsers(email string) (domain.User, error) {
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
		TwoFAEnabled: user.TwoFAEnabled,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Storage) SelectUsersByID(id int64) (domain.User, error) {
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
		TwoFAEnabled: user.TwoFAEnabled,
		CreatedAt: user.CreatedAt,
	}, nil
}

