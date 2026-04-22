package votes

import "github.com/PegasusMKD/travel-dream-board/internal/db"

type Vote struct {
	Uuid        string     `json:"uuid"`
	UserUuid    string     `json:"user_uuid"`
	UserName    string     `json:"user_name,omitempty"`
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

func FromRow(row db.FindAllVotesByVotedOnUuidRow) *Vote {
	return &Vote{
		Uuid:        row.Uuid.String(),
		UserUuid:    row.VotedBy.String(),
		UserName:    row.UserName,
		Rank:        row.Rank,
		VotedOn:     row.VotedOn,
		VotedOnUuid: row.VotedOnUuid.String(),
	}
}
