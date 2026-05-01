package memories

import (
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Memory struct {
	Uuid       string    `json:"uuid"`
	BoardUuid  string    `json:"board_uuid"`
	UploadedBy string    `json:"uploaded_by"`
	ImageUrl   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
}

func FromEntity(entity db.Memory) *Memory {
	return &Memory{
		Uuid:       entity.Uuid.String(),
		BoardUuid:  entity.BoardUuid.String(),
		UploadedBy: entity.UploadedBy.String(),
		ImageUrl:   entity.ImageUrl,
		CreatedAt:  entity.CreatedAt.Time,
	}
}
