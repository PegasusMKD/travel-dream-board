package accomodations

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateAccomodation(ctx context.Context, data *Accomodation) (*Accomodation, error)
	GetAccomodationById(ctx context.Context, uuid string) (*Accomodation, error)
	GetAllAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error)
	UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error
	DeleteAccomodationById(ctx context.Context, uuid string) error
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

func (repo *repositoryImpl) CreateAccomodation(ctx context.Context, data *Accomodation) error {
	return nil
}

func (repo *repositoryImpl) GetAccomodationById(ctx context.Context, uuid string) (*Accomodation, error) {
	return nil, nil
}

func (repo *repositoryImpl) GetAllAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error) {
	return nil, nil
}

func (repo *repositoryImpl) UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error {
	return nil
}

func (repo *repositoryImpl) DeleteAccomodationById(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid)
		return err
	}

	return repo.queries.DeleteAccomodationById(ctx, id)
}
