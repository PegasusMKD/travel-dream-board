package votes

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateVote(ctx context.Context, data *Vote) (*Vote, error)
	FindAllVotesByRelatedEntity(ctx context.Context, _type db.VotedOn, uuid string) ([]*Vote, error)
	UpdateVoteByUuid(ctx context.Context, uuid string, rank int32) error
	DeleteVoteByUuid(ctx context.Context, uuid string) error
}

type votesRepositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &votesRepositoryImpl{
		queries: queries,
	}
}

func (repo *votesRepositoryImpl) CreateVote(ctx context.Context, data *Vote) (*Vote, error) {
	votedOnId, err := utility.UuidFromString(data.VotedOnUuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", data.VotedOnUuid, "error", err)
		return nil, err
	}

	params := db.CreateVoteParams{
		VotedBy:     data.VotedBy,
		Rank:        data.Rank,
		VotedOn:     data.VotedOn,
		VotedOnUuid: votedOnId,
	}

	ent, err := repo.queries.CreateVote(ctx, params)
	if err != nil {
		log.Error(
			"Failed creating entity",
			"voted_by", data.VotedBy,
			"rank", data.Rank,
			"voted_on", data.VotedOn,
			"voted_on_uuid", data.VotedOnUuid,
			"error", err,
		)
		return nil, err
	}

	return FromEntity(ent), nil
}

func (repo *votesRepositoryImpl) FindAllVotesByRelatedEntity(ctx context.Context, _type db.VoteedOn, uuid string) ([]*Vote, error) {
	votedOnId, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	params := db.FindAllVotesByVotedOnUuidParams{
		VotedOn:     _type,
		VotedOnUuid: votedOnId,
	}

	entities, err := repo.queries.FindAllVotesByVotedOnUuid(ctx, params)
	if err != nil {
		log.Error("Failed fetching votes", "type", _type, "voted_on_uuid", uuid, "error", err)
		return nil, err
	}

	data := []*Vote{}
	for _, val := range entities {
		data = append(data, FromEntity(val))
	}

	return data, nil
}

func (repo *votesRepositoryImpl) UpdateVoteByUuid(ctx context.Context, uuid string, rank int32) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateVoteByUuidParams{
		Uuid: id,
		Rank: rank,
	}

	return repo.queries.UpdateVoteByUuid(ctx, params)
}

func (repo *votesRepositoryImpl) DeleteVoteByUuid(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteVoteByUuid(ctx, id)
}
