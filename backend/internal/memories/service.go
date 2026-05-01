package memories

import (
	"context"
	"errors"
	"os"
)

var ErrNotFound = errors.New("memory not found")

type Service interface {
	CreateMemory(ctx context.Context, boardUuid, userUuid, imageUrl string) (*Memory, error)
	GetMemoryByUuid(ctx context.Context, uuid string) (*Memory, error)
	GetMemoriesByBoardId(ctx context.Context, boardUuid string) ([]*Memory, error)
	DeleteMemoryByUuid(ctx context.Context, uuid string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo: repo}
}

func (s *serviceImpl) CreateMemory(ctx context.Context, boardUuid, userUuid, imageUrl string) (*Memory, error) {
	return s.repo.CreateMemory(ctx, boardUuid, userUuid, imageUrl)
}

func (s *serviceImpl) GetMemoryByUuid(ctx context.Context, uuid string) (*Memory, error) {
	return s.repo.GetMemoryByUuid(ctx, uuid)
}

func (s *serviceImpl) GetMemoriesByBoardId(ctx context.Context, boardUuid string) ([]*Memory, error) {
	return s.repo.FindAllByBoardUuid(ctx, boardUuid)
}

func (s *serviceImpl) DeleteMemoryByUuid(ctx context.Context, uuid string) error {
	memory, err := s.repo.GetMemoryByUuid(ctx, uuid)
	if err != nil {
		return err
	}
	if memory == nil {
		return ErrNotFound
	}
	if err := s.repo.DeleteMemoryByUuid(ctx, uuid); err != nil {
		return err
	}
	// Best-effort cleanup of the file on disk.
	_ = os.Remove(memory.ImageUrl)
	return nil
}
