package accomodations

import (
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Accomodation struct {
	Uuid string `json:"uuid"`

	Url       string  `json:"url" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	BoardUuid string  `json:"board_uuid" binding:"required"`
	ImageUrl  *string `json:"image_url"`

	Notes            *string                `json:"notes"`
	Status           db.AccomodationsStatus `json:"status"`
	BookingReference *string                `json:"booking_reference"`
	Selected         bool                   `json:"selected"`
	UserUuid         string                 `json:"user_uuid"`
	AvgRating        float64                `json:"avg_rating"`
	RatingCount      int32                  `json:"rating_count"`

	Price       *string          `json:"price"`
	Currency    *db.CurrencyCode `json:"currency"`
	Description *string          `json:"description"`
}

func FromEntity(entity db.Accomodation) *Accomodation {
	return &Accomodation{
		Uuid:             entity.Uuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		BoardUuid:        entity.BoardUuid.String(),
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
		Price:            utility.NumericToString(entity.Price),
		Currency:         utility.CurrencyPtrFromNull(entity.Currency),
		Description:      entity.Description,
	}
}

func FromGetAccomodationRow(entity db.GetAccomodationByUuidRow) *Accomodation {
	return &Accomodation{
		Uuid:             entity.Uuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		BoardUuid:        entity.BoardUuid.String(),
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
		AvgRating:        entity.AvgRating,
		RatingCount:      entity.RatingCount,
		Price:            utility.NumericToString(entity.Price),
		Currency:         utility.CurrencyPtrFromNull(entity.Currency),
		Description:      entity.Description,
	}
}

func FromFindAccomodationsRow(entity db.FindAllAccomodationsByBoardUuidRow) *Accomodation {
	return &Accomodation{
		Uuid:             entity.Uuid.String(),
		Url:              entity.Url,
		Title:            entity.Title,
		BoardUuid:        entity.BoardUuid.String(),
		ImageUrl:         entity.ImageUrl,
		Notes:            entity.Notes,
		Status:           entity.Status,
		BookingReference: entity.BookingReference,
		Selected:         entity.Selected,
		UserUuid:         entity.UserUuid.String(),
		AvgRating:        entity.AvgRating,
		RatingCount:      entity.RatingCount,
		Price:            utility.NumericToString(entity.Price),
		Currency:         utility.CurrencyPtrFromNull(entity.Currency),
		Description:      entity.Description,
	}
}

type AggregatedAccomodation struct {
	Accomodation

	Comments []*comments.Comment `json:"comments"`
	Votes    []*votes.Vote       `json:"votes"`
}
