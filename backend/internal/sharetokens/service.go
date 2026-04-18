package sharetokens

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

type Service interface {
	CreateShareToken(ctx context.Context, boardUuid string) (*ShareToken, error)
	GetShareToken(ctx context.Context, token string) (*ShareToken, error)
	DeleteShareToken(ctx context.Context, token string, boardUuid string) error
	GetShareTokensForBoard(ctx context.Context, boardUuid string) ([]*ShareToken, error)
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func generateRandomToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (svc *serviceImpl) CreateShareToken(ctx context.Context, boardUuid string) (*ShareToken, error) {
	token, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	return svc.repo.CreateShareToken(ctx, token, boardUuid)
}

func (svc *serviceImpl) GetShareToken(ctx context.Context, token string) (*ShareToken, error) {
	return svc.repo.GetShareToken(ctx, token)
}

func (svc *serviceImpl) DeleteShareToken(ctx context.Context, token string, boardUuid string) error {
	return svc.repo.DeleteShareToken(ctx, token, boardUuid)
}

func (svc *serviceImpl) GetShareTokensForBoard(ctx context.Context, boardUuid string) ([]*ShareToken, error) {
	return svc.repo.GetShareTokensForBoard(ctx, boardUuid)
}
