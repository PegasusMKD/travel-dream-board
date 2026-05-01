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

func (repo *repositoryImpl) CreateAccomodation(ctx context.Context, data *Accomodation) (*Accomodation, error) {
	boardUuid, err := utility.UuidFromString(data.BoardUuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", data.BoardUuid, "error", err)
		return nil, err
	}

	userUuid, err := utility.UuidFromString(data.UserUuid)
	if err != nil {
		log.Error("Failed parsing UserUUID", "uuid", data.UserUuid, "error", err)
		return nil, err
	}

	params := db.CreateAccomodationParams{
		Url:         data.Url,
		Title:       data.Title,
		ImageUrl:    data.ImageUrl,
		BoardUuid:   boardUuid,
		UserUuid:    userUuid,
		Price:       utility.NumericFromString(data.Price),
		Currency:    utility.NullCurrencyFromPtr(data.Currency),
		Description: data.Description,
	}

	entity, err := repo.queries.CreateAccomodation(ctx, params)
	if err != nil {
		log.Error("Failed creating entity", "url", data.Url, "error", err)
		return nil, err
	}

	return FromEntity(entity), nil
}

func (repo *repositoryImpl) GetAccomodationById(ctx context.Context, uuid string) (*Accomodation, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.GetAccomodationByUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	return FromGetAccomodationRow(ent), nil
}

func (repo *repositoryImpl) GetAllAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.FindAllAccomodationsByBoardUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	data := []*Accomodation{}
	for _, val := range ent {
		data = append(data, FromFindAccomodationsRow(val))
	}

	return data, nil
}

func (repo *repositoryImpl) UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateAccomodationByUuidParams{
		Url:              data.Url,
		Title:            data.Title,
		ImageUrl:         data.ImageUrl,
		Notes:            data.Notes,
		Status:           data.Status,
		BookingReference: data.BookingReference,
		Selected:         data.Selected,
		Uuid:             id,
		Price:            utility.NumericFromString(data.Price),
		Currency:         utility.NullCurrencyFromPtr(data.Currency),
		Description:      data.Description,
	}

	return repo.queries.UpdateAccomodationByUuid(ctx, params)
}

func (repo *repositoryImpl) DeleteAccomodationById(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteAccomodationByUuid(ctx, id)
}
