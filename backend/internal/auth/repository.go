package auth

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Repository interface {
	UpsertUser(ctx context.Context, params db.UpsertUserParams) (*db.User, error)
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
