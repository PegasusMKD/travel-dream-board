package activities

import (
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Activity struct {
	Uuid             string              `json:"uuid"`
	Url              string              `json:"url" binding:"required"`
	Title            string              `json:"title" binding:"required"`
	BoardUuid        string              `json:"board_uuid" binding:"required"`
	ImageUrl         *string             `json:"image_url"`
	Notes            *string             `json:"notes"`
	Status           db.ActivitiesStatus `json:"status"`
	BookingReference *string             `json:"booking_reference"`
	Selected         bool                `json:"selected"`
}

func FromEntity(entity db.Activity) *Activity {
	return &Activity{
		Uuid:             entity.Uuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
	}
}

type AggregatedActivity struct {
	Activity

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
