package sharetokens

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type ShareToken struct {
	Token     string    `json:"token"`
	BoardUuid string    `json:"board_uuid"`
	CreatedAt time.Time `json:"created_at"`
}

func FromEntity(ent *db.ShareToken) *ShareToken {
	return &ShareToken{
		Token:     ent.Token,
		BoardUuid: ent.BoardUuid.String(),
		CreatedAt: ent.CreatedAt.Time,
	}
}
