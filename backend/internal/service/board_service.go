package service

import (
	"context"

	"github.com/your-team/taskmanager-chat/backend/internal/domain"
)

type BoardStorage interface {
	CreateBoard(ctx context.Context, title string, description string, ownerID int64) (domain.Board, error)
	GetBoardByID(ctx context.Context, id int64) (domain.Board, error)
	GetBoardsByOwner(ctx context.Context, ownerID int64) ([]domain.Board, error)
	UpdateBoard(ctx context.Context, id int64, title, description *string) (domain.Board, error)
	DeleteBoard(ctx context.Context, id int64) error

	CreateColumn(ctx context.Context, boardID int64, title string, position int) (domain.Column, error)
	GetColumnByID(ctx context.Context, id int64) (domain.Column, error)
	GetColumnsByBoardID(ctx context.Context, boardID int64) ([]domain.Column, error)
	UpdateColumn(ctx context.Context, id int64, title *string, position *int) (domain.Column, error)
	DeleteColumn(ctx context.Context, id int64) error

	CreateTask(ctx context.Context, boardID int64, columnID *int64, title, description string, position *int) (domain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (domain.Task, error)
	GetTasksByBoardID(ctx context.Context, boardID int64) ([]domain.Task, error)
	UpdateTask(ctx context.Context, id int64, title, description *string, columnID *int64, position *int) (domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
}

type BoardService struct {
	storage BoardStorage
}

func NewBoardService(storage BoardStorage) *BoardService {
	return &BoardService{
		storage: storage,
	}
}

func (s *BoardService) CreateBoard(ctx context.Context, title, description string, ownerID int64) (domain.Board, error) {
	return s.storage.CreateBoard(ctx, title, description, ownerID)
}

func (s *BoardService) GetBoardByID(ctx context.Context, id int64) (domain.Board, error) {
	return s.storage.GetBoardByID(ctx, id)
}

func (s *BoardService) GetBoardsByOwner(ctx context.Context, ownerID int64) ([]domain.Board, error) {
	return s.storage.GetBoardsByOwner(ctx, ownerID)
}

func (s *BoardService) UpdateBoard(ctx context.Context, id int64, title, description string) (domain.Board, error) {
	titlePtr := &title
	if title == "" {
		titlePtr = nil
	}
	descriptionPtr := &description
	if description == "" {
		descriptionPtr = nil
	}
	
	return s.storage.UpdateBoard(ctx, id, titlePtr, descriptionPtr)
}

func (s *BoardService) DeleteBoard(ctx context.Context, id int64) error {
	return s.storage.DeleteBoard(ctx, id)
}

func (s *BoardService) CreateColumn(ctx context.Context, boardID int64, title string, position int) (domain.Column, error) {
	return s.storage.CreateColumn(ctx, boardID, title, position)
}

func (s *BoardService) GetColumnsByBoardID(ctx context.Context, boardID int64) ([]domain.Column, error) {
	return s.storage.GetColumnsByBoardID(ctx, boardID)
}

func (s *BoardService) GetColumnByID(ctx context.Context, id int64) (domain.Column, error) {
	return s.storage.GetColumnByID(ctx, id)
}

func (s *BoardService) UpdateColumn(ctx context.Context, id int64, title string, position int) (domain.Column, error) {
	titlePtr := &title
	if title == "" {
		titlePtr = nil
	}
	positionPtr := &position
	if position < 0 {
		positionPtr = nil
	}
	
	return s.storage.UpdateColumn(ctx, id, titlePtr, positionPtr)
}

func (s *BoardService) DeleteColumn(ctx context.Context, id int64) error {
	return s.storage.DeleteColumn(ctx, id)
}

func (s *BoardService) CreateTask(ctx context.Context, boardID int64, columnID *int64, title, description string, position *int) (domain.Task, error) {
	return s.storage.CreateTask(ctx, boardID, columnID, title, description, position)
}

func (s *BoardService) GetTasksByBoardID(ctx context.Context, boardID int64) ([]domain.Task, error) {
	return s.storage.GetTasksByBoardID(ctx, boardID)
}

func (s *BoardService) GetTaskByID(ctx context.Context, id int64) (domain.Task, error) {
	return s.storage.GetTaskByID(ctx, id)
}

func (s *BoardService) UpdateTask(ctx context.Context, id int64, title, description *string, columnID *int64, position *int) (domain.Task, error) {
	return s.storage.UpdateTask(ctx, id, title, description, columnID, position)
}

func (s *BoardService) DeleteTask(ctx context.Context, id int64) error {
	return s.storage.DeleteTask(ctx, id)
}