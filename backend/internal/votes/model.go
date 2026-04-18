package votes

import "github.com/PegasusMKD/travel-dream-board/internal/db"

type Vote struct {
	Uuid        string     `json:"uuid"`
	UserUuid    string     `json:"user_uuid"`
	Rank        int32      `json:"rank" binding:"required"`
	VotedOn     db.VotedOn `json:"voted_on" binding:"required"`
	VotedOnUuid string     `json:"voted_on_uuid" binding:"required"`
}

func FromEntity(entity db.Vote) *Vote {
	return &Vote{
		Uuid:        entity.Uuid.String(),
		UserUuid:    entity.VotedBy.String(),
		Rank:        entity.Rank,
		VotedOn:     entity.VotedOn,
		VotedOnUuid: entity.VotedOnUuid.String(),
	}
}
