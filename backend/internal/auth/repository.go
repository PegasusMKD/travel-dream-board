package auth

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	UpsertUser(ctx context.Context, params db.UpsertUserParams) (*db.User, error)
	GetUserById(ctx context.Context, uuid pgtype.UUID) (*db.User, error)
	CreateGuestUser(ctx context.Context, name string) (*db.User, error)
}

type authRepositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &authRepositoryImpl{
		queries: queries,
	}
}

func (r *authRepositoryImpl) UpsertUser(ctx context.Context, params db.UpsertUserParams) (*db.User, error) {
	user, err := r.queries.UpsertUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) GetUserById(ctx context.Context, uuid pgtype.UUID) (*db.User, error) {
	user, err := r.queries.GetUserById(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepositoryImpl) CreateGuestUser(ctx context.Context, name string) (*db.User, error) {
	user, err := r.queries.CreateGuestUser(ctx, name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
