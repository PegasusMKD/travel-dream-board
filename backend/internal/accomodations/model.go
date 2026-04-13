package accomodations

import (
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Accomodation struct {
	Uuid             string                 `json:"uuid"`
	Url              string                 `json:"url" binding:"required"`
	Title            string                 `json:"title" binding:"required"`
	BoardUuid        string                 `json:"board_uuid" binding:"required"`
	ImageUrl         *string                `json:"image_url"`
	Notes            *string                `json:"notes"`
	Status           db.AccomodationsStatus `json:"status"`
	BookingReference *string                `json:"booking_reference"`
	Selected         bool                   `json:"selected"`
}

type AggregatedAccomodation struct {
	Accomodation

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
