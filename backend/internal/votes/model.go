package votes

import "github.com/PegasusMKD/travel-dream-board/internal/db"

type Vote struct {
	Uuid        string     `json:"uuid"`
	VotedBy     string     `json:"voted_by" binding:"required"`
	Rank        int32      `json:"rank" binding:"required"`
	VotedOn     db.VotedOn `json:"voted_on" binding:"required"`
	VotedOnUuid string     `json:"voted_on_uuid" binding:"required"`
}

func FromEntity(entity db.Vote) *Vote {
	return &Vote{
		Uuid:        entity.Uuid.String(),
		VotedBy:     entity.VotedBy,
		Rank:        entity.Rank,
		VotedOn:     entity.VotedOn,
		VotedOnUuid: entity.VotedOnUuid.String(),
	}
}
