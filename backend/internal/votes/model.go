package votes

import "github.com/PegasusMKD/travel-dream-board/internal/db"

type Vote struct {
	Uuid        string     `json:"uuid"`
	VotedBy     string     `json:"voted_by" binding:"required"`
	Rank        int32      `json:"rank" binding:"required"`
	VotedOn     db.VotedOn `json:"voted_on" binding:"required"`
	VotedOnUuid string     `json:"voted_on_uuid" binding:"required"`
}
