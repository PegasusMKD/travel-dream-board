package transport

import (
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Transport struct {
	Uuid             string             `json:"uuid"`
	Url              string             `json:"url" binding:"required"`
	Title            string             `json:"title" binding:"required"`
	BoardUuid        string             `json:"board_uuid" binding:"required"`
	ImageUrl         *string            `json:"image_url"`
	Notes            *string            `json:"notes"`
	Status           db.TransportStatus `json:"status"`
	BookingReference *string            `json:"booking_reference"`
	Selected         bool               `json:"selected"`
	UserUuid         string             `json:"user_uuid"`
}

func FromEntity(entity db.Transport) *Transport {
	return &Transport{
		Uuid:             entity.Uuid.String(),
		BoardUuid:        entity.BoardUuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
	}
}

type AggregatedTransport struct {
	Transport

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
